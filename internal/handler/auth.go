package handler

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/FlareZone/melon-backend/common/consts"
	"github.com/FlareZone/melon-backend/common/jwt"
	"github.com/FlareZone/melon-backend/common/signature"
	"github.com/FlareZone/melon-backend/common/uuid"
	"github.com/FlareZone/melon-backend/config"
	"github.com/FlareZone/melon-backend/internal/ginctx"
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/FlareZone/melon-backend/internal/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"io"
	"net/http"
	"time"
)

type AuthHandler struct {
	user             service.UserService
	sigNonce         service.SigNonceService
	verificationCode service.VerificationCodeService
}

func NewAuthHandler(user service.UserService, sigNonce service.SigNonceService) *AuthHandler {
	return &AuthHandler{user: user, sigNonce: sigNonce, verificationCode: service.NewVerificationCode()}
}

func (a *AuthHandler) GoogleOauthCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := config.GoogleOauthCfg.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to exchange token", "error": err.Error()})
		return
	}
	// To retrieve user's information from Google's UserInfo endpoint
	client := config.GoogleOauthCfg.Client(c, token)
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
		return
	}
	defer userinfo.Body.Close()

	bodyBytes, err := io.ReadAll(userinfo.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
		return
	}
	var googleOauthInfo types.GoogleOAuthInfo
	err = json.Unmarshal(bodyBytes, &googleOauthInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
		return
	}
	if !googleOauthInfo.EmailVerified {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
		return
	}
	var actionType = AuthActionLogin.String()
	user := a.user.FindUserByEmail(googleOauthInfo.Email)
	if user.UUID == "" {
		user = googleOauthInfo.User()
		if !a.user.Register(*user) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo"})
			return
		}
		actionType = AuthActionRegister.String()
	}

	jwtToken, err := jwt.Generate(user.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
		return
	}

	c.SetCookie(consts.JwtCookie, jwtToken, 24*3600, "/", config.App.Domain(), false, true)

	result := map[string]interface{}{
		"action_type": actionType,
		"jwt_token":   jwtToken,
		"expired_at":  time.Now().Add(time.Hour * 24).UnixMilli(),
	}
	response.JsonSuccessWithMessage(c, result, "Login successful!")
	return
}

func (a *AuthHandler) EthereumEip712Signature(c *gin.Context) {
	var params EthereumEip712SignatureRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	ethAddress, err := signature.MelonLoginWithEip712(params.GetTypedData(), params.TypedDataHash, params.Signature, a.getEthAddressNonceToken)
	if err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	user := a.user.FindUserByEthAddress(ethAddress)
	var actionType = AuthActionLogin.String()
	if user.UUID == "" {
		user = &model.User{
			EthAddress: &ethAddress,
			NickName:   &ethAddress,
			UUID:       uuid.Uuid(),
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}
		if !a.user.Register(*user) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo"})
			return
		}
		actionType = AuthActionRegister.String()
	}
	jwtToken, err := jwt.Generate(user.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
		return
	}
	c.SetCookie(consts.JwtCookie, jwtToken, 24*3600, "/", config.App.Domain(), false, true)
	result := map[string]interface{}{
		"action_type": actionType,
		"jwt_token":   jwtToken,
		"expired_at":  time.Now().Add(time.Hour * 24).UnixMilli(),
	}
	response.JsonSuccessWithMessage(c, result, "Login successful!")
	return
}

func (a *AuthHandler) EthereumEip712SignatureNonce(c *gin.Context) {
	var params EthereumEip712SignatureNonceRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	sigNonce := a.sigNonce.FindSigNonceByEthAddress(params.EthAddress)
	if sigNonce.ID > 0 {
		response.JsonSuccess(c, a.sigNonce.ReGenerate(sigNonce))
		return
	}
	nonce, err := a.sigNonce.CreateNonce(params.EthAddress)
	if err != nil {
		err = fmt.Errorf("createNonce fail, err is %v", err)
		response.JsonFail(c, response.StatusInternalServerError, err.Error())
		return
	}
	response.JsonSuccess(c, nonce)
	return
}

