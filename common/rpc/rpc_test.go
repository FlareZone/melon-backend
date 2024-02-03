package rpc

import (
	"context"
	"github.com/FlareZone/melon-backend/common/hexutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvmRpc_ChainID(t *testing.T) {
	rpc, err := NewRpc(context.Background(), "http://192.168.0.6:8545")
	assert.NoError(t, err)
	id, err := rpc.ChainID()
	assert.NoError(t, err)
	t.Log("chain_id:", id)
}

func TestEvmRpc_Nonce(t *testing.T) {
	rpc, err := NewRpc(context.Background(), "http://192.168.0.6:8545")
	assert.NoError(t, err)
	nonce, err := rpc.Nonce(context.Background(), "0x5AD07396DCe3a6aA66eeAEA2f2207D81e35Ba48a", "pending")
	assert.NoError(t, err)
	t.Log(nonce)
}

func TestEvmRpc_BalanceAt(t *testing.T) {
	rpc, err := NewRpc(context.Background(), "http://192.168.0.6:8545")
	assert.NoError(t, err)
	balance, err := rpc.BalanceAt(context.Background(), hexutil.StringToHex("0x5AD07396DCe3a6aA66eeAEA2f2207D81e35Ba48a"))
	assert.NoError(t, err)
	t.Log(balance)
}

func TestEvmRpc_BalanceOf(t *testing.T) {
	rpc, err := NewRpc(context.Background(), "http://192.168.0.6:8545")
	assert.NoError(t, err)
	of, err := rpc.BalanceOf(context.Background(), hexutil.StringToHex("0x87A31aDa5ffDa118D83543CEF9EDCCf9806a0f5f"),
		hexutil.StringToHex("0xcB796Ee82d8824A0A0C8A14E177a99323Cc5A814"))
	assert.NoError(t, err)
	t.Log(of)
}

func TestEvmRpc_TokenName(t *testing.T) {
	rpc, err := NewRpc(context.Background(), "http://192.168.0.6:8545")
	assert.NoError(t, err)
	result, err := rpc.TokenName(context.Background(), hexutil.StringToHex("0x87A31aDa5ffDa118D83543CEF9EDCCf9806a0f5f"))
	assert.NoError(t, err)
	t.Log(result)
}

func TestEvmRpc_TokenSymbol(t *testing.T) {
	rpc, err := NewRpc(context.Background(), "http://192.168.0.6:8545")
	assert.NoError(t, err)
	result, err := rpc.TokenSymbol(context.Background(), hexutil.StringToHex("0x87A31aDa5ffDa118D83543CEF9EDCCf9806a0f5f"))
	assert.NoError(t, err)
	t.Log(result)
}

func TestEvmRpc_TokenDecimals(t *testing.T) {
	rpc, err := NewRpc(context.Background(), "http://192.168.0.6:8545")
	assert.NoError(t, err)
	result, err := rpc.TokenDecimal(context.Background(), hexutil.StringToHex("0x87A31aDa5ffDa118D83543CEF9EDCCf9806a0f5f"))
	assert.NoError(t, err)
	t.Log(result)
}
