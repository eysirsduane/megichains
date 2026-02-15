package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"megichains/pkg/biz"
	"megichains/pkg/converter"
	"megichains/pkg/crypt"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"os"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gagliardetto/solana-go"
	"github.com/jinzhu/copier"
	"github.com/mr-tron/base58/base58"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/sha3"
	"gorm.io/gorm"
)

type AddressService struct {
	db *gorm.DB
}

func NewAddressService(db *gorm.DB) *AddressService {
	return &AddressService{db: db}
}

func (s *AddressService) Get(ctx context.Context, id int64) (resp *converter.AddressItem, err error) {
	addr, err := gorm.G[entity.Address](s.db).Where("id = ?", id).First(ctx)
	if err != nil {
		logx.Errorf("address detail get failed, id:%v, err:%v", id, err)
		return
	}

	resp = &converter.AddressItem{}
	copier.Copy(resp, addr)

	return
}

func (s *AddressService) FirstFree(ctx context.Context) (resp *converter.AddressItem, err error) {
	now := time.Now().Add(time.Minute * -31).UnixMilli()
	addr, err := gorm.G[entity.Address](s.db).Where("last_used_timestamp <= ? or last_used_timestamp is null", now).First(ctx)
	if err != nil {
		logx.Errorf("address detail get failed, err:%v", err)
		return
	}

	resp = &converter.AddressItem{}
	copier.Copy(resp, addr)

	return
}

func (s *AddressService) Save(ctx context.Context, req *converter.AddressItem) (err error) {
	addr := &entity.Address{}
	copier.Copy(&addr, req)

	if req.Id > 0 {
		_, err = gorm.G[entity.Address](s.db).Updates(ctx, *addr)
		if err != nil {
			logx.Errorf("address update failed, id:%v, err:%v", req.Id, err)
			err = biz.AddressSaveFailed
			return
		}
	} else {
		err = gorm.G[entity.Address](s.db).Create(ctx, addr)
		if err != nil {
			logx.Errorf("address create failed, id:%v, err:%v", req.Id, err)
			err = biz.AddressCreateFailed
			return
		}
	}

	return
}

func (s *AddressService) ChangeStatus(addr string, status global.AddressStatus) (err error) {
	err = s.db.Model(&entity.Address{}).Where("address = ?", addr).Update("status", status).Error
	if err != nil {
		return
	}

	return
}

func (s *AddressService) SetLastUsed(addr string) (err error) {
	err = s.db.Model(&entity.Address{}).Where("address = ?", addr).Update("last_used_timestamp", time.Now().UnixMilli()).Error
	if err != nil {
		return
	}

	return
}

func (s *AddressService) GroupAll(ctx context.Context) (resp *converter.RespConverter[entity.AddressGroup], err error) {
	all, err := gorm.G[entity.AddressGroup](s.db).Order("id asc").Find(ctx)
	if err != nil {
		logx.Errorf("db address group all failed, err:%v", err)
		err = biz.AddressGroupFindFailed
		return
	}

	resp = converter.ConvertToRecordsResp(all, 0, 0, 0)

	return
}

func (s *AddressService) GroupFind(ctx context.Context, req *converter.AddressGroupListReq) (resp *converter.RespConverter[entity.AddressGroup], err error) {
	db := gorm.G[entity.AddressGroup](s.db).Order("id asc")
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	items, err := db.Offset(global.Offset(req.Current, req.Size)).Limit(req.Size).Find(ctx)
	if err != nil {
		logx.Errorf("address group paging failed, err:%v", err)
		err = biz.AddressGroupFindFailed
		return
	}

	total, err := db.Count(ctx, "id")
	if err != nil {
		logx.Errorf("address group count failed, err:%v", err)
		err = biz.AddressCountFailed
		return
	}

	resp = converter.ConvertToPagingRecordsResp(items, req.Current, req.Size, total)

	return
}

