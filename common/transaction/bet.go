package transaction

import (
	"context"
	"crypto/ecdsa"
	"github.com/FlareZone/melon-backend/common/hexutil"
	"github.com/FlareZone/melon-backend/common/rpc"
	"github.com/FlareZone/melon-backend/common/signature/eip712"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/log"
	"github.com/umbracle/ethgo/abi"
	"math/big"
)

type BetProposalStatus uint8

func (b BetProposalStatus) Status() uint8 {
	return uint8(b)
}

const (
	BetProposalNotStarted       = BetProposalStatus(0)
	BetProposalVotingInProgress = BetProposalStatus(1)
	BetProposalVotingEnded      = BetProposalStatus(2)
	BetProposalWaitSetting      = BetProposalStatus(3)
	BetProposalSettling         = BetProposalStatus(4)
	BetProposalClosed           = BetProposalStatus(5)
)

type Bet struct {
	contractAddress hexutil.Hex
	rpc             *rpc.EvmRpc
	privateKey      *ecdsa.PrivateKey
	eip712          *eip712.Bet
}

func (b *Bet) WithContractAddress(contractAddress hexutil.Hex) *Bet {
	b.contractAddress = contractAddress
	return b
}

func (b *Bet) WthRpc(rpc *rpc.EvmRpc) *Bet {
	b.rpc = rpc
	return b
}

func (b *Bet) WithPrivateKey(pk *ecdsa.PrivateKey) *Bet {
	b.privateKey = pk
	return b
}

// userNonce 用户的nonce
func (b *Bet) userNonce(owner hexutil.Hex) (uint64, error) {
	nonceAbi, _ := abi.NewMethod("function nonces(address owner)")
	toAddress := common.HexToAddress(b.contractAddress.Hex())
	encode, _ := nonceAbi.Encode([]interface{}{owner.Hex()})
	param := rpc.CallParameter{
		To:   &toAddress,
		Data: encode,
	}
	var result math.HexOrDecimal256
	err := b.rpc.CallContext(context.Background(), &result, "eth_call", param.ToArg(), "latest")
	if err != nil {
		log.Error("nonceAbi is fail", "owner", owner, "err", err)
		return 0, err
	}
	data := big.Int(result)
	return data.Uint64(), nil
}

// CreateProposal 创建提案
func (b *Bet) CreateProposal(proposer string, description string, options []string, amount, userNonce, deadline *big.Int, sign hexutil.Hex) (txHash string, err error) {
	method, err := abi.NewMethod("function createProposal(address proposer, string description, string[] options, uint256 nonce, uint256 amount, uint256 deadline, bytes signature)")
	if err != nil {
		return
	}
	contractCall := new(ContractCall).WithRpc(b.rpc).
		WithPrivateKey(b.privateKey).
		WithMethod(method, proposer, description, options, userNonce, amount, deadline, sign)
	err = contractCall.Call(context.Background(), b.contractAddress, big.NewInt(0))
	return contractCall.TxHash, err
}

// CreateVote 创建投票
func (b *Bet) CreateVote(voter string, proposalId, optionId, userNonce, amount, deadline *big.Int, sign hexutil.Hex) (txHash string, err error) {
	method, err := abi.NewMethod("function vote(address _voter, uint256 _proposalId, uint256 _optionId, uint256 _nonce, uint256 _amount, uint256 _deadline, bytes _signature)")
	if err != nil {
		return
	}
	contractCall := new(ContractCall).WithRpc(b.rpc).
		WithPrivateKey(b.privateKey).
		WithMethod(method, voter, proposalId, optionId, userNonce, amount, deadline, sign)
	err = contractCall.Call(context.Background(), b.contractAddress, big.NewInt(0))
	return contractCall.TxHash, err
}

// Settle 结算
func (b *Bet) Settle(proposalId, winnerOptionId *big.Int) (txHash string, err error) {
	method, err := abi.NewMethod("function settle(uint256 proposalId, uint256 winnerOptionId)")
	if err != nil {
		return
	}
	contractCall := new(ContractCall).WithRpc(b.rpc).
		WithPrivateKey(b.privateKey).
		WithMethod(method, proposalId, winnerOptionId)
	err = contractCall.Call(context.Background(), b.contractAddress, big.NewInt(0))
	return contractCall.TxHash, err
}

// SetFee 设置费率
func (b *Bet) SetFee(open bool, percentage uint8) (txHash string, err error) {
	method, err := abi.NewMethod("function setFee(bool enable, uint8 percentage)")
	if err != nil {
		return
	}
	contractCall := new(ContractCall).WithRpc(b.rpc).
		WithPrivateKey(b.privateKey).
		WithMethod(method, open, percentage)
	err = contractCall.Call(context.Background(), b.contractAddress, big.NewInt(0))
	return contractCall.TxHash, err
}

