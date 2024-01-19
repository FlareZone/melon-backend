package jwt

import (
	"fmt"
	"github.com/FlareZone/melon-backend/config"
	"github.com/golang-jwt/jwt"
	"time"
)

const (
	UserID        = "UserID"
	Bearer        = "Bearer "
	Authorization = "Authorization"
)

type CustomClaims struct {
	jwt.StandardClaims
	UserID string
}

func Generate(userID string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    config.JwtCfg.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken, err := token.SignedString([]byte(config.JwtCfg.Secret))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

func Parse(jwtToken string) (string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JwtCfg.Secret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims[UserID].(string), nil
	} else {
		return "", fmt.Errorf("parse jwt token fail")
	}
}
