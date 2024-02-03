package transaction

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/FlareZone/melon-backend/common/hexutil"
	"github.com/FlareZone/melon-backend/common/rpc"
	"github.com/ethereum/go-ethereum/common"
	hexutil2 "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/umbracle/ethgo/abi"
	"math/big"
)

type ContractCall struct {
	callData   string
	privateKey *ecdsa.PrivateKey
	rpc        *rpc.EvmRpc
	TxHash     string
}

func (t *ContractCall) WithMethod(method *abi.Method, args ...interface{}) *ContractCall {
	encode, _ := method.Encode(args)
	t.callData = hexutil.BytesToHex(encode)
	return t
}

func (t *ContractCall) WithPrivateKey(private *ecdsa.PrivateKey) *ContractCall {
	t.privateKey = private
	return t
}

func (t *ContractCall) WithRpc(rpc *rpc.EvmRpc) *ContractCall {
	t.rpc = rpc
	return t
}

func (t *ContractCall) address() string {
	return crypto.PubkeyToAddress(t.privateKey.PublicKey).Hex()
}

func (t *ContractCall) Call(ctx context.Context, to hexutil.Hex, amount *big.Int) error {
	nonce, err := t.rpc.Nonce(ctx, t.address(), rpc.Pending)
	if err != nil {
		return fmt.Errorf("getNonce call fail, err is %v", err)
	}
	gasPrice, err := t.rpc.GasPrice(ctx)
	if err != nil {
		return fmt.Errorf("gasPrice call fail, err is %v", err)
	}
	address := common.HexToAddress(to.Hex())
	txData := &types.LegacyTx{
		Nonce:    uint64(nonce),
		GasPrice: gasPrice,
		Value:    amount,
		To:       &address,
		Data:     common.FromHex(t.callData),
	}
	gas, err := t.rpc.EstimateGas(ctx, common.HexToAddress(t.address()), types.NewTx(txData))
	if err != nil {
		return fmt.Errorf("estimateGas fail, err is %v", err)
	}
	txData.Gas = gas.Uint64()
	signedTx, err := t.signTx(txData)
	if err != nil {
		return err
	}
	result, err := t.rpc.SendTx(ctx, signedTx)
	if err != nil {
		return fmt.Errorf("SendTx fail, err is %v", err)
	}
	t.TxHash = result
	return nil
}

func (t *ContractCall) signTx(tx types.TxData) (string, error) {
	chainId, err := t.rpc.ChainID()
	if err != nil {
		return "", fmt.Errorf("chainId is fail, err is %v", err)
	}
	newTx := types.NewTx(tx)
	signedTx, err := types.SignTx(newTx, types.NewLondonSigner(new(big.Int).SetUint64(uint64(chainId))), t.privateKey)
	if err != nil {
		return "", fmt.Errorf("signTx is fail, err is %v", err)
	}
	binary, err := signedTx.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("MarshalBinary signTx is fail, err is %v", err)
	}
	return hexutil2.Encode(binary), nil
}
