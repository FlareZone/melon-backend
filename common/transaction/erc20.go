package transaction

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/FlareZone/melon-backend/common/hexutil"
	"github.com/FlareZone/melon-backend/common/rpc"
	"github.com/FlareZone/melon-backend/common/signature"
	"github.com/FlareZone/melon-backend/common/signature/eip712"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/log"
	"github.com/umbracle/ethgo/abi"
	"math/big"
)

type ERC20 struct {
	contractAddress hexutil.Hex
	rpc             *rpc.EvmRpc
	privateKey      *ecdsa.PrivateKey
	eip712          *eip712.PermitToken
	TxHash          string
	name            string
	symbol          string
	decimals        uint64
}

func (e *ERC20) WithContractAddress(address hexutil.Hex) *ERC20 {
	e.contractAddress = address
	return e
}

func (e *ERC20) WithPrivateKey(private *ecdsa.PrivateKey) *ERC20 {
	e.privateKey = private
	return e
}

func (e *ERC20) WithRpc(rpc *rpc.EvmRpc) *ERC20 {
	e.rpc = rpc
	return e
}

func (e *ERC20) WithEip712(permit *eip712.PermitToken) *ERC20 {
	e.eip712 = permit
	return e
}

func (e *ERC20) Name() string {
	if e.name != "" {
		return e.name
	}
	result, err := e.rpc.TokenName(context.Background(), e.contractAddress)
	if err != nil {
		log.Error("get token name fail", "address", e.contractAddress, "err", err)
	}
	e.name = result
	return e.name
}

func (e *ERC20) Symbol() string {
	if e.symbol != "" {
		return e.symbol
	}
	result, err := e.rpc.TokenSymbol(context.Background(), e.contractAddress)
	if err != nil {
		log.Error("get token symbol fail", "address", e.contractAddress, "err", err)
	}
	e.symbol = result
	return e.symbol
}

func (e *ERC20) Decimal() uint64 {
	if e.decimals > 0 {
		return e.decimals
	}
	decimals, err := e.rpc.TokenDecimal(context.Background(), e.contractAddress)
	if err != nil {
		log.Error("get token decimal fail", "address", e.contractAddress, "err", err)
	}
	e.decimals = decimals
	return e.decimals
}

func (e *ERC20) BalanceOf(address hexutil.Hex) *big.Int {
	of, err := e.rpc.BalanceOf(context.Background(), e.contractAddress, address)
	if err != nil {
		log.Error("balance of fail", "address", address, "token", e.contractAddress, "err", err)
		return new(big.Int).SetInt64(0)
	}
	return of
}

func (e *ERC20) QuoDecimals(address hexutil.Hex) *big.Float {
	of, err := e.rpc.BalanceOf(context.Background(), e.contractAddress, address)
	if err != nil {
		log.Error("QuoDecimals, query balanceOf  fail", "address", address, "token", e.contractAddress, "err", err)
		return new(big.Float).SetInt64(0)
	}
	decimal := e.Decimal()
	if decimal <= 0 {
		log.Error("QuoDecimals, query decimal  fail", "address", address, "token", e.contractAddress, "err", err)
		return new(big.Float).SetInt64(0)
	}
	return new(big.Float).Quo(new(big.Float).SetInt(of), new(big.Float).SetInt(
		new(big.Int).Exp(new(big.Int).SetUint64(10), big.NewInt(int64(e.Decimal())), nil)))
}

func (e *ERC20) Allowance(owner, spender string) (*big.Int, error) {
	method, err := abi.NewMethod("function allowance(address owner, address spender)")
	if err != nil {
		return nil, err
	}
	encode, err := method.Encode([]interface{}{owner, spender})
	if err != nil {
		return nil, err
	}

	to := common.HexToAddress(e.contractAddress.Hex())
	param := rpc.CallParameter{
		To:   &to,
		Data: encode,
	}
	var result math.HexOrDecimal256
	err = e.rpc.CallContext(context.Background(), &result, "eth_call", param.ToArg(), "latest")
	if err != nil {
		return nil, err
	}
	balance := big.Int(result)
	return &balance, err

}

func (e *ERC20) TokenAmountToBigInt(amount string) (*big.Int, error) {
	amountFloat, b := new(big.Float).SetString(amount)
	if !b {
		return nil, fmt.Errorf("parse amount: %s  to big.int fail", amount)
	}
	mul := new(big.Float).Mul(amountFloat,
		new(big.Float).SetInt(new(big.Int).Exp(new(big.Int).SetUint64(10), big.NewInt(int64(e.Decimal())), nil)))
	s := mul.String()
	tokenAmount, b := new(big.Int).SetString(s, 10)
	if !b {
		return nil, fmt.Errorf("parse mulAmount: %s to big.int fail", s)
	}
	return tokenAmount, nil
}

func (e *ERC20) Permit(owner, spender string, amount, deadline *big.Int, sign hexutil.Hex) (txHash string, err error) {
	r, s, v := signature.RSV(sign)
	method, err := abi.NewMethod("function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s)")
	if err != nil {
		return
	}

	contractCall := new(ContractCall).WithRpc(e.rpc).
		WithPrivateKey(e.privateKey).
		WithMethod(method, owner, spender, amount, deadline, v, r.Bytes(), s.Bytes())
	err = contractCall.Call(context.Background(), e.contractAddress, big.NewInt(0))
	return contractCall.TxHash, err
}

func (e *ERC20) PermitNonce(owner hexutil.Hex) (uint64, error) {
	nonceAbi, _ := abi.NewMethod("function nonces(address owner)")
	toAddress := common.HexToAddress(e.contractAddress.Hex())
	encode, _ := nonceAbi.Encode([]interface{}{owner.Hex()})
	param := rpc.CallParameter{
		To:   &toAddress,
		Data: encode,
	}
	var result math.HexOrDecimal256
	err := e.rpc.CallContext(context.Background(), &result, "eth_call", param.ToArg(), "latest")
	if err != nil {
		log.Error("nonceAbi is fail", "owner", owner, "err", err)
		return 0, err
	}
	data := big.Int(result)
	return data.Uint64(), nil
}

func (e *ERC20) Transfer(to hexutil.Hex, amount string) (txHash string, err error) {
	tokenAmount, err := e.TokenAmountToBigInt(amount)
	if err != nil {
		return
	}
	method, _ := abi.NewMethod("function transfer(address to, uint256 value)")
	contractCall := new(ContractCall).
		WithRpc(e.rpc).
		WithPrivateKey(e.privateKey).
		WithMethod(method, to, tokenAmount)
	err = contractCall.Call(context.Background(),
		hexutil.StringToHex("0x87A31aDa5ffDa118D83543CEF9EDCCf9806a0f5f"),
		new(big.Int).SetUint64(0))
	return contractCall.TxHash, err
}

func (e *ERC20) TransferFrom(from, to hexutil.Hex, amount string) (txHash string, err error) {
	tokenAmount, err := e.TokenAmountToBigInt(amount)
	if err != nil {
		return
	}
	method, _ := abi.NewMethod("function transferFrom(address from, address to, uint256 value)")
	contractCall := new(ContractCall).WithRpc(e.rpc).
		WithPrivateKey(e.privateKey).WithMethod(method, from, to, tokenAmount)
	err = contractCall.Call(context.Background(), e.contractAddress, new(big.Int).SetUint64(0))
	return contractCall.TxHash, err
}
