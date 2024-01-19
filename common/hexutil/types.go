package hexutil

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strconv"
	"strings"
)

type BigInt big.Int

func (b BigInt) String() string {
	v := big.Int(b)
	return v.String()
}

func StringToBigInt(hexStr string, base int) *big.Int {
	if strings.HasPrefix(hexStr, "0x") || strings.HasPrefix(hexStr, "0X") {
		hexStr = hexStr[2:]
	}
	var t = new(big.Int)
	t.SetString(hexStr, base)
	return t
}

func (b *BigInt) Scan(value interface{}) error {
	var t big.Int
	switch v := value.(type) {
	case []byte:
		t.SetBytes(v)
		*b = BigInt(t)
		return nil
	case string:
		t.SetString(v, 10)
		*b = BigInt(t)
		return nil
	}

	return fmt.Errorf("Can't convert %T to common.BigInt", value)
}

func (b *BigInt) FromDB(value []byte) error {
	return b.Scan(value)
}

func (b BigInt) ToDB() ([]byte, error) {
	var bInt = big.Int(b)
	return bInt.Bytes(), nil
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *BigInt) UnmarshalJSON(value []byte) error {
	var s string
	if err := json.Unmarshal(value, &s); err != nil {
		return err
	}

	return b.Scan(s)
}

func (b *BigInt) Cmp(other *BigInt) int {
	v1 := big.Int(*b)
	v2 := big.Int(*other)
	return v1.Cmp(&v2)
}

func (b *BigInt) Divide(t big.Int) error {
	v := big.Int(*b)
	*b = BigInt(*v.Div(&v, &t))
	return nil
}

func (b *BigInt) Int64() int64 {
	v := big.Int(*b)
	return v.Int64()
}

type LNoPrefixHex []byte

func StringToLNoPrefixHex(s string) (r LNoPrefixHex) {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		s = s[2:]
	}
	_ = r.Scan(s)
	return
}

func (n LNoPrefixHex) Hex() string {
	return "0x" + n.String()
}

func (n LNoPrefixHex) Bytes() []byte {
	return []byte(n)
}

func (n LNoPrefixHex) String() string {
	return strings.ToUpper(hexutil.Encode(n)[2:])
}

func (n *LNoPrefixHex) Scan(value interface{}) (err error) {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*n = LNoPrefixHex(v)
		return
	case string:
		if len(v) == 0 {
			return
		}

		*n, err = hexutil.Decode("0x" + v)
		return
	}

	return fmt.Errorf("Can't convert %T to NoPrefixHex", value)
}

func (n *LNoPrefixHex) FromDB(value []byte) error {
	return n.Scan(value)
}

func (n LNoPrefixHex) ToDB() ([]byte, error) {
	return []byte(n), nil
}

func (n LNoPrefixHex) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}

func (n *LNoPrefixHex) UnmarshalJSON(value []byte) error {
	if len(value) == 0 {
		return nil
	}

	var s string
	if err := json.Unmarshal(value, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

func (n LNoPrefixHex) EqualTo(o LNoPrefixHex) bool {
	if len(n) != len(o) {
		return false
	}

	for index := 0; index < len(o); index++ {
		if n[index] != o[index] {
			return false
		}
	}

	return true
}

func StringToNoPrefixHex(s string) (r NoPrefixHex) {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		s = s[2:]
	}
	_ = r.Scan(s)
	return
}

type NoPrefixHex []byte

func (n NoPrefixHex) Hex() string {
	return "0x" + n.String()
}

func (n NoPrefixHex) Bytes() []byte {
	return []byte(n)
}

func (n NoPrefixHex) String() string {
	return hexutil.Encode(n)[2:]
}

func (n *NoPrefixHex) Scan(value interface{}) (err error) {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*n = NoPrefixHex(v)
		return
	case string:
		if len(v) == 0 {
			return
		}

		*n, err = hexutil.Decode("0x" + v)
		return
	}

	return fmt.Errorf("Can't convert %T to NoPrefixHex", value)
}

func (n *NoPrefixHex) FromDB(value []byte) error {
	return n.Scan(value)
}

func (n NoPrefixHex) ToDB() ([]byte, error) {
	return []byte(n), nil
}

func (n NoPrefixHex) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}

func (n *NoPrefixHex) UnmarshalJSON(value []byte) error {
	if len(value) == 0 {
		return nil
	}

	var s string
	if err := json.Unmarshal(value, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

func (n NoPrefixHex) EqualTo(o NoPrefixHex) bool {
	if len(n) != len(o) {
		return false
	}

	for index := 0; index < len(o); index++ {
		if n[index] != o[index] {
			return false
		}
	}

	return true
}

func StringToHex(s string) (r Hex) {
	_ = r.Scan(s)
	return
}

type Hex []byte

func (h Hex) Hex() string {
	return h.String()
}

func (h Hex) NoPrefixHex() string {
	return h.String()[2:]
}

func (h Hex) Bytes() []byte {
	return []byte(h)
}

func (h Hex) String() string {
	return hexutil.Encode(h)
}

func (h Hex) Split(width int) []Hex {
	if width <= 0 {
		panic("invalid parameter")
	}

	r := make([]Hex, 0)
	pos := 0

	for pos < len(h) {
		start := pos
		end := pos + width

		if end >= len(h) {
			end = len(h)
		}

		pos = end
		r = append(r, h[start:end])
	}

	return r
}

func (h *Hex) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case []byte:
		*h = v
		return
	case string:
		if len(v) <= 2 {
			return
		}

		*h, err = hexutil.Decode(v)
		return
	}

	return fmt.Errorf("Can't convert %T to Hex", value)
}

func (h *Hex) FromDB(value []byte) error {
	return h.Scan(value)
}

func (h Hex) ToDB() ([]byte, error) {
	return []byte(h), nil
}

func (h Hex) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *Hex) UnmarshalJSON(value []byte) error {
	if len(value) == 0 {
		return nil
	}

	var s string
	if err := json.Unmarshal(value, &s); err != nil {
		return err
	}

	return h.Scan(s)
}

func (h Hex) EqualTo(o Hex) bool {
	if len(h) != len(o) {
		return false
	}

	for index := 0; index < len(o); index++ {
		if h[index] != o[index] {
			return false
		}
	}

	return true
}

type HexUint64 uint64

func (i *HexUint64) UnmarshalJSON(input []byte) error {
	var s string
	if err := json.Unmarshal(input, &s); err != nil {
		s = string(input)
	}

	int, ok := ParseUint64(s)
	if !ok {
		return fmt.Errorf("invalid hex or decimal integer %q", input)
	}

	*i = HexUint64(int)
	return nil
}

func (i HexUint64) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint64(i))
}

func ParseUint64(s string) (uint64, bool) {
	if s == "" {
		return 0, true
	}
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		v, err := strconv.ParseUint(s[2:], 16, 64)
		return v, err == nil
	}
	v, err := strconv.ParseUint(s, 10, 64)
	return v, err == nil
}

func ConcatHex(params ...[]Hex) []Hex {
	r := make([]Hex, 0)

	for _, param := range params {
		r = append(r, param...)
	}

	return r
}

func SwitchTokenType(t int) string {
	switch t {
	case 2:
		return "ERC721"
	case 6:
		return "ERC1155"
	}
	return ""
}
