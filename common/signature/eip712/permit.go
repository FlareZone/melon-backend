package eip712

import (
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/umbracle/ethgo/abi"
	"math/big"
)

// PermitToken Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)
type PermitToken struct {
	Eip712
}

func NewPermitToken(name string, version string, chainId *math.HexOrDecimal256,
	verifyingContract string) *PermitToken {
	p := new(PermitToken)
	p.withTypeDataDomain(name, version, chainId, verifyingContract)
	return p
}

func (p *PermitToken) PermitMethod() (*abi.Method, error) {
	return abi.NewMethod("function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s)")
}

func (p *PermitToken) WithMessage(owner, spender string, amount, nonce, deadline *big.Int) *PermitToken {
	p.typedData = &apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"Permit": []apitypes.Type{
				{Name: "owner", Type: "address"},
				{Name: "spender", Type: "address"},
				{Name: "value", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
			},
		},
		PrimaryType: "Permit",
		Domain:      *p.domain,
		Message: apitypes.TypedDataMessage{
			"owner":    owner,
			"spender":  spender,
			"value":    amount,
			"deadline": deadline,
			"nonce":    nonce,
		},
	}
	return p
}
