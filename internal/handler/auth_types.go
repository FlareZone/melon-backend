package handler

import (
	"encoding/hex"
)

type EthereumEip712SignatureRequest struct {
	TypedData     string `json:"typedData" binding:"required,hexString"`
	TypedDataHash string `json:"typedDataHash" binding:"required,hexString"`
	Signature     string `json:"signature" binding:"required,hexString"`
}

func (e EthereumEip712SignatureRequest) GetTypedData() string {
	typeDataBytes, _ := hex.DecodeString(e.TypedData[2:])
	return string(typeDataBytes)
}

type EthereumEip712SignatureNonceRequest struct {
	EthAddress string `json:"eth_address" binding:"required,hexString"`
}
