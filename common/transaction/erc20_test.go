package transaction

import (
	"context"
	"fmt"
	"github.com/FlareZone/melon-backend/common/hexutil"
	"github.com/FlareZone/melon-backend/common/rpc"
	"github.com/FlareZone/melon-backend/common/signature"
	"github.com/FlareZone/melon-backend/common/signature/eip712"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/assert"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"math/big"
	"testing"
	"time"
)

func TestERC20_BalanceOf(t *testing.T) {
	node, err := rpc.NewRpc(context.Background(), "http://192.168.0.6:8545")
	assert.NoError(t, err)
	erc20 := new(ERC20).WithContractAddress(hexutil.StringToHex("0xdEdDA3fD0ACd5860a362c260D659ec09148B7757")).
		WithRpc(node)
	balanceOf := erc20.BalanceOf(hexutil.StringToHex("0x2863bb864557291664d5faf0cA579c45991e7962"))
	tokenNumber := erc20.QuoDecimals(hexutil.StringToHex("0x2863bb864557291664d5faf0cA579c45991e7962"))
	t.Log("balanceOf: ", balanceOf)
	t.Log("tokenNumber: ", tokenNumber.String())

	name := erc20.Name()
	t.Log("name: ", name)
	symbol := erc20.Symbol()
	t.Log("symbol: ", symbol)
	decimal := erc20.Decimal()
	t.Log("decimal: ", decimal)
}

func TestERC20_QuoDecimals(t *testing.T) {
	node, err := rpc.NewRpc(context.Background(), "http://192.168.0.6:8545")
	assert.NoError(t, err)
	erc20 := new(ERC20).WithContractAddress(hexutil.StringToHex("0x89450afD5Db88f755D6c9a676Dded90032609155")).
		WithRpc(node)

	balanceOf := erc20.BalanceOf(hexutil.StringToHex("0x3415c8EfF9f4066FA7c4e89526735C866B606081"))
	tokenNumber := erc20.QuoDecimals(hexutil.StringToHex("0x3415c8EfF9f4066FA7c4e89526735C866B606081"))
	t.Log("balanceOf: ", balanceOf)
	t.Log("tokenNumber: ", tokenNumber.String())

	name := erc20.Name()
	t.Log("name: ", name)
	symbol := erc20.Symbol()
	t.Log("symbol: ", symbol)
	decimal := erc20.Decimal()
	t.Log("decimal: ", decimal)
}

func TestERC20_PermitNonce(t *testing.T) {
	node, err := rpc.NewRpc(context.Background(), "http://192.168.0.6:8545")
	assert.NoError(t, err)
	erc20 := new(ERC20).WithContractAddress(hexutil.StringToHex("0x87A31aDa5ffDa118D83543CEF9EDCCf9806a0f5f")).
		WithRpc(node)
	nonce, err := erc20.PermitNonce(hexutil.StringToHex("0xcB796Ee82d8824A0A0C8A14E177a99323Cc5A814"))
	assert.NoError(t, err)
	t.Log("nonce: ", nonce)
}

func TestERC20_Transfer(t *testing.T) {
	node, err := rpc.NewRpc(context.Background(), "http://192.168.0.6:8545")
	assert.NoError(t, err)
	erc20 := new(ERC20).WithContractAddress(hexutil.StringToHex("0x87A31aDa5ffDa118D83543CEF9EDCCf9806a0f5f")).
		WithRpc(node).WithPrivateKey(
		signature.GetPrivateKeyFromPrivateHex(
			hexutil.StringToHex("0xd1ffffdd096b4ea087406f792d5de768dbaea9b55a8ff802e9389c0b112dcb4d")))

	hash, err := erc20.Transfer(hexutil.StringToHex("0x3415c8EfF9f4066FA7c4e89526735C866B606081"), "100")
	assert.NoError(t, err)
	t.Log("hash: ", hash)
}

func TestERC20_TransferFrom(t *testing.T) {
	node, err := rpc.NewRpc(context.Background(), "http://192.168.0.6:8545")
	assert.NoError(t, err)
	erc20 := new(ERC20).WithContractAddress(hexutil.StringToHex("0xdEdDA3fD0ACd5860a362c260D659ec09148B7757")).
		WithRpc(node).WithPrivateKey(
		signature.GetPrivateKeyFromPrivateHex(
			hexutil.StringToHex("0xb7e9a3f55f201e349470e8b9ad9f6b927e7900f62520898bad9c2d19fa1162a1")))

	hash, err := erc20.TransferFrom(hexutil.StringToHex("0x3415c8EfF9f4066FA7c4e89526735C866B606081"),
		hexutil.StringToHex("0x2863bb864557291664d5faf0cA579c45991e7962"), "0.8")
	assert.NoError(t, err)
	t.Log("hash: ", hash)
}

