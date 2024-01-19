package service

import (
	"fmt"
	"github.com/FlareZone/melon-backend/common/uuid"
	"github.com/FlareZone/melon-backend/internal/model"
	"time"
	"xorm.io/xorm"
)

type SigNonceService interface {
	CreateNonce(ethAddress string) (nonce string, err error)
	FindSigNonceByEthAddress(ethAddress string) *model.SigNonce
	ReGenerate(sigNonce *model.SigNonce) string
	UseNonce(sigNonce *model.SigNonce) bool
}

type SigNonce struct {
	xorm           *xorm.Engine
	expireDuration time.Duration
}

func NewNonce(xorm *xorm.Engine) SigNonceService {
	return &SigNonce{xorm: xorm, expireDuration: time.Minute * 5}
}

func (s *SigNonce) CreateNonce(ethAddress string) (nonce string, err error) {
	nonce = fmt.Sprintf("%s-%08x", uuid.RawUuid(), 1)
	sigNonceModel := model.SigNonce{
		EthAddress: ethAddress,
		NonceToken: nonce,
		UsedNonce:  0,
	}
	_, err = s.xorm.Table(&model.SigNonce{}).Insert(&sigNonceModel)
	return
}

func (s *SigNonce) FindSigNonceByEthAddress(ethAddress string) *model.SigNonce {
	var sigNonce model.SigNonce
	_, _ = s.xorm.Table(&model.SigNonce{}).Where("eth_address = ?", ethAddress).Get(&sigNonce)
	return &sigNonce
}

func (s *SigNonce) ReGenerate(sigNonce *model.SigNonce) string {
	sigNonce.NonceToken = fmt.Sprintf("%s-%08x", uuid.RawUuid(), sigNonce.UsedNonce+1)
	update, _ := s.xorm.Table(&model.SigNonce{}).Update(sigNonce)
	if update > 0 {
		return sigNonce.NonceToken
	}
	return ""
}

func (s *SigNonce) UseNonce(sigNonce *model.SigNonce) bool {
	sigNonce.UsedNonce++
	update, _ := s.xorm.Table(&model.SigNonce{}).Update(sigNonce)
	if update > 0 {
		return true
	}
	return false
}
