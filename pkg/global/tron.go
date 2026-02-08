package global

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
)

func BuildTronTransferParam(to string, amount *big.Int) string {
	methodID := crypto.Keccak256([]byte("transfer(address,uint256)"))[:4]

	toHex := Base58ToHex(to)[2:] // 去掉 41
	paddedTo := fmt.Sprintf("%064s", toHex)
	amountHex := fmt.Sprintf("%064x", amount)

	return hex.EncodeToString(methodID) + paddedTo + amountHex
}

func Base58ToHex(addr string) string {
	decoded := base58.Decode(addr)
	return hex.EncodeToString(decoded[:len(decoded)-4])
}