func TestERC20_Permit(t *testing.T) {
	node, err := rpc.NewRpc(context.Background(), "http://192.168.0.6:8545")
	assert.NoError(t, err)
	erc20 := new(ERC20).WithContractAddress(hexutil.StringToHex("0xdEdDA3fD0ACd5860a362c260D659ec09148B7757")).
		WithRpc(node).
		WithPrivateKey(signature.GetPrivateKeyFromPrivateHex(hexutil.StringToHex("0xf83ad49b9fe3e8d889026b4205e96dbe11c93a6d369df8fe2d125c4cded33679")))

	chainID, err := node.ChainID()
	assert.NoError(t, err)

	nonce, err := erc20.PermitNonce(hexutil.StringToHex("0x3415c8EfF9f4066FA7c4e89526735C866B606081"))
	assert.NoError(t, err)
	unix := time.Now().Add(time.Hour * 24).Unix()

	var (
		owner     = "0x3415c8EfF9f4066FA7c4e89526735C866B606081"
		spender   = "0x5AD07396DCe3a6aA66eeAEA2f2207D81e35Ba48a"
		amount    = new(big.Int).SetUint64(100000000)
		nonceData = new(big.Int).SetUint64(nonce)
		timestamp = new(big.Int).SetInt64(unix)
	)

	permitToken := eip712.NewPermitToken("Melon", "1", math.NewHexOrDecimal256(int64(chainID)), "0xdEdDA3fD0ACd5860a362c260D659ec09148B7757")
	permitToken.
		WithMessage(owner,
			spender,
			amount,
			nonceData,
			timestamp)
	sign, err := permitToken.Sign(
		signature.GetPrivateKeyFromPrivateHex(hexutil.StringToHex("0xa8344e3d88588d3d1cd587f72b94130566d3b4bb0ce02d71192c372df8c24b76")))
	assert.NoError(t, err)
	hash, err := erc20.Permit(
		owner,
		spender,
		amount,
		timestamp,
		hexutil.StringToHex(sign))
	assert.NoError(t, err)
	t.Log("hash: ", hash)
}

/**

calldata: 0xd505accf0000000000000000000000003415c8eff9f4066fa7c4e89526735c866b6060810000000000000000000000005ad07396dce3a6aa66eeaea2f2207d81e35ba48a0000000000000000000000000000000000000000000000000000000005f5e1000000000000000000000000000000000000000000000000000000000065bb527b000000000000000000000000000000000000000000000000000000000000001c69efaabbecf46341e15bc41dedafaba0422fbb8aa6803f4781aa1f22a3152b8c207b7926b404434ba77755c00af32c4be3d0c6ec4810469361e779ad6e5ba735
{
	owner: 0x3415c8EfF9f4066FA7c4e89526735C866B606081
    spender: 0x5AD07396DCe3a6aA66eeAEA2f2207D81e35Ba48a
    value: 100000000
    deadline: 1706775163
    v: R
	R: 0x69efaabbecf46341e15bc41dedafaba0422fbb8aa6803f4781aa1f22a3152b8c
    S: 0x207b7926b404434ba77755c00af32c4be3d0c6ec4810469361e779ad6e5ba735
}

// return (address, bytes32, bytes32, bytes32, bytes memory)
0x0000000000000000000000003d42df88b37b3fbd04ba2d9c08d0176d698a8d896e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c96016ccfaf671a47da87794bfcb5f3635c3ccd799a4c8fd91b6baab06af421c9173d5da8c266049076729febb271fe516afbd8023e20ee50ae3542ba78e8e603e00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000c06e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c90000000000000000000000003415c8eff9f4066fa7c4e89526735c866b6060810000000000000000000000005ad07396dce3a6aa66eeaea2f2207d81e35ba48a0000000000000000000000000000000000000000000000000000000005f5e10000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000065bb527b
signer: 	0x3D42df88b37b3Fbd04BA2d9C08d0176D698a8D89
hash:  		0x6016ccfaf671a47da87794bfcb5f3635c3ccd799a4c8fd91b6baab06af421c91
structHash: 0x73d5da8c266049076729febb271fe516afbd8023e20ee50ae3542ba78e8e603e
permitHash: 0x6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c9
abiEncode:  0x6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c90000000000000000000000003415c8eff9f4066fa7c4e89526735c866b6060810000000000000000000000005ad07396dce3a6aa66eeaea2f2207d81e35ba48a0000000000000000000000000000000000000000000000000000000005f5e10000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000065bb527b
*/

