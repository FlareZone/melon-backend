package contractinteraction

import (
	contractAbi "github.com/FlareZone/melon-backend/common/contract-interaction/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/contract"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/wallet"
	"testing"
)

func init() {

}

func Test(t *testing.T) {

	privKey, _ := wallet.NewWalletFromPrivKey(hexutil.MustDecode("0x641943fa6c3fab18fed274c2b3194f0d71383ecfb9a58b2d70188e693c245510"))

	proposalLogicAbi, _ := abi.NewABIFromList(contractAbi.ProposalLogicABI)

	proposalLogicAddress := ethgo.HexToAddress("0xd92B00D9D2Fc3e51828672bfA6D77D43DF35191c")

	client, _ := jsonrpc.NewClient("https://eth-sepolia.g.alchemy.com/v2/hCl9ueHpm7dG-w_fK-TMw6xK2JFR8SxH")

	proposalLogicContract := contract.NewContract(proposalLogicAddress, proposalLogicAbi, contract.WithJsonRPC(client.Eth()), contract.WithSender(privKey))

	//tx, err := proposalLogicContract.Txn("exchangePoints", big.NewInt(10))
	//doErr := tx.Do()
	//if doErr != nil {
	//	t.Error("交换积分错误", err)
	//}
	//txInfo, _ := tx.Wait()
	//t.Log("交易信息", txInfo)
	//if err != nil {
	//	t.Error("交换积分错误", err)
	//}
	balance, _ := proposalLogicContract.Call("balances", ethgo.Latest, ethgo.HexToAddress("0xc0ee714715108b1a6795391f7e05a044d795ba70"))
	t.Log("账户余额", balance["0"])
}
