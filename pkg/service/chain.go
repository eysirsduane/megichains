package service

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"megichains/pkg/biz"
	"megichains/pkg/converter"
	"megichains/pkg/entity"
	"megichains/pkg/erc20"
	"megichains/pkg/global"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

var mu sync.RWMutex
var erc20ABI = `[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"}]`

type ChainService struct {
	cfg     *global.BackendesConfig
	db      *gorm.DB
	running bool
}

func NewChainService(cfg *global.BackendesConfig, db *gorm.DB) *ChainService {
	return &ChainService{cfg: cfg, db: db}
}

func (s *ChainService) newEvmClient(chain global.ChainName) (client *ethclient.Client, err error) {
	port := ""
	switch chain {
	case global.ChainNameBsc:
		port = s.cfg.Bsc.WssNetwork2
	case global.ChainNameEth:
		port = fmt.Sprintf("%v%v", s.cfg.Eth.WssNetwork, s.cfg.Eth.ApiKey)
		// port = s.cfg.Eth.WssNetwork2
	// case global.ChainNameTron:
	// 	port = s.cfg.Tron.WssNetwork
	default:
		return nil, fmt.Errorf("unknown chain:%v", chain)
	}

	client, err = ethclient.Dial(port)
	if err != nil {
		logx.Errorf("evm client Dial failed, chain:%v, err:%v", chain, err)
		return
	}

	return
}

func (s *ChainService) EncryptEthPrivateKey() {
	ctx := context.Background()
	addrs, err := gorm.G[entity.Address](s.db).Where("chain = ?", global.ChainNameEvm).Order("id asc").Find(ctx)
	if err != nil {
		logx.Errorf("encrypt private key address find failed, err:%v", err)
		return
	}
	for _, addr := range addrs {
		encrypted, err := global.EncryptEthPrivateKey(addr.PrivateKey, global.SecretKey)
		if err != nil {
			logx.Errorf("encrypt private key encrypted failed, addr:%v, err:%v", addr.Address, err)
			continue
		}
		addr.PrivateKey = encrypted
		if err := s.db.Save(&addr).Error; err != nil {
			logx.Errorf("encrypt private key save encrypted key failed, addr:%v, err:%v", addr.Address, err)
		}
	}
}

func (s *ChainService) ScanAddressesFunds() {
	if s.running {
		return
	}
	s.running = true
	defer func() {
		s.running = false
	}()

	ctx := context.Background()
	addrs, err := gorm.G[entity.Address](s.db).Order("id asc").Find(ctx)
	if err != nil {
		logx.Errorf("scan addresses funds cron address find failed, err:%v", err)
		return
	}

	for _, addr := range addrs {
		switch global.ChainName(addr.Chain) {
		case global.ChainNameEvm:

			s.EvmFunds(addr.Address, global.ChainNameBsc)
			s.EvmFunds(addr.Address, global.ChainNameEth)
		case global.ChainNameTron:
			s.TronFunds(addr.Address)
		default:
			logx.Errorf("scan addresses funds found unknown chain, chain:%v", addr.Chain)
		}

		time.Sleep(time.Millisecond * 500)
	}
}

func (s *ChainService) EvmFunds(addr string, chain global.ChainName) {
	address := common.HexToAddress(addr)

	taddr, err := s.getContractAddress(global.CurrencyTypoUsdt, string(chain))
	if err != nil {
		logx.Errorf("evm fund get usdt contract address failed, err:%v", err)
		return
	}
	caddr, err := s.getContractAddress(global.CurrencyTypoUsdc, string(chain))
	if err != nil {
		logx.Errorf("evm fund get usdc contract address failed, err:%v", err)
		return
	}

	usdtAddr := common.HexToAddress(taddr)
	usdcAddr := common.HexToAddress(caddr)

	usdt, err := s.ERC20Balance(chain, usdtAddr, address)
	if err != nil {
		logx.Errorf("evm fund get usdt balance failed, chain:%v, err:%v", chain, err)
		return
	}
	usdc, err := s.ERC20Balance(chain, usdcAddr, address)
	if err != nil {
		logx.Errorf("evm fund get usdc balance failed, chain:%v, err:%v", chain, err)
		return
	}

	s.updateDB(addr, chain, usdt, usdc)
}

