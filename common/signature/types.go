package signature

import (
	"crypto/ecdsa"
	"github.com/FlareZone/melon-backend/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"math/big"
)

func RSV(sign hexutil.Hex) (r, s *big.Int, v uint) {
	r = new(big.Int).SetBytes(sign[:32])
	s = new(big.Int).SetBytes(sign[32:64])
	v = uint(sign[64])

	// 根据以太坊的惯例，如果 v 不是 27 或 28，需要调整
	if v < 27 {
		v += 27
	}
	return
}

func GetPrivateKeyFromMnemonic(mnemonic string) *ecdsa.PrivateKey {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, _ := bip32.NewMasterKey(seed)
	purposeKey, _ := masterKey.NewChildKey(0x8000002C)
	coinTypeKey, _ := purposeKey.NewChildKey(0x8000003C)
	accountKey, _ := coinTypeKey.NewChildKey(0x80000000)
	changeKey, _ := accountKey.NewChildKey(0)
	addressKey, _ := changeKey.NewChildKey(0)

	// 使用addressKey的 PrivateKey 方法
	privateKey, _ := crypto.ToECDSA(addressKey.Key)
	return privateKey
}

func GetPrivateKeyFromPrivateHex(hex hexutil.Hex) *ecdsa.PrivateKey {
	privateKey, _ := crypto.ToECDSA(hex)
	return privateKey
}