func (s *AddressService) Find(ctx context.Context, req *converter.AddressListReq) (resp *converter.RespConverter[converter.AddressWithGroup], err error) {
	db := s.db.Model(&entity.Address{}).Order("typo")
	if req.Address != "" {
		db = db.Where("address = ?", req.Address)
	}
	if req.Address2 != "" {
		db = db.Where("address2 = ?", req.Address2)
	}
	if req.Chain != "" {
		db = db.Where("chain = ?", req.Chain)
	}
	if req.Status != "" {
		db = db.Where("addresses.status = ?", req.Status)
	}
	if req.Typo != "" {
		db = db.Where("typo = ?", req.Typo)
	}
	if req.GroupId > 0 {
		db = db.Where("group_id = ?", req.GroupId)
	}
	if req.Start > 0 {
		db = db.Where("created_at >= ?", req.Start)
	}
	if req.End > 0 {
		db = db.Where("created_at <= ?", req.End)
	}
	addrs := make([]*entity.Address, 0, req.Size)
	err = db.Session(&gorm.Session{}).Preload("AddressBalance").Offset(global.Offset(req.Current, req.Size)).Limit(req.Size).Find(&addrs).Error
	if err != nil {
		logx.Errorf("address list find failed, err:%v", err)
		err = biz.AddressFindFailed
		return
	}
	total := int64(0)
	err = db.Session(&gorm.Session{}).Count(&total).Error
	if err != nil {
		logx.Errorf("address list count failed, err:%v", err)
		err = biz.AddressCountFailed
		return
	}

	gids := make([]int64, 0, len(addrs))
	for _, addr := range addrs {
		gids = append(gids, addr.GroupId)
	}
	groups := make([]*entity.AddressGroup, 0)
	err = s.db.Model(&entity.AddressGroup{}).Where("id in (?)", gids).Find(&groups).Error
	if err != nil {
		logx.Errorf("address list find groups failed, err:%v", err)
		err = biz.AddressFindGroupFailed
		return
	}

	items := make([]converter.AddressWithGroup, 0, len(addrs))
	for _, addr := range addrs {
		aitem := converter.AddressItem{}
		copier.Copy(&aitem, addr)

		gname := ""
		for _, group := range groups {
			if addr.GroupId == group.Id {
				gname = group.Name
				break
			}
		}

		items = append(items, converter.AddressWithGroup{AddressItem: aitem, GroupName: gname})
	}

	resp = converter.ConvertToPagingRecordsResp(items, req.Current, req.Size, total)

	return
}

func (s *AddressService) UseAddress(id int64) (addr *entity.Address, err error) {
	err = s.db.Where("status = ?", global.AddressStatusInFree).First(&addr).Error
	if err != nil {
		logx.Errorf("db use address failed, id:%v, err:%v", id, err)
		return
	}

	return
}

func (s *AddressService) CreateAddress(req *converter.ChainAddressCreateReq) (err error) {
	switch req.Chain {
	case "EVM":
		for i := range req.Count {
			fmt.Println(i)
			s.createEvmAddress(req.Chain, req.GroupId)
		}
	case "TRON":
		for i := range req.Count {
			fmt.Println(i)
			s.createTronAddress(req.Chain, req.GroupId)
		}
	case "SOLANA":
		for i := range req.Count {
			fmt.Println(i)
			s.createSolanaAddresses(req.Chain, req.GroupId)
		}
	default:
		err = biz.AddressCreateFailed
	}

	return
}

func (s *AddressService) createEvmAddress(chain string, gid int64) {
	// 1. 生成私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		logx.Errorf("address create gen key failed, chain:%v, err:%v", chain, err)
		return
	}

	// 2. 私钥转字节 / hex
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hexutil.Encode(privateKeyBytes)
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	// 3. 生成公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logx.Errorf("address create public key gen failed, chain:%v, err:%v", chain, err)
		return
	}

	// 4. 生成地址（ETH / BSC 通用）
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	enc, err := crypt.Encrypt(privateKeyHex, crypt.PrivateKeySecretPassword, crypt.PrivateKeySecretSalt)
	if err != nil {
		logx.Errorf("create address failed, chain:%v, err:%v", chain, err)
		return
	}
	addr := &entity.Address{
		Chain:          chain,
		GroupId:        gid,
		Typo:           string(global.AddressTypoIn),
		Status:         string(global.AddressStatusInFree),
		Address:        address,
		PrivateKey:     enc,
		AddressBalance: entity.AddressBalance{Address: address, BscUsdt: 10, BscUsdc: 10, EthUsdt: 10, EthUsdc: 10},
	}

	err = s.db.Create(addr).Error
	if err != nil {
		panic(err)
	}
}

