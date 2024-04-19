package _type

import (
	"encoding/hex"
)

type ActionType string

func (a ActionType) String() string {
	return string(a)
}

const (
	AuthActionLogin    = ActionType("login")
	AuthActionRegister = ActionType("register")
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

type EmailVerificationCodeRequest struct {
	To string `json:"to" binding:"required,email"`
}

type EmailLoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}
