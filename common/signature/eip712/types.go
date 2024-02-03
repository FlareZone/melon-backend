package eip712

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/FlareZone/melon-backend/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type Eip712 struct {
	domain    *apitypes.TypedDataDomain
	typedData *apitypes.TypedData
}

func (e *Eip712) withTypeDataDomain(name string, version string, chainId *math.HexOrDecimal256,
	verifyingContract string) {
	e.domain = &apitypes.TypedDataDomain{
		Name:              name,
		Version:           version,
		ChainId:           chainId,
		VerifyingContract: verifyingContract,
	}
}

func (e *Eip712) Sign(privateKey *ecdsa.PrivateKey) (string, error) {
	hash, _, err := typeDataAndHash(*e.typedData)
	fmt.Println("hash: ", hexutil.Hex(hash).Hex())
	if err != nil {
		return "", err
	}
	sig, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return "", err
	}
	if sig[64] < 27 {
		sig[64] += 27
	}
	return fmt.Sprintf("0x%s", hex.EncodeToString(sig)), nil
}

func typeDataAndHash(typedData apitypes.TypedData) ([]byte, string, error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, "", err
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, "", err
	}
	hashStruct, err := HashStruct(typedData)
	if err != nil {
		return nil, "", err
	}
	fmt.Println("typedDataHash ", hexutil.Hex(hashStruct).Hex())
	rawData := fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash))
	return crypto.Keccak256([]byte(rawData)), rawData, nil
}

func HashStruct(typedData apitypes.TypedData) ([]byte, error) {
	fmt.Println("primaryTypeHash: ", hexutil.Hex(typedData.TypeHash(typedData.PrimaryType)).Hex())

	encodedData, err := typedData.EncodeData(typedData.PrimaryType, typedData.Message, 2)
	if err != nil {
		return nil, err
	}
	fmt.Println("encode data: ", hexutil.Hex(encodedData))
	return crypto.Keccak256(encodedData), nil
}
