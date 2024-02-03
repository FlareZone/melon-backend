package rpc

import (
	hexutil2 "github.com/FlareZone/melon-backend/common/hexutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/umbracle/ethgo/abi"
	"math/big"
)

type NonceStatus string

func (n NonceStatus) String() string {
	return string(n)
}

const (
	Pending  = NonceStatus("pending")
	Latest   = NonceStatus("latest")
	Earliest = NonceStatus("earliest")
)

type CallParameter struct {
	From     common.Address
	To       *common.Address
	Data     []byte
	Gas      uint64
	GasPrice *big.Int
	Value    *big.Int
}

func (c CallParameter) ToArg() interface{} {
	arg := make(map[string]interface{})
	arg["from"] = c.From
	if c.To != nil {
		arg["to"] = c.To
	}
	if c.Data != nil {
		arg["data"] = hexutil2.Hex(c.Data).Hex()
	}
	if c.Value != nil {
		arg["value"] = (*hexutil.Big)(c.Value)
	}
	if c.Gas != 0 {
		arg["gas"] = hexutil.Uint64(c.Gas)
	}
	if c.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(c.GasPrice)
	}
	return arg
}

var (
	balanceOfAbi, _   = abi.NewMethod("function balanceOf(address)")
	balanceTypeAbi, _ = abi.NewType("uint256")

	tokenNameAbi, _     = abi.NewMethod("function name()")
	tokenNameTypeAbi, _ = abi.NewType("tuple(string name)")

	tokenSymbolAbi, _     = abi.NewMethod("function symbol()")
	tokenSymbolTypeAbi, _ = abi.NewType("tuple(string symbol)")

	tokenDecimalAbi, _     = abi.NewMethod("function decimals()")
	tokenDecimalTypeAbi, _ = abi.NewType("tuple(uint256 decimals)")
)

func DecodeString(value string) string {
	var s string
	if len(value) > 64*2+2 {
		s = value[2+128:]
	} else {
		s = value[2:]
	}
	v, _ := hexutil.Decode("0x" + s)
	return string(StripBytes(v))
}

func StripBytes(input []byte) []byte {
	var result []byte
	for _, v := range input {
		if v != 0 {
			result = append(result, v)
		}
	}
	return result
}
