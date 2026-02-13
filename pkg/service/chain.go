package service

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"megichains/pkg/biz"
	"megichains/pkg/converter"
	"megichains/pkg/crypt"
	"megichains/pkg/entity"
	"megichains/pkg/erc20"
	"megichains/pkg/global"
	"net/http"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/kslamph/tronlib/pb/api"
	"github.com/kslamph/tronlib/pkg/client"
	"github.com/kslamph/tronlib/pkg/signer"
	ttypes "github.com/kslamph/tronlib/pkg/types"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ChainService struct {
	cfg          *global.BackendesConfig
	db           *gorm.DB
	running      bool
	cli          any
	BscClinet    *ethclient.Client
	EthClient    *ethclient.Client
	SolanaClient *rpc.Client
}

func NewChainService(cfg *global.BackendesConfig, db *gorm.DB) *ChainService {
	return &ChainService{cfg: cfg, db: db}
}

func (s *ChainService) initChainClient(chain global.ChainName) (err error) {
	switch chain {
	case global.ChainNameBsc:
		client, ierr := ethclient.Dial(s.cfg.Bsc.WssNetwork2)
		if ierr != nil {
			logx.Errorf("chain client Dial failed, chain:%v, err:%v", chain, err)
			return
		}

		s.cli = client
	case global.ChainNameEth:
		client, ierr := ethclient.Dial(s.cfg.Eth.WssNetwork2)
		if ierr != nil {
			logx.Errorf("chain client Dial failed, chain:%v, err:%v", chain, err)
			return
		}

		s.cli = client
	case global.ChainNameTron:
	case global.ChainNameSolana:
		s.cli = rpc.New(s.cfg.Solana.GrpcNetwork)
	default:
		return fmt.Errorf("chain client found unknown chain:%v", chain)
	}

	return
}

