package contractinteraction

import (
	contractAbi "github.com/FlareZone/melon-backend/common/contract-interaction/abi"
	"github.com/FlareZone/melon-backend/config"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/inconshreveable/log15"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/contract"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/wallet"
	"math/big"
)

var (
	log = log15.New("m", "contract-interaction")
)

type ProposalLogicContract struct {
	Instance      *contract.Contract
	name          string
	deployNetwork string
}

func NewProposalLogicContract(privateKey string) *ProposalLogicContract {
	proposalLogicContract := new(ProposalLogicContract)

	proposalLogicAbi, _ := abi.NewABIFromList(contractAbi.ProposalLogicABI)

	proposalLogicAddress := ethgo.HexToAddress(config.SmartContract.ProposalLogicContractAddress)

	client, _ := jsonrpc.NewClient(config.SmartContract.Jsonrpc)

	privKey, _ := wallet.NewWalletFromPrivKey(hexutil.MustDecode(privateKey))

	proposalLogicContract.Instance = contract.NewContract(proposalLogicAddress, proposalLogicAbi, contract.WithJsonRPC(client.Eth()), contract.WithSender(privKey))

	proposalLogicContract.name = "ProposalLogic"

	proposalLogicContract.deployNetwork = "eth_sepolia"

	return proposalLogicContract
}

func (proposalLogicContract ProposalLogicContract) ExchangePoints(points uint64) bool {
	tx, txErr := proposalLogicContract.Instance.Txn("exchangePoints", big.NewInt(int64(points)))
	if txErr != nil {
		log.Error("创建活动错误", "err", txErr)
		return false
	}
	err := tx.Do()
	if err != nil {
		log.Error("交换积分错误", "err", err)
		return false
	}
	txInfo, waitErr := tx.Wait()
	if waitErr != nil {
		log.Error("等待活动错误", "err", waitErr)
		return false
	}
	log.Info("交换积分信息", "txInfo", txInfo)
	return true
}

func (proposalLogicContract ProposalLogicContract) balances(publicKey string) *big.Int {
	balance, err := proposalLogicContract.Instance.Call("balances", ethgo.Latest, ethgo.HexToAddress(publicKey))
	if err != nil {
		log.Error("查询账户余额错误", err)
	}
	return balance["0"].(*big.Int)
}