func (s *ChainService) TronFunds(addr string) {
	taddr, err := s.getContractAddress(global.CurrencyTypoUsdt, string(global.ChainNameTron))
	if err != nil {
		logx.Errorf("tron fund get usdc contract address failed, err:%v", err)
		return
	}
	caddr, err := s.getContractAddress(global.CurrencyTypoUsdc, string(global.ChainNameTron))
	if err != nil {
		logx.Errorf("tron fund get usdc contract address failed, err:%v", err)
		return
	}

	usdt, err := s.GetTRC20(addr, taddr)
	if err != nil {
		logx.Errorf("tron fund get usdt balance failed, chain:%v, err:%v", global.ChainNameTron, err)
		return
	}
	usdc, err := s.GetTRC20(addr, caddr)
	if err != nil {
		logx.Errorf("tron fund get usdc balance failed, chain:%v, err:%v", global.ChainNameTron, err)
		return
	}

	s.updateDB(addr, global.ChainNameTron, usdt, usdc)
}

func (s *ChainService) updateDB(addr string, chain global.ChainName, usdt, usdc *big.Int) {
	ctx := context.Background()

	db := gorm.G[entity.Address](s.db)
	address, err := db.Where("address = ?", addr).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logx.Errorf("address fund find unknown address, addr:%v", addr)
	} else if err != nil {
		logx.Errorf("address fund first failed, err:%v", err)
	} else {
		switch chain {
		case global.ChainNameBsc:
			address.BscUsdt = global.Amount(usdt.Int64(), global.AmountTypoBsc)
			address.BscUsdc = global.Amount(usdc.Int64(), global.AmountTypoBsc)
		case global.ChainNameEth:
			address.EthUsdt = global.Amount(usdt.Int64(), global.AmountTypoEth)
			address.EthUsdc = global.Amount(usdc.Int64(), global.AmountTypoEth)
		case global.ChainNameTron:
			address.TronUsdt = global.Amount(usdt.Int64(), global.AmountTypoTron)
			address.TronUsdc = global.Amount(usdc.Int64(), global.AmountTypoTron)
		}

		_, err = db.Updates(ctx, address)
		if err != nil {
			logx.Errorf("address fund first updates failed, err:%v", err)
		}
	}
}

