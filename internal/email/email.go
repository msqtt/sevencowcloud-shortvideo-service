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
func (ms *MailSender) SendActivateEmail(toAddr, code, ip string) error {
	content := fmt.Sprintf("尊敬的用户：\n您好，您在 ShortVideo 应用注册了账户，"+
		"您的验证码是：%s，您注册时使用的IP地址为 %s，请确保收到该邮件后的 10分钟 内注册该账户，"+
		"否则本次注册失败。", code, ip)
	return ms.SendEmail([]string{toAddr}, "ShortVideo 账户注册", content)
}
