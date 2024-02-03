package eip712

import (
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"math/big"
)

type Bet struct {
	Eip712
}

func NewBet(name string, version string, chainId *math.HexOrDecimal256,
	verifyingContract string) *Bet {
	b := new(Bet)
	b.withTypeDataDomain(name, version, chainId, verifyingContract)
	return b
}

// WithProposal ProposalData(address proposer,string description,string[] options,uint256 nonce, uint256 amount, uint256 deadline)
func (b *Bet) WithProposal(proposer string, description string, options []string, nonce *big.Int, amount *big.Int, deadline *big.Int) *Bet {
	b.typedData = &apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"ProposalData": []apitypes.Type{
				{Name: "proposer", Type: "address"},
				{Name: "description", Type: "string"},
				{Name: "options", Type: "string[]"},
				{Name: "nonce", Type: "uint256"},
				{Name: "amount", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
			},
		},
		PrimaryType: "ProposalData",
		Domain:      *b.domain,
		Message: apitypes.TypedDataMessage{
			"proposer":    proposer,
			"description": description,
			"options":     options,
			"nonce":       nonce,
			"amount":      amount,
			"deadline":    deadline,
		},
	}
	return b
}

// WithVote VoteData(address voter,uint256 proposalId, uint256 optionId, uint256 nonce, uint256 amount, uint256 deadline)
func (b *Bet) WithVote(voter string, proposalId *big.Int, optionId *big.Int, nonce *big.Int, amount *big.Int, deadline *big.Int) *Bet {
	b.typedData = &apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"VoteData": []apitypes.Type{
				{Name: "voter", Type: "address"},
				{Name: "proposalId", Type: "uint256"},
				{Name: "optionId", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
				{Name: "amount", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
			},
		},
		PrimaryType: "VoteData",
		Domain:      *b.domain,
		Message: apitypes.TypedDataMessage{
			"voter":      voter,
			"proposalId": proposalId,
			"optionId":   optionId,
			"nonce":      nonce,
			"amount":     amount,
			"deadline":   deadline,
		},
	}
	return b
}
