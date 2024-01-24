package worker

import (
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/ethereum/go-ethereum/log"
)

func initVerificationCode() {
	code := service.NewVerificationCode()
	go func() {
		if err := code.Run(); err != nil {
			log.Error("run verification code service fail", "err", err)
		}
	}()
}

func InitWorker() {
	initVerificationCode()
}