func TestERC20_111(t *testing.T) {
	permitHashBytes := ethgo.Keccak256([]byte("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)"))
	t.Log(hexutil.BytesToHex(permitHashBytes))
	newType, err := abi.NewType("tuple(bytes32 permitHash, address owner, address spender, uint256 value, uint256 nonce, uint256 deadline)")
	assert.NoError(t, err)
	encode, err := newType.Encode([]interface{}{permitHashBytes,
		"0x3415c8EfF9f4066FA7c4e89526735C866B606081",
		"0x5AD07396DCe3a6aA66eeAEA2f2207D81e35Ba48a",
		new(big.Int).SetInt64(100000000), new(big.Int).SetInt64(1), new(big.Int).SetInt64(1706775163)})
	assert.NoError(t, err)
	t.Log(hexutil.BytesToHex(encode))
	structHash := ethgo.Keccak256(encode)
	fmt.Println("structHash: ", hexutil.BytesToHex(structHash))

	//var (
	//	owner     = "0x3415c8EfF9f4066FA7c4e89526735C866B606081"
	//	spender   = "0x5AD07396DCe3a6aA66eeAEA2f2207D81e35Ba48a"
	//	amount    = new(big.Int).SetUint64(100000000)
	//	nonceData = new(big.Int).SetUint64(1)
	//	timestamp = new(big.Int).SetInt64(1706775163)
	//)

	//permitToken := new(eip712.PermitToken).
	//	WithTypeDataDomain("Melon", "1", math.NewHexOrDecimal256(int64(1337)), "0xdEdDA3fD0ACd5860a362c260D659ec09148B7757")
	//permitToken.
	//	WithMessage(owner,
	//		spender,
	//		amount,
	//		nonceData,
	//		timestamp)
	//sign, err := permitToken.Sign(
	//	signature.GetPrivateKeyFromPrivateHex(hexutil.StringToHex("0xa8344e3d88588d3d1cd587f72b94130566d3b4bb0ce02d71192c372df8c24b76")))
	//assert.NoError(t, err)
	//t.Log(sign)
}

func TestERC20_Result(t *testing.T) {
	newType, err := abi.NewType("tuple(address signer, bytes32 permitHash, bytes32 hash, bytes32 structHash, bytes abiEncode)")
	assert.NoError(t, err)

	data := hexutil.StringToHex("0x0000000000000000000000003415c8eff9f4066fa7c4e89526735c866b6060816e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c97bfaee2ea78aa73e8ee145372d9ec60577bf00550f6194d4ac2fa611863085aa5afa93e8057dadcfa27ab8468aadcefcae40484764f52e1a98902ca0bc78d70200000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000c06e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c90000000000000000000000003415c8eff9f4066fa7c4e89526735c866b6060810000000000000000000000009429ce1b2e9fcbb7edabcb745d4550da2669ae170000000000000000000000000000000000000000000000000000000005f5e10000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000065bc766d")
	type PermitData struct {
		Signer     common.Address `mapstructure:"signer"`
		Hash       []byte         `mapstructure:"hash"`
		StructHash []byte         `mapstructure:"structHash"`
		PermitHash []byte         `mapstructure:"permitHash"`
		AbiEncode  []byte         `mapstructure:"abiEncode"`
	}
	var s PermitData

	err = newType.DecodeStruct(data, &s)
	assert.NoError(t, err)
	t.Log(s.Signer)
	t.Log(hexutil.BytesToHex(s.Hash))
	t.Log(hexutil.BytesToHex(s.StructHash))
	t.Log(hexutil.BytesToHex(s.PermitHash))
	t.Log(hexutil.BytesToHex(s.AbiEncode))
}

func TestERC20_Decimal(t *testing.T) {
	var (
		mlnToken = hexutil.StringToHex("0xdEdDA3fD0ACd5860a362c260D659ec09148B7757")
		relayer  = "0x1e6383a945765Fe83AbE7f89CC5eBB4Cd4614A10"

		account0    = "0x3415c8EfF9f4066FA7c4e89526735C866B606081"
		account3    = "0x9429Ce1b2E9FcbB7eDaBCB745d4550da2669Ae17"
		tokenAmount = "100"
		deadline    = new(big.Int).SetInt64(time.Now().Add(time.Hour * 24).Unix())
	)
	t.Log("deadline: ", deadline.String())
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

	r, s, v := signature.RSV(hexutil.StringToHex(sign))
	var getPermit = "function getPermit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s)"
	method, err := abi.NewMethod(getPermit)
	assert.NoError(t, err)
	encode, err := method.Encode([]interface{}{account0, account3, amount, deadline, v, r.Bytes(), s.Bytes()})
	assert.NoError(t, err)

	toAddress := common.HexToAddress(erc20Token.contractAddress.Hex())

	param := rpc.CallParameter{
		To:   &toAddress,
		Data: encode,
	}
	var result string
	err = rpcNode.CallContext(context.Background(), &result, "eth_call", param.ToArg(), "latest")
	fmt.Println(result)
	assert.NoError(t, err)
}