func (s *AddressService) createTronAddress(chain string, gid int64) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		panic(err)
	}

	pubKey := privateKey.PubKey().SerializeUncompressed()

	// Keccak256(pubKey[1:])
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubKey[1:])
	addressHash := hash.Sum(nil)

	// 取后20字节，加 TRON 前缀 0x41
	address := append([]byte{0x41}, addressHash[12:]...)

	// Base58Check
	checksum := checksum(address)
	fullAddress := append(address, checksum...)

	base58Address := base58.Encode(fullAddress)

	enc, err := crypt.Encrypt(fmt.Sprintf("%x", privateKey.Serialize()), crypt.PrivateKeySecretPassword, crypt.PrivateKeySecretSalt)
	if err != nil {
		logx.Errorf("address create failed, chain:%v, err:%v", chain, err)
		return
	}
	addr := &entity.Address{
		Chain:          chain,
		GroupId:        gid,
		Typo:           string(global.AddressTypoIn),
		Status:         string(global.AddressStatusInFree),
		Address:        base58Address,
		Address2:       fmt.Sprintf("%x", address),
		PrivateKey:     enc,
		PublicKey:      fmt.Sprintf("%x", pubKey),
		AddressBalance: entity.AddressBalance{Address: base58Address, TronUsdt: 10, TronUsdc: 10},
	}

	err = s.db.Create(addr).Error
	if err != nil {
		panic(err)
	}
}

func (s *AddressService) createSolanaAddresses(chain string, gid int64) {
	keypair := solana.NewWallet()
	data, _ := json.MarshalIndent(keypair, "", "  ")
	fileName := fmt.Sprintf("wallet_%d.json", 1)
	_ = os.WriteFile(fileName, data, 0644)

	prikey := keypair.PrivateKey.String()
	encrypted, err := crypt.Encrypt(prikey, crypt.PrivateKeySecretPassword, crypt.PrivateKeySecretSalt)
	if err != nil {
		fmt.Printf("address create failed, chain:%v, err:%v", chain, err)
		return
	}

	addr := &entity.Address{
		Chain:      chain,
		GroupId:    gid,
		Typo:       string(global.AddressTypoIn),
		Status:     string(global.AddressStatusInFree),
		Address:    keypair.PublicKey().String(),
		Address2:   "",
		PrivateKey: encrypted,
		PublicKey:  "",
		AddressBalance: entity.AddressBalance{
			Address:    keypair.PublicKey().String(),
			SolanaUsdt: 10,
			SolanaUsdc: 10,
		},
	}

	err = s.db.Create(addr).Error
	if err != nil {
		panic(err)
	}
}

func checksum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	return second[:4]
}

func (s *AddressService) GroupGet(ctx context.Context, id int64) (resp *converter.AddressGroupItem, err error) {
	group, err := gorm.G[entity.AddressGroup](s.db).Where("id = ?", id).First(ctx)
	if err != nil {
		logx.Errorf("address group get detail failed, err:%v", err)
		err = biz.AddressGroupFindFailed
		return
	}

	resp = &converter.AddressGroupItem{}
	copier.Copy(resp, group)

	return
}

func (s *AddressService) GroupSave(ctx context.Context, req *converter.AddressGroupItem) (err error) {
	group := &entity.AddressGroup{}
	copier.Copy(group, req)

	if req.Name == "" || req.Status == "" {
		err = biz.AddressGroupFieldInvalid
		return
	}

	if req.Id > 0 {
		_, err = gorm.G[entity.AddressGroup](s.db).Where("id = ?", req.Id).Updates(ctx, *group)
		if err != nil {
			logx.Errorf("address group save failed, err:%v", err)
			err = biz.AddressGroupSaveFailed
			return
		}
	} else {
		err = gorm.G[entity.AddressGroup](s.db).Create(ctx, group)
		if err != nil {
			logx.Errorf("address group create failed, err:%v", err)
			err = biz.AddressGroupSaveFailed
			return
		}
	}

	return
}

func (s *AddressService) GroupCreate(ctx context.Context, req *converter.AddressGroupItem) (err error) {
	group := &entity.AddressGroup{}
	copier.Copy(group, req)

	err = gorm.G[entity.AddressGroup](s.db).Create(ctx, group)
	if err != nil {
		logx.Errorf("address group save failed, err:%v", err)
		err = biz.AddressGroupSaveFailed
		return
	}

	return
}