func (s *ChainService) post(path string, body any) (byts []byte, err error) {
	b, _ := json.Marshal(body)
	resp, err := http.Post(s.cfg.Tron.HttpNetwork+path, "application/json", bytes.NewBuffer(b))
	if err != nil {
		logx.Errorf("scan funds tron http post failed, err:%v", err)
		return
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (s *ChainService) Base58ToHex(addr string) (string, error) {
	decoded := base58.Decode(addr)
	if len(decoded) < 5 {
		return "", errors.New("invalid address length")
	}

	payload := decoded[:len(decoded)-4] // 去掉 checksum
	checksum := decoded[len(decoded)-4:]

	hash := sha256.Sum256(payload)
	hash = sha256.Sum256(hash[:])

	if string(checksum) != string(hash[:4]) {
		return "", errors.New("checksum mismatch")
	}

	if payload[0] != 0x41 {
		return "", errors.New("invalid TRON version byte")
	}

	return hex.EncodeToString(payload), nil
}

func (s *ChainService) GetTRC20(addr, contract string) (balance *big.Int, err error) {
	userHex, _ := s.Base58ToHex(addr)
	contractHex, _ := s.Base58ToHex(contract)

	param := "000000000000000000000000" + userHex[2:]

	resp, err := s.post("/wallet/triggersmartcontract", map[string]any{
		"owner_address":     userHex,
		"contract_address":  contractHex,
		"function_selector": "balanceOf(address)",
		"parameter":         param,
	})
	if err != nil {
		return
	}

	var r struct {
		ConstantResult []string `json:"constant_result"`
	}
	json.Unmarshal(resp, &r)

	if len(r.ConstantResult) == 0 {
		return big.NewInt(0), nil
	}

	balance = new(big.Int)
	balance.SetString(r.ConstantResult[0], 16)

	return
}

func (s *ChainService) getContractAddress(currency global.CurrencyTypo, chain string) (caddr string, err error) {
	switch currency {
	case global.CurrencyTypoUsdt:
		for _, addr := range s.cfg.ContractAddresses {
			if strings.EqualFold(addr.Chain, string(chain)) && strings.EqualFold(addr.Currency, string(currency)) {
				caddr = addr.Address
			}
		}
	case global.CurrencyTypoUsdc:
		for _, addr := range s.cfg.ContractAddresses {
			if strings.EqualFold(addr.Chain, string(chain)) && strings.EqualFold(addr.Currency, string(currency)) {
				caddr = addr.Address
			}
		}
	default:
		err = fmt.Errorf("contract address get failed...!, chain:%v, currency:%v", chain, currency)
		return
	}

	return
}

func (s *ChainService) ERC20Balance(chain global.ChainName, uaddr common.Address, owner common.Address) (*big.Int, error) {
	var cli *ethclient.Client
	switch chain {
	case global.ChainNameBsc:
		cli, _ = s.newEvmClient(global.ChainNameBsc)
	case global.ChainNameEth:
		cli, _ = s.newEvmClient(global.ChainNameEth)
	}
	defer cli.Close()

	erc20c, err := erc20.NewErc20(uaddr, cli)
	if err != nil {
		logx.Errorf("evm collect new erc20 instance failed, uaddr:%v, err:%v", uaddr.Hex(), err)
		err = biz.AddressFundCollectNewErc20InstanceFailed
		return nil, err
	}

	balance, err := erc20c.BalanceOf(&bind.CallOpts{}, owner)
	if err != nil {
		logx.Errorf("evm collect get erc20 balance failed, uaddr:%v, err:%v", uaddr.Hex(), err)
		err = biz.AddressFundCollectErc20TransferFailed
		return nil, err
	}

	return balance, nil
}

func (s *ChainService) Collect(ctx context.Context, uid string, req *converter.AddressFundCollectReq) (resp *converter.AddressFundCollectResp, err error) {
	collect := &entity.AddressFundCollect{
		UserId:         uid,
		AddressGroupId: req.AddressGroupId,
		Chain:          req.Chain,
		Currency:       req.Currency,
		AmountMin:      req.AmountMin,
		FeeMax:         req.FeeMax,
		Status:         string(global.CollectStatusCreated),
	}

	err = gorm.G[entity.AddressFundCollect](s.db).Create(ctx, collect)
	if err != nil {
		logx.Errorf("address fund collect log create failed, err:%v", err)
		err = biz.AddressFundCollectLogCreateFailed
		return
	}

	chain := global.ChainName(req.Chain)
	currency := global.CurrencyTypo(req.Currency)
	switch chain {
	case global.ChainNameBsc, global.ChainNameEth:
		err = s.EvmCollect(ctx, collect, req.AddressGroupId, chain, currency, req.AmountMin, req.FeeMax)
	case global.ChainNameTron:
		err = s.TronCollect(currency, req.AmountMin, req.FeeMax, req.AddressGroupId)
	default:
		logx.Errorf("collect found unknown chain, chain:%v", chain)
		err = biz.AddressFundCollectUnknownChain
	}

	return
}

func (s *ChainService) EvmCollect(ctx context.Context, collect *entity.AddressFundCollect, gid int64, chain global.ChainName, currency global.CurrencyTypo, amin, fmax float64) (err error) {
	var cli *ethclient.Client
	switch chain {
	case global.ChainNameBsc:
		cli, _ = s.newEvmClient(global.ChainNameBsc)
	case global.ChainNameEth:
		cli, _ = s.newEvmClient(global.ChainNameEth)
	}

	receiver := &entity.Address{}
	err = s.db.Model(&entity.Address{}).Where("chain = ? and typo = ?", global.ChainNameEvm, global.AddressTypoCollect).First(receiver).Error
	if err != nil {
		logx.Errorf("evm collect get receiver address failed, chain:%v, err:%v", chain, err)
		err = biz.AddressFundCollectReceiverAddressNotFound
		cli.Close()
		return
	}

	collect.ReceiverAddress = receiver.Address

	jquery := `"group_id" = ? and "chain" = ? and typo = ? `
	switch chain {
	case global.ChainNameBsc:
		switch currency {
		case global.CurrencyTypoUsdt:
			jquery += " and bsc_usdt >= ?"
		case global.CurrencyTypoUsdc:
			jquery += " and bsc_usdc >= ?"
		}
	case global.ChainNameEth:
		switch currency {
		case global.CurrencyTypoUsdt:
			jquery += " and eth_usdt >= ?"
		case global.CurrencyTypoUsdc:
			jquery += " and eth_usdc >= ?"
		}
	}

	froms := make([]*entity.Address, 0)
	err = s.db.Model(entity.Address{}).Where(jquery, gid, global.ChainNameEvm, global.AddressTypoIn, amin).Find(&froms).Error
	if err != nil {
		logx.Errorf("evm collect get from address failed, group:%v, chain:%v, err:%v", gid, chain, err)
		err = biz.AddressFundCollectFromAddressNotFound
		cli.Close()
		return
	}

	collect.TotalCount = int64(len(froms))

	for _, from := range froms {
		go func() {
			amount := 0.000011

			log := &entity.AddressFundCollectLog{
				CollectId:       collect.Id,
				Chain:           string(chain),
				Currency:        string(currency),
				Amount:          amount,
				FromAddress:     from.Address,
				ReceiverAddress: receiver.Address,
				Status:          string(global.CollectLogStatusCreated),
			}

			err = s.db.Save(log).Error
			if err != nil {
				logx.Errorf("evm collect log create failed, collect:%v, from:%v, to:%v, err:%v", collect.Id, from.Address, receiver.Address, err)
				return
			}

			ctx := context.Background()
			tx, err := s.sendEvmTransaction(ctx, cli, chain, currency, from, receiver, amount, fmax)
			if err != nil {
				logx.Errorf("evm collect send transaction failed, err:%v", err)

				log.Status = string(global.CollectLogStatusFailed)
				log.Description = err.Error()
				s.db.Save(log)

				err = biz.AddressFundCollectSendTxFailed
				return
			}

			log.TransactionId = tx.Hash().Hex()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			receipt, err := bind.WaitMined(ctx, cli, tx)
			if err != nil {
				logx.Errorf("evm collect wait mined failed, from:%v, err:%v", from.Address, err)

				log.Status = string(global.CollectLogStatusFailed)
				log.Description = err.Error()
				s.db.Save(log)

				err = biz.AddressFundCollectWaitMinedFailed
				return
			}

			if receipt.Status == types.ReceiptStatusSuccessful {
				logx.Infof("🎉 evm collect transfer success, chain:%v, currency:%v, from:%v, to:%v,  txid:%v", chain, currency, from.Address, receiver.Address, tx.Hash().Hex())
				log.Status = string(global.CollectLogStatusSuccess)
			} else {
				logx.Errorf("evm collect transfer failed, chain:%v, currency:%v, from:%v, to:%v,  txid:%v, status:%d", chain, currency, from.Address, receiver.Address, tx.Hash().Hex(), receipt.Status)
				log.Status = string(global.CollectLogStatusFailed)
				log.Description = fmt.Sprintf("evm collect log transfer failed, receipt status:%v", receipt.Status)
			}

			err = s.db.Save(log).Error
			if err != nil {
				logx.Errorf("evm collect log save failed, collect:%v, from:%v, to:%v, err:%v", collect.Id, from.Address, receiver.Address, err)
			}

			//更新链上真实余额
			switch chain {
			case global.ChainNameBsc:
				s.EvmFunds(from.Address, global.ChainNameBsc)
			case global.ChainNameEth:
				s.EvmFunds(from.Address, global.ChainNameEth)
			default:
				logx.Errorf("evm collect scan address fund found unknown chain, chain:%v, from:%v", from.Chain, from.Address)
			}
		}()
	}

	collect.Status = string(global.CollectStatusProcessing)
	err = s.db.Save(collect).Error
	if err != nil {
		logx.Errorf("address fund collect log updates failed, err:%v", err)
		err = biz.AddressFundCollectLogUpdateFailed
		return
	}

	return
}

func (s *ChainService) sendEvmTransaction(ctx context.Context, cli *ethclient.Client, chain global.ChainName, currency global.CurrencyTypo, from, receiver *entity.Address, amount float64, fmax float64) (tx *types.Transaction, err error) {
	decrypted, err := global.DecryptEthPrivateKey(from.PrivateKey, global.SecretKey)
	if err != nil {
		return
	}
	privateKey, err := crypto.HexToECDSA(decrypted)
	if err != nil {
		return
	}

	publicKeyECDSA := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	caddr, err := s.getContractAddress(currency, string(chain))
	if err != nil {
		return
	}
	uaddr := common.HexToAddress(caddr)
	toAddress := common.HexToAddress(receiver.Address)

	sun := int64(0)
	switch chain {
	case global.ChainNameBsc:
		sun = global.Sun(amount, global.AmountTypoBsc)
	case global.ChainNameEth:
		sun = global.Sun(amount, global.AmountTypoEth)
	}

	efee, tcap, fcap, glimit, err := erc20.EstimateTransactionFee(ctx, cli, chain, uaddr, fromAddress, toAddress, sun)
	if err != nil {
		return
	}

	overed, err := erc20.CheckFeeOverLimit(efee, fmax)
	if err != nil {
		return
	}
	if overed {
		err = fmt.Errorf("evm collect fee limit overflow, efee:%v", efee)
		return
	}

	cid, err := cli.ChainID(ctx)
	if err != nil {
		return
	}

	// clearPendingTransaction(ctx, cli, cid, fromAddress, privateKey)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, cid)
	if err != nil {
		return
	}

	nonce, err := cli.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = glimit
	switch chain {
	case global.ChainNameBsc:
		auth.GasPrice = fcap
	case global.ChainNameEth:
		auth.GasTipCap = tcap
		auth.GasFeeCap = fcap
	}

	erc20c, err := erc20.NewErc20(uaddr, cli)
	if err != nil {
		return
	}

	tx, err = erc20c.Transfer(auth, toAddress, big.NewInt(sun))
	if err != nil {
		return
	}

	return
}

func estimateGasLimit(ctx context.Context, cli ethclient.Client, toAddress, fromAddress, uaddr common.Address, sun int64) (glimit uint64, err error) {
	txlimit, err := erc20.PackTransfer(toAddress, big.NewInt(sun))
	if err != nil {
		logx.Errorf("evm collect pack transfer for estimate gas limit failed, err:%v", err)
		err = biz.AddressFundCollectPackTransferFailed
		return
	}
	msg := ethereum.CallMsg{
		From:  fromAddress,   // 付款账户
		To:    &uaddr,        // 目标合约地址
		Data:  txlimit,       // 合约调用数据
		Value: big.NewInt(0), // ERC20转账不转账ETH，值为0
	}
	glimit, err = cli.EstimateGas(ctx, msg)
	if err != nil {
		logx.Errorf("=== evm collect estimate gas limit failed ===, err:%v", err)
		glimit = 120000
	}
	glimit = uint64(float64(glimit) * 1.25)

	return
}

func clearPendingTransaction(ctx context.Context, cli *ethclient.Client, cid *big.Int, addr common.Address, prikey *ecdsa.PrivateKey) {
	confirmedNonce, err := cli.NonceAt(ctx, addr, nil)
	if err != nil {
		logx.Errorf("获取确认Nonce失败: %v, cnonce:%v", err, confirmedNonce)
		return
	}

	gasTipCap := big.NewInt(100_000_000_000) // 100 Gwei 矿工小费
	gasFeeCap := big.NewInt(300_000_000_000) // 300 Gwei 总费用上限
	gasLimit := uint64(21000)                // 原生转账标准GasLimit

	// 4. 构造覆盖交易：相同阻塞Nonce + 0ETH转给自己
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   cid,
		Nonce:     confirmedNonce, // 关键：使用阻塞的Nonce
		To:        &addr,          // 收款方=自己，无资产转移
		Value:     big.NewInt(0),
		Gas:       gasLimit,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
	})

	// 5. 签名交易
	signer := types.LatestSignerForChainID(cid)
	signedTx, err := types.SignTx(tx, signer, prikey)
	if err != nil {
		logx.Errorf("交易签名失败: %v", err)
		return
	}

	// 6. 发送覆盖交易
	err = cli.SendTransaction(ctx, signedTx)
	if err != nil {
		logx.Errorf("发送覆盖交易失败: %v", err)
		return
	}

	logx.Infof("clean txid:%v", signedTx.Hash().Hex())
}

func (s *ChainService) TronCollect(currency global.CurrencyTypo, amin, fmax float64, gid int64) (err error) {
	return
}
