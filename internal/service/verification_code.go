package service

import (
	"context"
	"github.com/FlareZone/melon-backend/common/rdbkey"
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	emailChannel = make(chan *VerificationMail, 5)
)

type VerificationCodeService interface {
	SendLoginVerificationCode(to string)
	VerifyEmailCode(to string, code string) bool
	Run() error
}

type VerificationCode struct {
	Mail   MailService
	Redis  *redis.Client
	locker *redislock.Client
}

func NewVerificationCode() VerificationCodeService {
	rdb := components.Redis
	return &VerificationCode{Mail: NewGoogleMail(), Redis: rdb, locker: redislock.New(rdb)}
}

func (v *VerificationCode) generateCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// 生成随机数的一种方法，避免0和1
	number := rand.Intn(8) + 2 // 2-9
	for i := 0; i < 5; i++ {
		nextDigit := rand.Intn(10) // 0-9
		if nextDigit == 0 {
			nextDigit = rand.Intn(8) + 2 // 再次避免0和1
		}
		number = number*10 + nextDigit
	}
	return strconv.Itoa(number)
}

func (v *VerificationCode) VerifyEmailCode(to string, code string) bool {
	mailLoginKey := rdbkey.MailLogin(to)
	result, err := v.Redis.Get(context.Background(), mailLoginKey).Result()
	if err != nil {
		return false
	}
	_, _ = v.Redis.Del(context.Background(), mailLoginKey).Result()
	return strings.EqualFold(result, code)
}

// This function is used to send a login verification code to the user's email address.
func (v *VerificationCode) SendLoginVerificationCode(to string) {
	// Generate a verification code
	code := v.generateCode()
	// Create a key for the email login
	mailLoginKey := rdbkey.MailLogin(to)
	// Lock the email login key for 1 minute
	lock, err := v.locker.Obtain(context.Background(), rdbkey.MailLoginLock(), time.Minute, nil)
	// If the lock is not obtained, log an error
	if err == redislock.ErrNotObtained {
		log.Error("get mail login locker fail", "mailto", to, "err", err)
		return
	} else if err != nil {
		log.Error("get mail login locker occur a error", "mailto", to, "err", err)
		return
	}
	// Release the lock
	defer lock.Release(context.Background())
	// Set the verification code to the Redis database
	err = v.Redis.Set(context.Background(), mailLoginKey, code, 5*time.Minute).Err()
	// If the set operation fails, log an error
	if err != nil {
		log.Error("Set code to redis fail", "key", mailLoginKey, "verification_code", code, "err", err)
		return
	}
	// Send the verification code to the email channel
	emailChannel <- &VerificationMail{To: to, Code: code, ExpiredMinute: 5}
}

func (v *VerificationCode) Run() error {
	for {
		select {
		case email, ok := <-emailChannel:
			if !ok {
				log.Error("send email channel is closed")
				break
			}
			v.Mail.SendVerificationCode(email.To, email.Code, email.ExpiredMinute)
		case <-time.After(2 * time.Minute):
			log.Info("no user login email in 2 minutes, timeout!")
		}
	}
}
