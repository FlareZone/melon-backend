package handler

import (
	"encoding/json"
	"fmt"
	"github.com/FlareZone/melon-backend/common/jwt"
	"github.com/FlareZone/melon-backend/common/signature"
	"github.com/FlareZone/melon-backend/common/uuid"
	"github.com/FlareZone/melon-backend/config"
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/FlareZone/melon-backend/internal/types"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"time"
)

type AuthHandler struct {
	user     service.UserService
	sigNonce service.SigNonceService
}

func NewAuthHandler(user service.UserService, sigNonce service.SigNonceService) *AuthHandler {
	return &AuthHandler{user: user, sigNonce: sigNonce}
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
	user := a.user.FindUserByEmail(googleOauthInfo.Email)
	if user.UUID == "" {
		user = googleOauthInfo.User()
		if !a.user.Register(*user) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo"})
			return
		}
	}

	jwtToken, err := jwt.Generate(user.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
		return
	}
	response.JsonSuccessWithMessage(c, jwtToken, "Login successful!")
	return
}

func (a *AuthHandler) EthereumEip712Signature(c *gin.Context) {
	var params EthereumEip712SignatureRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	ethAddress, err := signature.MelonLoginWithEip712(params.GetTypedData(), params.TypedDataHash, params.Signature, a.checkLoginNonce)
	if err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	user := a.user.FindUserByEthAddress(ethAddress)
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
	}
	jwtToken, err := jwt.Generate(user.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
		return
	}
	response.JsonSuccessWithMessage(c, jwtToken, "Login successful!")
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

func (a *AuthHandler) checkLoginNonce(ethAddress, nonce string) error {
	sigNonce := a.sigNonce.FindSigNonceByEthAddress(ethAddress)
	if strings.EqualFold(sigNonce.NonceToken, nonce) {
		a.sigNonce.UseNonce(sigNonce)
		return nil
	}
	return fmt.Errorf("nonce is not equal, %s != %s", sigNonce.NonceToken, nonce)
}