package signature

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestProposalDataHash(t *testing.T) {
	// 定义要计算哈希的字符串
	data := "ProposalData(address proposer,string description,string[] options,uint256 nonce,uint256 amount, uint256 deadline)"

	// 使用 keccak256 计算哈希
	hash := crypto.Keccak256Hash([]byte(data))

	// 打印哈希值的十六进制表示
	fmt.Printf("Keccak256 Hash: %s\n", hash.Hex())
}

func TestVoteDataHash(t *testing.T) {
	// 定义要计算哈希的字符串
	data := "VoteData(address voter,uint256 proposalId, uint256 optionId, uint256 nonce, uint256 amount, uint256 deadline)"

	// 使用 keccak256 计算哈希
	hash := crypto.Keccak256Hash([]byte(data))

	// 打印哈希值的十六进制表示
	fmt.Printf("Keccak256 Hash: %s\n", hash.Hex())
}
