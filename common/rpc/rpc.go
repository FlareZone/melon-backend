package rpc

import (
	"context"
	"fmt"
	"github.com/FlareZone/melon-backend/common/hexutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/umbracle/ethgo/abi"
	"math/big"
)

type EvmRpc struct {
	*rpc.Client
	chainID math.HexOrDecimal64
}

func NewRpc(ctx context.Context, rawUrl string, opts ...rpc.ClientOption) (*EvmRpc, error) {
	dialClient, err := rpc.DialOptions(ctx, rawUrl, opts...)
	if err != nil {
		return nil, err
	}
	return &EvmRpc{Client: dialClient}, nil
}

func (e *EvmRpc) ChainID() (chainId math.HexOrDecimal64, err error) {
	if e.chainID > 0 {
		return e.chainID, nil
	}
	err = e.Call(&chainId, "eth_chainId")
	e.chainID = chainId
	return
}

func (e *EvmRpc) Nonce(ctx context.Context, address string, status NonceStatus) (nonce math.HexOrDecimal64, err error) {
	err = e.CallContext(ctx, &nonce, "eth_getTransactionCount", address, status.String())
	return
}

func (e *EvmRpc) GasPrice(ctx context.Context) (gasPrice *big.Int, err error) {
	var result math.HexOrDecimal256
	err = e.CallContext(ctx, &result, "eth_gasPrice")
	b := big.Int(result)
	return &b, err
}

func (e *EvmRpc) EstimateGas(ctx context.Context, from common.Address, txData *types.Transaction) (*big.Int, error) {
	parameter := CallParameter{
		From: from,
		Data: txData.Data(),
		To:   txData.To(),
	}
	fmt.Println(parameter.From, "to", parameter.To, "callData", hexutil.Hex(txData.Data()).Hex())
	var result math.HexOrDecimal256
	err := e.CallContext(ctx, &result, "eth_estimateGas", parameter.ToArg())
	if err != nil {
		return nil, err
	}
	data := big.Int(result)
	return &data, nil
}

func (e *EvmRpc) BalanceAt(ctx context.Context, address hexutil.Hex) (*big.Int, error) {
	var result math.HexOrDecimal256
	err := e.CallContext(ctx, &result, "eth_getBalance", address.Hex(), "latest")
	balance := big.Int(result)
	return &balance, err
}

func (e *EvmRpc) BalanceOf(ctx context.Context, to hexutil.Hex, address hexutil.Hex) (*big.Int, error) {
	toAddress := common.HexToAddress(to.Hex())
	encode, _ := balanceOfAbi.Encode([]interface{}{address.Hex()})
	param := CallParameter{
		From: common.HexToAddress(address.Hex()),
		To:   &toAddress,
		Data: encode,
	}
	var result math.HexOrDecimal256
	err := e.CallContext(ctx, &result, "eth_call", param.ToArg(), "latest")
	balance := big.Int(result)
	return &balance, err
}

func (e *EvmRpc) SendTx(ctx context.Context, signTx string) (result string, err error) {
	err = e.CallContext(ctx, &result, "eth_sendRawTransaction", signTx)
	return
}

func (e *EvmRpc) TokenName(ctx context.Context, address hexutil.Hex) (result string, err error) {
	toAddress := common.HexToAddress(address.Hex())
	encode, _ := tokenNameAbi.Encode([]interface{}{})
	param := CallParameter{
		To:   &toAddress,
		Data: encode,
	}
	err = e.CallContext(ctx, &result, "eth_call", param.ToArg(), "latest")
	if err != nil {
		return
	}
	var m map[string]string
	err = tokenNameTypeAbi.DecodeStruct(hexutil.StringToHex(result), &m)
	return m["name"], err
}

func (e *EvmRpc) TokenSymbol(ctx context.Context, address hexutil.Hex) (result string, err error) {
	toAddress := common.HexToAddress(address.Hex())
	encode, _ := tokenSymbolAbi.Encode([]interface{}{})
	param := CallParameter{
		To:   &toAddress,
		Data: encode,
	}
	err = e.CallContext(ctx, &result, "eth_call", param.ToArg(), "latest")
	if err != nil {
		return
	}
	var m map[string]string
	err = tokenSymbolTypeAbi.DecodeStruct(hexutil.StringToHex(result), &m)
	return m["symbol"], err
}

func (e *EvmRpc) TokenDecimal(ctx context.Context, address hexutil.Hex) (decimals uint64, err error) {
	toAddress := common.HexToAddress(address.Hex())
	encode, _ := tokenDecimalAbi.Encode([]interface{}{})
	param := CallParameter{
		To:   &toAddress,
		Data: encode,
	}
	var result string
	err = e.CallContext(ctx, &result, "eth_call", param.ToArg(), "latest")
	if err != nil {
		return
	}
	var m map[string]interface{}
	err = tokenDecimalTypeAbi.DecodeStruct(hexutil.StringToHex(result), &m)
	b := m["decimals"].(*big.Int)
	return b.Uint64(), err
}

func (e *EvmRpc) CallWithCallData(ctx context.Context, result interface{}, to hexutil.Hex, callData hexutil.Hex) (err error) {
	toAddress := common.HexToAddress(to.Hex())
	param := CallParameter{
		To:   &toAddress,
		Data: callData,
	}
	err = e.CallContext(ctx, &result, "eth_call", param.ToArg(), "latest")
	return
}

func (e *EvmRpc) CallWithMethod(ctx context.Context, result interface{}, to hexutil.Hex, method *abi.Method, args ...interface{}) (err error) {
	encode, err := method.Encode(args)
	if err != nil {
		return err
	}
	return e.CallWithCallData(ctx, result, to, encode)
}
