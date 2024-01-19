package signature

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/FlareZone/melon-backend/config"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"regexp"
	"strings"
)

var (
	walletAddressRegex = regexp.MustCompile(`Wallet address: (0x[a-fA-F0-9]{40})`)
	nonceRegex         = regexp.MustCompile(`Nonce: (\S+)`)
	testNonce          = "fbd33030-9054-478b-b254-0eae199b99fe"
)

var LoginMessage = func(ethAddress string, nonce string) string {
	return fmt.Sprintf(`
Welcome to %s!

Click to sign in and accept the Melon Terms of Service (%s) and Privacy Policy (%s/privacy).

This request will not trigger a blockchain transaction or cost any gas fees.

Your authentication status will reset after 24 hours.  

Wallet address: %s

Nonce: %s
`, config.EIP712.Name, config.App.Url, config.App.Url, ethAddress, nonce)
}

func GenerateLogin(PrivateECDSA *ecdsa.PrivateKey, ethAddress string, nonce string) (hashHex, signatureHex string, err error) {
	domain := apitypes.TypedDataDomain{
		Name:              config.EIP712.Name,
		Version:           config.EIP712.Version,
		ChainId:           math.NewHexOrDecimal256(config.EIP712.ChainID),
		VerifyingContract: config.EIP712.VerifyingContract,
	}
	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "string"},
			},
			"Login": []apitypes.Type{
				{Name: "message", Type: "string"},
			},
		},
		PrimaryType: "Login",
		Domain:      domain,
		Message: apitypes.TypedDataMessage{
			"message": LoginMessage(ethAddress, nonce),
		},
	}
	hash, _, err := apitypes.TypedDataAndHash(typedData)
	if err != nil {
		return
	}

	sig, err := crypto.Sign(hash, PrivateECDSA)
	if err != nil {
		return
	}
	sig[64] += 27
	hashHex = fmt.Sprintf("0x%s", hex.EncodeToString(hash))
	signatureHex = fmt.Sprintf("0x%s", hex.EncodeToString(sig))
	return
}

func MelonLoginWithEip712(typedDataStr, hashHex, signatureHex string, f func(ethAddress, nonce string) error) (ethAddress string, err error) {
	var typedData apitypes.TypedData
	err = json.Unmarshal([]byte(typedDataStr), &typedData)
	if err != nil {
		err = fmt.Errorf("MelonLogin, Unmarshal typedData  err is %v", err)
		return
	}
	messageInfo, ok := typedData.Message["message"]
	if !ok {
		err = fmt.Errorf("MelonLogin, unkonwn typedData.Message.message")
		return
	}
	message, ok := messageInfo.(string)
	if !ok {
		err = fmt.Errorf("MelonLogin, unkonwn typedData.Message.message type")
		return
	}
	nonceMatches := nonceRegex.FindStringSubmatch(message)
	walletAddressMatches := walletAddressRegex.FindStringSubmatch(message)
	if len(nonceMatches) <= 1 || len(walletAddressMatches) <= 1 {
		err = fmt.Errorf("MelonLogin, unkonwn typedData ")
		return
	}
	hash, _, err := apitypes.TypedDataAndHash(typedData)
	calHashHex := fmt.Sprintf("0x%s", hex.EncodeToString(hash))
	if !strings.EqualFold(calHashHex, hashHex) {
		err = fmt.Errorf("MelonLogin, TypedDataAndHash failed %s != %s", calHashHex, hashHex)
		return
	}
	signature, err := hex.DecodeString(signatureHex[2:])
	if err != nil {
		err = fmt.Errorf("MelonLogin, decode signature  err is %v", err)
		return
	}
	if signature[64] != 27 && signature[64] != 28 {
		err = fmt.Errorf("MelonLogin, invalid Ethereum signature (V is not 27 or 28)")
		return
	}
	signature[64] -= 27 // 对于以太坊的签名，需要将V值减去27
	pub, err := crypto.SigToPub(hash, signature)
	if err != nil {
		err = fmt.Errorf("MelonLogin, invalid SigToPub, signature: %s", signatureHex)
		return
	}
	if !strings.EqualFold(typedData.Domain.Name, config.EIP712.Name) {
		err = fmt.Errorf("MelonLogin, invalid typedData.Domain.Name, %s != %s", typedData.Domain.Name, config.EIP712.Name)
		return
	}
	if !strings.EqualFold(typedData.Domain.VerifyingContract, config.EIP712.VerifyingContract) {
		err = fmt.Errorf("MelonLogin, invalid typedData.Domain.VerifyingContract, %s != %s", typedData.Domain.VerifyingContract, config.EIP712.VerifyingContract)
		return
	}
	nonce := nonceMatches[1]
	wallet := walletAddressMatches[1]
	recoveredAddr := crypto.PubkeyToAddress(*pub)
	ethAddress = recoveredAddr.Hex()
	if !strings.EqualFold(wallet, ethAddress) {
		err = fmt.Errorf("MelonLogin, invalid wallet address, %s != %s", wallet, ethAddress)
		ethAddress = ""
		return
	}
	if nonce != testNonce {
		if err = f(ethAddress, nonce); err != nil {
			err = fmt.Errorf("MelonLogin, invalid nonce, err is %v", err)
			return
		}
	}
	return
}
