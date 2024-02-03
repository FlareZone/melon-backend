package transaction

import (
	"context"
	"github.com/FlareZone/melon-backend/common/hexutil"
	"github.com/FlareZone/melon-backend/common/rpc"
	"github.com/FlareZone/melon-backend/common/signature"
	"github.com/FlareZone/melon-backend/common/signature/eip712"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

var (
	accounts = map[string]string{}
	rpcUrl   = "http://192.168.0.6:8545"
	rpcNode  = &rpc.EvmRpc{}
	account0 = "0xcB796Ee82d8824A0A0C8A14E177a99323Cc5A814"
	account1 = "0x3415c8EfF9f4066FA7c4e89526735C866B606081"
	account2 = "0x5AD07396DCe3a6aA66eeAEA2f2207D81e35Ba48a"
	account3 = "0x9429Ce1b2E9FcbB7eDaBCB745d4550da2669Ae17"
	account4 = "0x1e6383a945765Fe83AbE7f89CC5eBB4Cd4614A10"
)

func init() {
	accounts = map[string]string{
		account0: "0xd1ffffdd096b4ea087406f792d5de768dbaea9b55a8ff802e9389c0b112dcb4d",
		account1: "0xa8344e3d88588d3d1cd587f72b94130566d3b4bb0ce02d71192c372df8c24b76",
		account2: "0xb7e9a3f55f201e349470e8b9ad9f6b927e7900f62520898bad9c2d19fa1162a1",
		account3: "0x4f2da037c1423f028c8c9aa13bd25c52047c6332f9446d4a653958871886bb1e",
		account4: "0x73805f725bdec4b3800995678ed56d0d83ddbde995a6be046a32b580729bc25a",
	}
	rpcNode, _ = rpc.NewRpc(context.Background(), rpcUrl)
}

/**
 * accounts[0] 部署了 mln token: 0x87A31aDa5ffDa118D83543CEF9EDCCf9806a0f5f
 * accounts[1] 部署了 bet contract: 0xBD8392F1a5A2b05612a25d73073306E3E1825551
 *
 * 1. 签名：accounts[0] permit 100 mln 代币给到 accounts[3]
 * 2. relayer: accounts[4] 拿着签名上链
 * 3. accounts[3] transfer from 10 mln 给到 accounts[2]
 * 4. 签名：account2 permit 20 mln 代币给到合约 0xBD8392F1a5A2b05612a25d73073306E3E1825551
 * 4. 签名：accounts[2] 创建proposal, 并质押10 mln
 * 5. relayer: accounts[4] 拿着签名上链
 */

func TestBet_Step1_Account0_permit_100_mln(t *testing.T) {
	var (
		mlnToken = hexutil.StringToHex("0x87A31aDa5ffDa118D83543CEF9EDCCf9806a0f5f")
		relayer  = account4

		account0    = "0x3415c8EfF9f4066FA7c4e89526735C866B606081"
		account3    = "0x9429Ce1b2E9FcbB7eDaBCB745d4550da2669Ae17"
		tokenAmount = "100"
		deadline    = new(big.Int).SetInt64(time.Now().Add(time.Hour * 24).Unix())
	)
	chainId, err := rpcNode.ChainID()
	assert.NoError(t, err)
	t.Log("chainId: ", chainId)
	erc20Token := new(ERC20).WithContractAddress(mlnToken).
		WithRpc(rpcNode).
		WithPrivateKey(signature.GetPrivateKeyFromPrivateHex(hexutil.StringToHex(accounts[relayer])))
	of := erc20Token.BalanceOf(hexutil.StringToHex(account0))
	t.Log("account0 balance: ", of.String())
	amount, err := erc20Token.TokenAmountToBigInt(tokenAmount)
	assert.NoError(t, err)
	t.Log("account0 permit token: ", amount.String())
	nonce, err := erc20Token.PermitNonce(hexutil.StringToHex(account0))
	assert.NoError(t, err)
	permitToken := eip712.NewPermitToken("Melon", "1", math.NewHexOrDecimal256(int64(chainId)), mlnToken.Hex())
	permitToken.WithMessage(account0, account3, amount, new(big.Int).SetUint64(nonce), deadline)
	sign, err := permitToken.Sign(signature.GetPrivateKeyFromPrivateHex(hexutil.StringToHex(accounts[account0])))
	assert.NoError(t, err)
	erc20Token.WithEip712(permitToken)
	hash, err := erc20Token.Permit(account0, account3, amount, deadline, hexutil.StringToHex(sign))
	assert.NoError(t, err)
	t.Log("hash: ", hash)
}

func TestBet_Step2_Account3_transfer_from_account2_10_mln(t *testing.T) {
	var (
		mlnToken    = hexutil.StringToHex("0x87A31aDa5ffDa118D83543CEF9EDCCf9806a0f5f")
		account0    = "0x3415c8EfF9f4066FA7c4e89526735C866B606081"
		account3    = "0x9429Ce1b2E9FcbB7eDaBCB745d4550da2669Ae17"
		account2    = "0x5AD07396DCe3a6aA66eeAEA2f2207D81e35Ba48a"
		tokenAmount = "10"
	)

	erc20Token := new(ERC20).WithContractAddress(mlnToken).
		WithRpc(rpcNode).
		WithPrivateKey(signature.GetPrivateKeyFromPrivateHex(hexutil.StringToHex(accounts[account3])))
	balanceOf := erc20Token.BalanceOf(hexutil.StringToHex(account3))
	t.Log("balanceOf: ", balanceOf.String())

	hash, err := erc20Token.TransferFrom(hexutil.StringToHex(account0), hexutil.StringToHex(account2), tokenAmount)
	assert.NoError(t, err)
	t.Log("hash: ", hash)
}

func TestBet_BalanceOf(t *testing.T) {
	var (
		mlnToken = hexutil.StringToHex("0x87A31aDa5ffDa118D83543CEF9EDCCf9806a0f5f")
		account3 = "0x9429Ce1b2E9FcbB7eDaBCB745d4550da2669Ae17"
		account2 = "0x5AD07396DCe3a6aA66eeAEA2f2207D81e35Ba48a"
	)
	erc20Token := new(ERC20).WithContractAddress(mlnToken).
		WithRpc(rpcNode).
		WithPrivateKey(signature.GetPrivateKeyFromPrivateHex(hexutil.StringToHex(accounts[account3])))
	balanceOf := erc20Token.BalanceOf(hexutil.StringToHex(account2))
	t.Log("balanceOf: ", balanceOf.String())
}

func TestBet_step4_permit_20_mnl_to_contract(t *testing.T) {
	var (
		mlnToken    = hexutil.StringToHex("0x87A31aDa5ffDa118D83543CEF9EDCCf9806a0f5f")
		betContract = hexutil.StringToHex("0xaF336D318835929654A359b1fD1907aF019fE6A1")
		relayer     = account4
		tokenAmount = "20"
		deadline    = new(big.Int).SetInt64(time.Now().Add(time.Hour * 24).Unix())
	)

	chainId, err := rpcNode.ChainID()
	assert.NoError(t, err)
	t.Log("chainId: ", chainId)
	erc20Token := new(ERC20).WithContractAddress(mlnToken).
		WithRpc(rpcNode).
		WithPrivateKey(signature.GetPrivateKeyFromPrivateHex(hexutil.StringToHex(accounts[relayer])))
	of := erc20Token.BalanceOf(hexutil.StringToHex(account2))
	t.Log("account0 balance: ", of.String())
	amount, err := erc20Token.TokenAmountToBigInt(tokenAmount)
	assert.NoError(t, err)
	t.Log("account0 permit token: ", amount.String())
	nonce, err := erc20Token.PermitNonce(hexutil.StringToHex(account2))
	assert.NoError(t, err)
	permitToken := eip712.NewPermitToken("Melon", "1", math.NewHexOrDecimal256(int64(chainId)), mlnToken.Hex())
	permitToken.WithMessage(account2, betContract.Hex(), amount, new(big.Int).SetUint64(nonce), deadline)
	sign, err := permitToken.Sign(signature.GetPrivateKeyFromPrivateHex(hexutil.StringToHex(accounts[account2])))
	assert.NoError(t, err)
	erc20Token.WithEip712(permitToken)
	hash, err := erc20Token.Permit(account2, betContract.Hex(), amount, deadline, hexutil.StringToHex(sign))
	assert.NoError(t, err)
	t.Log("hash: ", hash)
}

func TestBet_betAllowance(t *testing.T) {
	var (
		mlnToken    = hexutil.StringToHex("0x87A31aDa5ffDa118D83543CEF9EDCCf9806a0f5f")
		betContract = hexutil.StringToHex("0xaF336D318835929654A359b1fD1907aF019fE6A1")
		relayer     = account4
	)

	chainId, err := rpcNode.ChainID()
	assert.NoError(t, err)
	t.Log("chainId: ", chainId)
	erc20Token := new(ERC20).WithContractAddress(mlnToken).
		WithRpc(rpcNode).
		WithPrivateKey(signature.GetPrivateKeyFromPrivateHex(hexutil.StringToHex(accounts[relayer])))
	of := erc20Token.BalanceOf(hexutil.StringToHex(account2))
	t.Log("account0 balance: ", of.String())
	allowance, err := erc20Token.Allowance(account2, betContract.Hex())
	assert.NoError(t, err)
	t.Log("allowance: ", allowance)

}

func TestBet_CreateProposal(t *testing.T) {
	var (
		mlnToken    = hexutil.StringToHex("0x87A31aDa5ffDa118D83543CEF9EDCCf9806a0f5f")
		betContract = hexutil.StringToHex("0xaF336D318835929654A359b1fD1907aF019fE6A1")
		proposer    = hexutil.StringToHex(account2)
		description = "Hello world!"
		options     = []string{"hello", "world"}
		tokenAmount = "15"
		deadline    = new(big.Int).SetInt64(time.Now().Add(time.Hour * 24).Unix())
		relayer     = account4
	)

	erc20Token := new(ERC20).WithContractAddress(mlnToken).
		WithRpc(rpcNode)
	amount, err := erc20Token.TokenAmountToBigInt(tokenAmount)
	assert.NoError(t, err)
	t.Log("deadline: ", deadline)
	betCall := new(Bet).WithContractAddress(betContract).WithPrivateKey(
		signature.GetPrivateKeyFromPrivateHex(
			hexutil.StringToHex(accounts[relayer]))).WthRpc(rpcNode)

	nonce, err := betCall.userNonce(proposer)
	assert.NoError(t, err)

	chainId, err := rpcNode.ChainID()
	assert.NoError(t, err)
	t.Log("chainId: ", chainId)
	bet := eip712.NewBet("BetContractDomain", "1", math.NewHexOrDecimal256(int64(chainId)), betContract.Hex())
	bet.WithProposal(proposer.Hex(), description, options, new(big.Int).SetUint64(nonce), amount, deadline)
	sign, err := bet.Sign(signature.GetPrivateKeyFromPrivateHex(
		hexutil.StringToHex(accounts[account2])))
	assert.NoError(t, err)

	hash, err := betCall.CreateProposal(proposer.Hex(), description, options, amount, new(big.Int).SetUint64(nonce), deadline, hexutil.StringToHex(sign))
	assert.NoError(t, err)
	t.Log("txHash: ", hash)
}

func TestBet_GrantRelayerRole(t *testing.T) {
	var (
		betContract = hexutil.StringToHex("0xaF336D318835929654A359b1fD1907aF019fE6A1")
		relayer     = account4
		deployer    = account1
	)

	betCall := new(Bet).WithContractAddress(betContract).WithPrivateKey(
		signature.GetPrivateKeyFromPrivateHex(
			hexutil.StringToHex(accounts[deployer]))).WthRpc(rpcNode)

	hash, err := betCall.GrantRelayerRole(relayer)
	assert.NoError(t, err)
	t.Log("txHash: ", hash)
}

func TestBet_UpdateProposalStatus(t *testing.T) {
	var (
		betContract = hexutil.StringToHex("0xaF336D318835929654A359b1fD1907aF019fE6A1")
		relayer     = account4
	)
	betCall := new(Bet).WithContractAddress(betContract).WithPrivateKey(
		signature.GetPrivateKeyFromPrivateHex(
			hexutil.StringToHex(accounts[relayer]))).WthRpc(rpcNode)

	hash, err := betCall.UpdateProposalStatus(new(big.Int).SetInt64(0), BetProposalVotingInProgress)
	assert.NoError(t, err)
	t.Log("txHash: ", hash)
}

func TestBet_GetProposalStatus(t *testing.T) {
	var (
		betContract = hexutil.StringToHex("0xaF336D318835929654A359b1fD1907aF019fE6A1")
		relayer     = account4
	)
	betCall := new(Bet).WithContractAddress(betContract).WithPrivateKey(
		signature.GetPrivateKeyFromPrivateHex(
			hexutil.StringToHex(accounts[relayer]))).WthRpc(rpcNode)
	status, err := betCall.QueryProposalStatus(new(big.Int).SetInt64(0))
	assert.NoError(t, err)
	t.Log("status: ", status.Status())
}
