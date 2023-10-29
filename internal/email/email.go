package mail

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

// MailSender used to send email.
type MailSender struct {
	fromAddr string
	smtpHost string
	smtpPort int
	smtpScrt string
}

// NewMailSender creates a MailSender then return.
func NewMailSender(smtpHost, smtpAddr, smtpScrt string, smtpPort int) *MailSender {
	return &MailSender{
		fromAddr: smtpAddr,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		smtpScrt: smtpScrt,
	}
}

// SendEmail sends an email to given addr.
func (ms *MailSender) SendEmail(toAddr []string, sub, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", ms.fromAddr)
	m.SetHeader("To", toAddr...)
	m.SetHeader("Subject", sub)
	m.SetBody("text/plain", body)
	d := gomail.NewDialer(ms.smtpHost, ms.smtpPort, ms.fromAddr, ms.smtpScrt)
	err := d.DialAndSend(m)
	if err != nil {
		return fmt.Errorf("cannot send email: %w", err)
	}
	return nil
}

// SendActivateEmail sends an activation email to addr
func (ms *MailSender) SendActivateEmail(toAddr, nickname, code, ip string) error {
	content := fmt.Sprintf("尊敬的%s：\n你好，你在 ShortVideo 应用注册了账户，"+
		"这是你的激活码： %s\n，你注册时使用的IP地址为 %s，请确保收到该邮件后的 10分钟 内激活该账户，"+
		"否则本次注册失败。", nickname, code, ip)
	return ms.SendEmail([]string{toAddr}, "ShortVideo 账户激活", content)
}
