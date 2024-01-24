package service

import (
	"fmt"
	"github.com/FlareZone/melon-backend/config"
	"net/smtp"
)

type MailService interface {
	SendVerificationCode(mailTo, code string, expired uint) bool
}

type GoogleMail struct {
	password string
	sender   string
}

type VerificationMail struct {
	To            string
	Code          string
	ExpiredMinute uint
}

func NewGoogleMail() MailService {
	return &GoogleMail{password: config.GoogleMail.Password, sender: config.GoogleMail.Sender}
}

func (m *GoogleMail) SendVerificationCode(mailTo, code string, expired uint) bool {
	// 设置 Gmail 的 SMTP 服务器详情
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"
	// 设置收件人地址、邮件主题和正文
	subject := "Your Verification Code!\n"
	body := fmt.Sprintf("Hello, this mail is sent from melon, your verification code is %s. It is valid for %d Minutes", code, expired)

	// 构建邮件内容
	message := []byte(subject + "\n" + body)

	// 使用 TLS 加密连接到 SMTP 服务器
	auth := smtp.PlainAuth("", m.sender, m.password, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, m.sender, []string{mailTo}, message)
	if err != nil {
		log.Error("smtp.SendMail failed", "mailTo", mailTo, "code", code, "err", err)
		return false
	}
	log.Info("email sent successfully", "mail_to", mailTo, "code", code)
	return true
}