// UpdateProposalStatus 修改提案状态
func (b *Bet) UpdateProposalStatus(proposalId *big.Int, newStatus BetProposalStatus) (txHash string, err error) {
	method, err := abi.NewMethod("function updateProposalStatus(uint256 proposalId, uint8 newStatus)")
	if err != nil {
		return
	}
	contractCall := new(ContractCall).WithRpc(b.rpc).
		WithPrivateKey(b.privateKey).
		WithMethod(method, proposalId, newStatus.Status())
	err = contractCall.Call(context.Background(), b.contractAddress, big.NewInt(0))
	return contractCall.TxHash, err
}

// QueryProposalStatus 查询提案的状态
func (b *Bet) QueryProposalStatus(proposalId *big.Int) (status BetProposalStatus, err error) {
	getProposalStatus, _ := abi.NewMethod("function getProposalStatus(uint256 proposalId)")
	toAddress := common.HexToAddress(b.contractAddress.Hex())
	encode, _ := getProposalStatus.Encode([]interface{}{proposalId})
	param := rpc.CallParameter{
		To:   &toAddress,
		Data: encode,
	}
	var result string
	err = b.rpc.CallContext(context.Background(), &result, "eth_call", param.ToArg(), "latest")
	if err != nil {
		log.Error("getProposalStatus is fail", "proposalId", proposalId.String(), "err", err)
		return 0, err
	}
	resultTypo, err := abi.NewType("tuple(uint8 status)")
	if err != nil {
		log.Error("new getProposalStatus result abi fail", "proposalId", proposalId.String(), "err", err)
		return 0, err
	}
	var data struct {
		Status uint8 `json:"status"`
	}
	err = resultTypo.DecodeStruct(hexutil.StringToHex(result), &data)
	if err != nil {
		log.Error("new getProposalStatus result abi fail", "proposalId", proposalId.String(), "err", err)
		return 0, err
	}
	status = BetProposalStatus(data.Status)
	return
}

// QueryProposalTotalAmount 查询提案总质押数量
func (b *Bet) QueryProposalTotalAmount(proposalId *big.Int) (amount *big.Int, err error) {
	method, err := abi.NewMethod("function getProposalTotalAmount(uint256 proposalId)")
	if err != nil {
		return
	}
	var result string
	err = b.rpc.CallWithMethod(context.Background(), &result, b.contractAddress, method, proposalId)
	if err != nil {
		return
	}
	resultTypoAbi, _ := abi.NewType("tuple(uint256 amount)")
	var resultTypo struct {
		Amount *big.Int `json:"amount"`
	}
	err = resultTypoAbi.DecodeStruct(hexutil.StringToHex(result), &resultTypo)
	if err != nil {
		return
	}
	amount = new(big.Int).Set(resultTypo.Amount)
	return
}

// QueryProposalOptionTotalAmount 查询提案每个选项质押数量
func (b *Bet) QueryProposalOptionTotalAmount(proposalId, optionId *big.Int) (amount *big.Int, err error) {
	method, err := abi.NewMethod("function getProposalOptionTotalAmount(uint256 proposalId, uint256 optionId)")
	if err != nil {
		return
	}
	var result string
	err = b.rpc.CallWithMethod(context.Background(), &result, b.contractAddress, method, proposalId, optionId)
	if err != nil {
		return
	}
	resultTypoAbi, _ := abi.NewType("tuple(uint256 amount)")
	var resultTypo struct {
		Amount *big.Int `json:"amount"`
	}
	err = resultTypoAbi.DecodeStruct(hexutil.StringToHex(result), &resultTypo)
	if err != nil {
		return
	}
	amount = new(big.Int).Set(resultTypo.Amount)
	return
}

// GrantRelayerRole 添加account为relayer 角色
func (b *Bet) GrantRelayerRole(account string) (txHash string, err error) {
	method, err := abi.NewMethod("function grantRelayerRole(address account)")
	if err != nil {
		return
	}
	contractCall := new(ContractCall).WithRpc(b.rpc).
		WithPrivateKey(b.privateKey).
		WithMethod(method, account)
	err = contractCall.Call(context.Background(), b.contractAddress, big.NewInt(0))
	return contractCall.TxHash, err
}

// RevokeRelayerRole 取消account为relayer角色
func (b *Bet) RevokeRelayerRole(account string) (txHash string, err error) {
	method, err := abi.NewMethod("function revokeRelayerRole(address account)")
	if err != nil {
		return
	}
	contractCall := new(ContractCall).WithRpc(b.rpc).
		WithPrivateKey(b.privateKey).
		WithMethod(method, account)
	err = contractCall.Call(context.Background(), b.contractAddress, big.NewInt(0))
	return contractCall.TxHash, err
}