func (a *AuthHandler) SendVerificationCode(c *gin.Context) {
	var params EmailVerificationCodeRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	a.verificationCode.SendLoginVerificationCode(params.To)
	response.JsonSuccess(c, "send success")
}

// LoginWithEmail 如果没有注册，则注册
func (a *AuthHandler) LoginWithEmail(c *gin.Context) {
	var params EmailLoginRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	if !a.verificationCode.VerifyEmailCode(params.Email, params.Code) {
		response.JsonFail(c, response.BadEmailVerificationCode, fmt.Sprintf("email code is invalid"))
		return
	}

	var actionType = AuthActionLogin.String()
	user := a.user.FindUserByEmail(params.Email)
	if user.UUID == "" {
		emailVerify := true
		nickName := fmt.Sprintf("melon_%s", uuid.Uuid())
		user = &model.User{
			UUID:        uuid.Uuid(),
			Email:       &params.Email,
			NickName:    &nickName,
			EmailVerify: &emailVerify,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
		if !a.user.Register(*user) {
			response.JsonFail(c, response.StatusInternalServerError, "register fail")
			return
		}
		actionType = AuthActionRegister.String()
	}

	jwtToken, err := jwt.Generate(user.UUID)
	if err != nil {
		response.JsonFail(c, response.StatusInternalServerError, "generate jwt token fail")
		return
	}

	c.SetCookie(consts.JwtCookie, jwtToken, 24*3600, "/", config.App.Domain(), false, true)

	result := map[string]interface{}{
		"action_type": actionType,
		"jwt_token":   jwtToken,
		"expired_at":  time.Now().Add(time.Hour * 24).UnixMilli(),
	}
	response.JsonSuccessWithMessage(c, result, "Login successful!")
	return
}

func (a *AuthHandler) getEthAddressNonceToken(ethAddress string) string {
	sigNonce := a.sigNonce.FindSigNonceByEthAddress(ethAddress)
	return sigNonce.NonceToken
}

func (a *AuthHandler) Refresh(c *gin.Context) {
	user := ginctx.AuthUser(c)
	jwtToken, err := jwt.Generate(user.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
		return
	}
	c.SetCookie(consts.JwtCookie, jwtToken, 24*3600, "/", config.App.Domain(), false, true)
	result := map[string]interface{}{
		"jwt_token":  jwtToken,
		"expired_at": time.Now().Add(time.Hour * 24).UnixMilli(),
	}
	response.JsonSuccessWithMessage(c, result, "Refresh successful!")
	return
}

func (a *AuthHandler) GetPayload(c *gin.Context) {
	nonce := c.Param("nonce")
	privateKey := a.GetPrivateKey()
	publicKey := privateKey.PublicKey
	address := crypto.PubkeyToAddress(publicKey).Hex()
	typedDataHex, hashHex, signatureHex, _ := signature.GenerateLogin(privateKey, address, nonce)
	var result struct {
		TypedData     string
		TypedDataHash string
		Signature     string
	}
	result.TypedData = typedDataHex
	result.TypedDataHash = hashHex
	result.Signature = signatureHex
	response.JsonSuccess(c, result)

}

func (a *AuthHandler) GetPrivateKey() *ecdsa.PrivateKey {
	mnemonic := "check antique innocent spice much neglect split lottery trouble twelve report tennis"
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, _ := bip32.NewMasterKey(seed)
	purposeKey, _ := masterKey.NewChildKey(0x8000002C)
	coinTypeKey, _ := purposeKey.NewChildKey(0x8000003C)
	accountKey, _ := coinTypeKey.NewChildKey(0x80000000)
	changeKey, _ := accountKey.NewChildKey(0)
	addressKey, _ := changeKey.NewChildKey(0)

	// 使用addressKey的 PrivateKey 方法
	privateKey, _ := crypto.ToECDSA(addressKey.Key)
	return privateKey
}
