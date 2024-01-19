package hexutil

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"strconv"
	"strings"
)

func Wei2Ether(balance *big.Int) *big.Float {
	return big.NewFloat(0).Quo(big.NewFloat(0).SetInt(balance), big.NewFloat(params.Ether))
}

func BytesToHex(b []byte) string {
	return "0x" + strings.ToLower(hex.EncodeToString(b))
}

func HexStringToUInt64(input string) uint64 {
	var raw string
	if strings.HasPrefix(input, "0x") {
		raw = input[2:]
	} else {
		raw = input
	}
	s, err := strconv.ParseUint(raw, 16, 64)
	if err != nil {
		return 0
	} else {
		return s
	}
}

func HexToBigInt(hex string) *big.Int {
	if !(strings.HasPrefix(hex, "0x") || strings.HasPrefix(hex, "0X")) {
		return big.NewInt(0)
	}
	return StringToBigInt(hex, 16)
}

func RemoveHexLeadingZeroDigits(hex string) string {
	if len(hex) > 2 && (strings.HasPrefix(hex, "0x") || strings.HasPrefix(hex, "0X")) {
		if strings.TrimLeft(hex[2:], "0") == "" {
			return "0x0"
		}
		return "0x" + strings.TrimLeft(hex[2:], "0")
	}
	return hex
}