func (s *ChainService) EncryptPrivateKey() {
	ctx := context.Background()
	addrs, err := gorm.G[entity.Address](s.db).Where("chain = ? and id = 1102", global.ChainNameSolana).Order("id asc").Find(ctx)
	if err != nil {
		logx.Errorf("encrypt private key address find failed, err:%v", err)
		return
	}
	for _, addr := range addrs {
		encrypted, err := crypt.EncryptEthPrivateKey(addr.PrivateKey, crypt.PrivateKeySecretSalt)
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

func (s *ChainService) ReEncryptPrivateKey() {
	ctx := context.Background()
	addrs, err := gorm.G[entity.Address](s.db).Where("chain = ?", global.ChainNameTron).Order("id asc").Find(ctx)
	if err != nil {
		logx.Errorf("encrypt private key address find failed, err:%v", err)
		return
	}
	for _, addr := range addrs {
		decrypted, err := crypt.DecryptEthPrivateKey(addr.PrivateKey, crypt.PrivateKeySecretSalt)
		if err != nil {
			logx.Errorf("encrypt private key encrypted failed, addr:%v, err:%v", addr.Address, err)
			continue
		}
		encrypted, err := crypt.Encrypt(decrypted, crypt.PrivateKeySecretPassword, crypt.PrivateKeySecretSalt)
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

func (s *ChainService) ScanAddressesFunds(chain global.ChainName) {
	if s.running {
		return
	}
	s.running = true
	defer func() {
		s.running = false
	}()

	ctx := context.Background()
	addrs, err := gorm.G[entity.Address](s.db).Where("chain = ?", chain).Order("id asc").Find(ctx)
	if err != nil {
		logx.Errorf("scan addresses funds cron address find failed, err:%v", err)
		return
	}

	err = s.initChainClient(chain)
	if err != nil {
		logx.Errorf("collect new chain client failed, chain:%v, err:%v", global.ChainNameBsc, err)
		return
	}

	for _, addr := range addrs {
		switch global.ChainName(addr.Chain) {
		case global.ChainNameEvm:
			s.EvmFunds(addr.Address, global.ChainNameBsc)
			s.EvmFunds(addr.Address, global.ChainNameEth)
		case global.ChainNameTron:
			s.TronFunds(addr.Address, global.ChainNameTron)
		case global.ChainNameSolana:
			s.SolanaFunds(addr.Address, global.ChainNameSolana)
		default:
			logx.Errorf("scan addresses funds found unknown chain, chain:%v", addr.Chain)
		}

		time.Sleep(time.Millisecond * 300)
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

	usdt, err := s.GetERC20Balance(chain, usdtAddr, address)
	if err != nil {
		logx.Errorf("evm fund get usdt balance failed, chain:%v, err:%v", chain, err)
		return
	}
	usdc, err := s.GetERC20Balance(chain, usdcAddr, address)
	if err != nil {
		logx.Errorf("evm fund get usdc balance failed, chain:%v, err:%v", chain, err)
		return
	}

	s.updateDB(addr, chain, usdt, usdc)
}

func (s *ChainService) GetERC20Balance(chain global.ChainName, uaddr common.Address, owner common.Address) (balance int64, err error) {
	cli, ok := s.cli.(*ethclient.Client)
	if !ok {
		err = biz.ConvertClientFailed
		return
	}

	erc20c, err := erc20.NewErc20(uaddr, cli)
	if err != nil {
		logx.Errorf("evm collect new erc20 instance failed, uaddr:%v, err:%v", uaddr.Hex(), err)
		err = biz.AddressFundCollectNewErc20InstanceFailed
		return
	}

	b, err := erc20c.BalanceOf(&bind.CallOpts{}, owner)
	if err != nil {
		logx.Errorf("evm collect get erc20 balance failed, uaddr:%v, err:%v", uaddr.Hex(), err)
		err = biz.AddressFundCollectErc20TransferFailed
		return
	}

	balance = b.Int64()

	return
}

func (s *ChainService) TronFunds(addr string, chain global.ChainName) {
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

	usdt, err := s.GetTRC20Balance(addr, taddr)
	if err != nil {
		logx.Errorf("tron fund get usdt balance failed, chain:%v, err:%v", global.ChainNameTron, err)
		return
	}
	usdc, err := s.GetTRC20Balance(addr, caddr)
	if err != nil {
		logx.Errorf("tron fund get usdc balance failed, chain:%v, err:%v", global.ChainNameTron, err)
		return
	}

	s.updateDB(addr, chain, usdt, usdc)
}

func (s *ChainService) GetTRC20Balance(addr, contract string) (balance int64, err error) {
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
		return 0, nil
	}

	b, _ := new(big.Int).SetString(r.ConstantResult[0], 16)
	balance = b.Int64()

	// c, err := client.NewClient(s.cfg.Tron.GrpcNetwork)
	// if err != nil {
	// 	logx.Errorf("tron collect client new clinet failed, collect:%v, err:%v", collect.Id, err)
	// 	err = biz.AddressFundCollectInitClientFailed
	// 	return
	// }

	return
}

func (s *ChainService) SolanaFunds(addr string, chain global.ChainName) {
	taddr, err := s.getContractAddress(global.CurrencyTypoUsdt, string(global.ChainNameSolana))
	if err != nil {
		logx.Errorf("solana fund get usdc contract address failed, err:%v", err)
		return
	}
	caddr, err := s.getContractAddress(global.CurrencyTypoUsdc, string(global.ChainNameSolana))
	if err != nil {
		logx.Errorf("solana fund get usdc contract address failed, err:%v", err)
		return
	}

	usdt, err := s.GetSolanaUsdBalance(addr, taddr)
	if err != nil {
		logx.Errorf("solana fund get usdt balance failed, chain:%v, err:%v", global.ChainNameSolana, err)
	}
	usdc, err := s.GetSolanaUsdBalance(addr, caddr)
	if err != nil {
		logx.Errorf("solana fund get usdc balance failed, chain:%v, err:%v", global.ChainNameSolana, err)
	}

	s.updateDB(addr, chain, usdt, usdc)
}

func (s *ChainService) GetSolanaUsdBalance(addr, contract string) (balance int64, err error) {
	cli, ok := s.cli.(*rpc.Client)
	if !ok {
		err = biz.ConvertClientFailed
		return
	}

	umint := solana.MustPublicKeyFromBase58(contract)
	owner := solana.MustPublicKeyFromBase58(addr)

	uata, _, err := solana.FindAssociatedTokenAddress(owner, umint)
	if err != nil {
		return
	}

	ub, err := cli.GetTokenAccountBalance(context.Background(), uata, rpc.CommitmentFinalized)
	if err != nil {
		return
	}

	b, _ := new(big.Int).SetString(ub.Value.Amount, 10)
	balance = b.Int64()

	return
}

func (s *ChainService) updateDB(addr string, chain global.ChainName, usdt, usdc int64) {
	balance := &entity.AddressBalance{
		Address: addr,
	}

	err := s.db.Model(&entity.AddressBalance{}).Where("address = ?", addr).FirstOrCreate(balance).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logx.Errorf("address fund find unknown address, addr:%v", addr)
	} else if err != nil {
		logx.Errorf("address fund first failed, err:%v", err)
	} else {
		switch chain {
		case global.ChainNameBsc:
			balance.BscUsdt = global.Amount(usdt, global.AmountTypo18e)
			balance.BscUsdc = global.Amount(usdc, global.AmountTypo18e)
		case global.ChainNameEth:
			balance.EthUsdt = global.Amount(usdt, global.AmountTypo6e)
			balance.EthUsdc = global.Amount(usdc, global.AmountTypo6e)
		case global.ChainNameTron:
			balance.TronUsdt = global.Amount(usdt, global.AmountTypo6e)
			balance.TronUsdc = global.Amount(usdc, global.AmountTypo6e)
		case global.ChainNameSolana:
			balance.SolanaUsdt = global.Amount(usdt, global.AmountTypo6e)
			balance.SolanaUsdc = global.Amount(usdc, global.AmountTypo6e)
		}

		err = s.db.Save(balance).Error
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

func (s *ChainService) Collect(ctx context.Context, uid string, req *converter.AddressFundCollectReq) (resp *converter.AddressFundCollectResp, err error) {
	// s.ScanAddressesFunds(global.ChainNameSolana)
	// return
	// s.ReEncryptPrivateKey()
	// return

	cchain := ""
	chain := global.ChainName(req.Chain)
	currency := global.CurrencyTypo(req.Currency)
	if chain == global.ChainNameBsc || chain == global.ChainNameEth {
		cchain = "EVM"
	} else {
		cchain = string(chain)
	}

	err = s.initChainClient(chain)
	if err != nil {
		logx.Errorf("collect new chain client failed, chain:%v, err:%v", chain, err)
		return
	}

	receiver := &entity.Address{}
	err = s.db.Model(&entity.Address{}).Where("chain = ? and typo = ?", cchain, global.AddressTypoCollect).First(receiver).Error
	if err != nil {
		logx.Errorf("collect get receiver address failed, err:%v", err)
		err = biz.AddressFundCollectReceiverAddressNotFound
		return
	}

	var caddr string
	jquery := `addresses.group_id = ? and addresses.chain = ? and addresses.typo = ? `
	var updateAddressFund func(addr string, chain global.ChainName)
	var sendTransaction func(log *entity.AddressFundCollectLog, chain global.ChainName, from, to, prikey, caddr string, amount, fmax float64) (err error)

	switch chain {
	case global.ChainNameBsc:
		switch currency {
		case global.CurrencyTypoUsdt:
			jquery += " and bsc_usdt >= ?"
		case global.CurrencyTypoUsdc:
			jquery += " and bsc_usdc >= ?"
		}

		updateAddressFund = s.EvmFunds
		sendTransaction = s.sendEvmTransaction

		caddr, err = s.getContractAddress(currency, string(global.ChainNameBsc))
		if err != nil {
			return
		}
	case global.ChainNameEth:
		switch currency {
		case global.CurrencyTypoUsdt:
			jquery += " and eth_usdt >= ?"
		case global.CurrencyTypoUsdc:
			jquery += " and eth_usdc >= ?"
		}

		updateAddressFund = s.EvmFunds
		sendTransaction = s.sendEvmTransaction

		caddr, err = s.getContractAddress(currency, string(global.ChainNameEth))
		if err != nil {
			return
		}
	case global.ChainNameTron:
		switch currency {
		case global.CurrencyTypoUsdt:
			jquery += " and address_balances.tron_usdt >= ?"
		case global.CurrencyTypoUsdc:
			jquery += " and address_balances.tron_usdc >= ?"
		}

		updateAddressFund = s.TronFunds
		sendTransaction = s.sendTronTransaction

		caddr, err = s.getContractAddress(currency, string(global.ChainNameTron))
		if err != nil {
			return
		}
	case global.ChainNameSolana:
		switch currency {
		case global.CurrencyTypoUsdt:
			jquery += " and address_balances.solana_usdt >= ?"
		case global.CurrencyTypoUsdc:
			jquery += " and address_balances.solana_usdc >= ?"
		}

		updateAddressFund = s.SolanaFunds
		sendTransaction = s.sendSolanaTransaction

		caddr, err = s.getContractAddress(currency, string(global.ChainNameSolana))
		if err != nil {
			return
		}
	}

	froms := make([]*entity.Address, 0)
	err = s.db.Model(entity.Address{}).Preload("AddressBalance").Joins("left join address_balances on addresses.address = address_balances.address").Where(jquery, req.AddressGroupId, cchain, global.AddressTypoIn, req.AmountMin).Find(&froms).Error
	if err != nil {
		logx.Errorf("collect get from address failed, group:%v, chain:%v, err:%v", req.AddressGroupId, chain, err)
		err = biz.AddressFundCollectFromAddressNotFound
		return
	}

	collect := &entity.AddressFundCollect{
		UserId:          uid,
		AddressGroupId:  req.AddressGroupId,
		Chain:           req.Chain,
		Currency:        req.Currency,
		AmountMin:       req.AmountMin,
		FeeMax:          req.FeeMax,
		Status:          string(global.CollectStatusProcessing),
		ReceiverAddress: receiver.Address,
		TotalCount:      int64(len(froms)),
	}

	err = s.db.Save(collect).Error
	if err != nil {
		logx.Errorf("collect save failed, group:%v, chain:%v, err:%v", req.AddressGroupId, chain, err)
		err = biz.AddressFundCollectSaveFailed
		return
	}

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
				logx.Errorf("collect log create failed, collect:%v, from:%v, to:%v, err:%v", collect.Id, from.Address, receiver.Address, err)
				return
			}

			prikey, err := crypt.Decrypt(from.PrivateKey, crypt.PrivateKeySecretPassword, crypt.PrivateKeySecretSalt)
			if err != nil {
				logx.Errorf("collect decrypted private key failed, from:%v, err:%v", from.Address, err)
				return
			}

			err = sendTransaction(log, chain, from.Address, receiver.Address, prikey, caddr, amount, req.FeeMax)
			if err != nil {
				logx.Errorf("collect send transaction failed, chain:%v, collect:%v, from:%v, to:%v, err:%v", chain, collect.Id, from.Address, receiver.Address, err)
				err = biz.AddressFundCollectSendTxFailed
			}

			err = s.db.Save(log).Error
			if err != nil {
				logx.Errorf("collect log update failed, collect:%v, from:%v, to:%v, err:%v", collect.Id, from.Address, receiver.Address, err)
			}

			updateAddressFund(from.Address, chain)
			updateAddressFund(receiver.Address, chain)
		}()
	}

	return
}

func (s *ChainService) sendEvmTransaction(log *entity.AddressFundCollectLog, chain global.ChainName, from, to, prikey, caddr string, amount, fmax float64) (err error) {
	cli, ok := s.cli.(*ethclient.Client)
	if !ok {
		logx.Errorf("send transaction convert cli failed, chain:%v, log:%v", chain, log.Id)
		return
	}

	tx, err := s.evmTransfer(cli, chain, to, prikey, caddr, amount, fmax)
	if err != nil {
		logx.Errorf("evm collect send transaction failed, err:%v", err)

		log.Status = string(global.CollectLogStatusFailed)
		log.Description = err.Error()

		err = biz.AddressFundCollectSendTxFailed
		return
	}

	log.TransactionId = tx.Hash().Hex()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, cli, tx)
	if err != nil {
		logx.Errorf("evm collect wait mined failed, from:%v, err:%v", from, err)

		log.Status = string(global.CollectLogStatusFailed)
		log.Description = err.Error()

		err = biz.AddressFundCollectWaitMinedFailed
		return
	}

	if receipt.Status == types.ReceiptStatusSuccessful {
		logx.Infof("🎉 evm collect transfer success, chain:%v, currency:%v, from:%v, to:%v, txid:%v", chain, log.Currency, from, to, tx.Hash().Hex())

		log.Status = string(global.CollectLogStatusSuccess)
		log.GasUsed = receipt.GasUsed
		log.EffectiveGasPrice = receipt.EffectiveGasPrice.Int64()
		log.GasPrice = tx.GasPrice().Int64()

		switch chain {
		case global.ChainNameBsc:
			log.TotalGasFee = int64(log.GasUsed) * log.GasPrice
		case global.ChainNameEth:
			log.TotalGasFee = int64(log.GasUsed) * log.EffectiveGasPrice
		}
	} else {
		logx.Errorf("evm collect transfer failed, chain:%v, currency:%v, from:%v, to:%v, txid:%v, status:%d", chain, log.Currency, from, to, tx.Hash().Hex(), receipt.Status)

		log.Status = string(global.CollectLogStatusFailed)
		log.Description = fmt.Sprintf("evm collect log transfer failed, receipt status:%v", receipt.Status)
	}

	return
}

func (s *ChainService) evmTransfer(cli *ethclient.Client, chain global.ChainName, to, prikey, caddr string, amount, fmax float64) (tx *types.Transaction, err error) {
	ctx := context.Background()
	privateKey, err := crypto.HexToECDSA(prikey)
	if err != nil {
		return
	}

	publicKeyECDSA := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	uaddr := common.HexToAddress(caddr)
	toAddress := common.HexToAddress(to)

	sun := int64(0)
	switch chain {
	case global.ChainNameBsc:
		sun = global.Sun(amount, global.AmountTypo18e)
	case global.ChainNameEth:
		sun = global.Sun(amount, global.AmountTypo6e)
	}

	efee, tcap, fcap, glimit, err := erc20.EstimateTransactionFee(ctx, cli, chain, uaddr, fromAddress, toAddress, sun)
	if err != nil {
		return
	}

	mwei := new(big.Float).Mul(new(big.Float).SetFloat64(fmax), new(big.Float).SetInt64(1e18))
	mfee := new(big.Int)
	mwei.Int(mfee)

	if efee.Cmp(mfee) > 0 {
		err = fmt.Errorf("evm collect fee limit overflow, need:%v, actual:%v, fmax:%v", efee, mfee, fmax)
		return
	}

	cid, err := cli.ChainID(ctx)
	if err != nil {
		return
	}

	// tx = clearPendingTransaction(ctx, cli, cid, fromAddress, privateKey)
	// return

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

func clearPendingTransaction(ctx context.Context, cli *ethclient.Client, cid *big.Int, addr common.Address, prikey *ecdsa.PrivateKey) (stx *types.Transaction) {
	confirmedNonce, err := cli.NonceAt(ctx, addr, nil)
	if err != nil {
		logx.Errorf("获取确认Nonce失败: %v, cnonce:%v", err, confirmedNonce)
		return
	}

	gasTipCap := big.NewInt(10_000_000_000) // 10 Gwei 矿工小费
	gasFeeCap := big.NewInt(30_000_000_000) // 30 Gwei 总费用上限
	gasLimit := uint64(21000)               // 原生转账标准GasLimit

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
	stx, err = types.SignTx(tx, signer, prikey)
	if err != nil {
		logx.Errorf("交易签名失败: %v", err)
		return
	}

	// 6. 发送覆盖交易
	err = cli.SendTransaction(ctx, stx)
	if err != nil {
		logx.Errorf("发送覆盖交易失败: %v", err)
		return
	}

	logx.Infof("clean pending transaction, addr:%v, txid:%v", addr.Hex(), stx.Hash().Hex())

	return
}

func (s *ChainService) sendTronTransaction(log *entity.AddressFundCollectLog, chain global.ChainName, from, to, prikey, caddr string, amount, fmax float64) (err error) {
	cli, ok := s.cli.(*client.Client)
	if !ok {
		logx.Errorf("send transaction convert cli failed, chain:%v, log:%v", chain, log.Id)
		return
	}

	ctx := context.Background()

	ffrom := ttypes.MustNewAddressFromBase58(from)
	tto := ttypes.MustNewAddressFromBase58(to)
	uaddr := ttypes.MustNewAddressFromBase58(caddr)

	// cli.SmartContract().EstimateEnergy()

	ustx, err := cli.TRC20(uaddr).Transfer(ctx, ffrom, tto, decimal.NewFromFloat(amount))
	if err != nil {
		return
	}

	signer, err := signer.NewPrivateKeySigner(prikey)
	if err != nil {
		return
	}

	result, err := cli.SignAndBroadcast(ctx, ustx, client.DefaultBroadcastOptions(), signer)
	if err != nil {
		return
	}

	log.TransactionId = result.TxID

	if result.Success && result.Code == api.Return_SUCCESS {
		logx.Infof("🎉 tron collect transfer success, chain:%v, currency:%v, from:%v, to:%v, txid:%v", global.ChainNameSolana, log.Currency, from, to, result.TxID)

		log.Status = string(global.CollectLogStatusSuccess)
		log.GasUsed = uint64(result.EnergyUsage)
		log.EffectiveGasPrice = result.NetUsage
	} else {
		logx.Errorf("tron collect transfer failed, chain:%v, currency:%v, from:%v, to:%v, txid:%v, status:%d", chain, log.Currency, from, to, result.TxID, result.Code)

		log.Status = string(global.CollectLogStatusFailed)
		log.Description = fmt.Sprintf("tron collect transfer failed, status:%v, msg:%v", result.Code, result.Message)
	}

	return
}

func (s *ChainService) sendSolanaTransaction(log *entity.AddressFundCollectLog, chain global.ChainName, from, to, prikey, caddr string, amount, fmax float64) (err error) {
	cli, ok := s.cli.(*rpc.Client)
	if !ok {
		err = fmt.Errorf("send transaction convert cli failed, chain:%v, log:%v", chain, log.Id)
		return
	}
	txid, statuses, detail, err := s.solanaTransfer(cli, caddr, prikey, from, to, amount, fmax)
	if err != nil {
		logx.Errorf("solona send transaction failed, err:%v", err)

		log.Status = string(global.CollectLogStatusFailed)
		emsg := err.Error()
		if len(emsg) > 1022 {
			log.Description = emsg[0:1022]
		} else {
			log.Description = emsg
		}

		err = biz.AddressFundCollectSendTxFailed
		return
	}

	log.TransactionId = txid

	if len(statuses.Value) > 0 && (statuses.Value[0].ConfirmationStatus == rpc.ConfirmationStatusConfirmed || statuses.Value[0].ConfirmationStatus == rpc.ConfirmationStatusFinalized) {
		logx.Infof("🎉 solana collect transfer success, chain:%v, currency:%v, from:%v, to:%v, txid:%v", global.ChainNameTron, log.Currency, from, to, txid)

		log.Status = string(global.CollectLogStatusSuccess)
		log.TotalGasFee = int64(detail.Meta.Fee)
	} else {
		logx.Errorf("solana collect transfer failed, chain:%v, currency:%v, from:%v, to:%v, txid:%v", chain, log.Currency, from, to, txid)

		log.Status = string(global.CollectLogStatusFailed)
		log.Description = fmt.Sprintf("solana collect transfer failed, currency:%v, from:%v, to:%v", log.Currency, from, to)
	}

	return
}

func (s *ChainService) solanaTransfer(cli *rpc.Client, caddr, prikey, from, to string, amount, fmax float64) (txid string, statuses *rpc.GetSignatureStatusesResult, detail *rpc.GetTransactionResult, err error) {
	ctx := context.Background()

	pkey, err := solana.PrivateKeyFromBase58(prikey)
	if err != nil {
		return
	}
	ffrom := pkey.PublicKey()
	tto, err := solana.PublicKeyFromBase58(to)
	if err != nil {
		return
	}

	umint, err := solana.PublicKeyFromBase58(caddr)
	if err != nil {
		return
	}

	fromt, err := getTokenAccountAddress(ffrom, umint)
	if err != nil {
		return
	}
	tot, err := getTokenAccountAddress(tto, umint)
	if err != nil {
		return
	}

	fbal, err := checkTokenAccountBalance(ctx, cli, fromt)
	if err != nil {
		return
	}
	if fbal < amount {
		err = fmt.Errorf("solana account token balance insufficient, from:%v, now:%v, need:%v, caddr:%v", from, fbal, amount, caddr)
		return
	}

	lblock, err := cli.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return
	}

	ins, err := token.NewTransferInstruction(uint64(global.Sun(amount, global.AmountTypo6e)), fromt, tot, ffrom, nil).ValidateAndBuild()
	if err != nil {
		return
	}

	tx, err := solana.NewTransaction([]solana.Instruction{ins}, lblock.Value.Blockhash, solana.TransactionPayer(ffrom))
	if err != nil {
		return
	}

	serialized, err := tx.Message.MarshalBinary()
	if err != nil {
		return
	}

	bmsg := base64.StdEncoding.EncodeToString(serialized)
	fee, err := cli.GetFeeForMessage(ctx, bmsg, rpc.CommitmentFinalized)
	if err != nil {
		return
	}
	need := global.Amount(int64(*fee.Value), global.AmountTypo9e)
	if fmax < need {
		err = fmt.Errorf("solana collect fee limit overflow, need:%v, actual:%v, fmax:%v", need, fee.Value, fmax)
		return
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(ffrom) {
			return &pkey
		}
		return nil
	})
	if err != nil {
		return
	}

	wscli, err := ws.Connect(context.Background(), s.cfg.Solana.WssNetwork2)
	if err != nil {
		return
	}
	sig, err := confirm.SendAndConfirmTransaction(context.TODO(), cli, wscli, tx)
	if err != nil {
		return
	}

	txid = sig.String()

	statuses, err = cli.GetSignatureStatuses(ctx, false, []solana.Signature{sig}...)
	if err != nil {
		return
	}

	detail, err = cli.GetTransaction(ctx, sig, &rpc.GetTransactionOpts{Commitment: rpc.CommitmentConfirmed, Encoding: solana.EncodingBase64})
	if err != nil {
		return
	}

	return
}

func checkTokenAccountBalance(ctx context.Context, rpcClient *rpc.Client, owner solana.PublicKey) (balance float64, err error) {
	tbal, err := rpcClient.GetTokenAccountBalance(ctx, owner, rpc.CommitmentFinalized)
	if err != nil {
		return
	}

	balance = *tbal.Value.UiAmount
	return
}

func getTokenAccountAddress(wallet solana.PublicKey, mint solana.PublicKey) (solana.PublicKey, error) {
	tokenAccount, _, err := solana.FindAssociatedTokenAddress(wallet, mint)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return tokenAccount, nil
}
