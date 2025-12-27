package service

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"megichains/pkg/biz"
	"megichains/pkg/converter"
	"megichains/pkg/entity"
	"megichains/pkg/global"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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

func (s *AddressService) GetAddress(id int64) (addr *entity.Address, err error) {
	err = s.db.Where("id = ?", id).First(&addr).Error
	if err != nil {
		logx.Errorf("db get address failed, id:%v, err:%v", id, err)
		return
	}

	return
}

func (s *AddressService) UseAddress(id int64) (addr *entity.Address, err error) {
	err = s.db.Where("status = ?", global.AddressTypoVacant).First(&addr).Error
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
			s.createEvmAddress(req.Chain)
		}
	case "TRON":
		for i := range req.Count {
			fmt.Println(i)
			s.createTronAddress(req.Chain)
		}
	default:
		err = biz.AddressCreateFailed
	}

	return
}

func (s *AddressService) createEvmAddress(chain string) {
	// 1. 生成私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	// 2. 私钥转字节 / hex
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hexutil.Encode(privateKeyBytes)

	// 3. 生成公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("cannot assert public key type")
	}

	// 4. 生成地址（ETH / BSC 通用）
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	fmt.Println("Private Key:", privateKeyHex)
	fmt.Println("BSC Address:", address)

	addr := &entity.Address{
		Chain:      chain,
		Typo:       string(global.AddressTypoIn),
		Status:     string(global.AddressTypoVacant),
		Address:    address,
		PrivateKey: privateKeyHex,
	}

	err = s.db.Create(addr).Error
	if err != nil {
		panic(err)
	}
}
func (s *AddressService) createTronAddress(chain string) {
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

	// fmt.Println("Private Key:", fmt.Sprintf("%x", privateKey.Serialize()))
	// fmt.Println("Public Key :", fmt.Sprintf("%x", pubKey))
	// fmt.Println("TRON Addr  :", base58Address)
	addr := &entity.Address{
		Chain:      chain,
		Typo:       string(global.AddressTypoIn),
		Status:     string(global.AddressTypoVacant),
		Address:    base58Address,
		Address2:   fmt.Sprintf("%x", address),
		PrivateKey: fmt.Sprintf("%x", privateKey.Serialize()),
		PublicKey:  fmt.Sprintf("%x", pubKey),
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
