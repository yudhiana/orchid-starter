package gomail

import (
	"crypto/tls"
	"fmt"
	"orchid-starter/internal/common"
	"os"

	"github.com/yudhiana/logos"
	gomailV2 "gopkg.in/gomail.v2"
)

type mailer struct {
	Message     *gomailV2.Message
	FromAddress string
	FromName    string
}

type Config struct {
	SmtpHost string
	SmtpPort int
	Username string
	Password string
}

func NewService() *mailer {
	return &mailer{
		Message:     gomailV2.NewMessage(),
		FromAddress: common.GetEnvWithDefault("EMAIL_FROM_ADDRESS", "no-reply@mail.loc"),
		FromName:    common.GetEnvWithDefault("EMAIL_FROM_NAME", "orchid-starter"),
	}
}

func (m *mailer) SendEmail(subject string, body string, to, cc, bcc []string) error {
	from := m.Message.FormatAddress(m.FromAddress, m.FromName)
	m.Message.SetHeader("From", from)
	m.Message.SetHeader("To", to...)
	m.Message.SetHeader("Subject", subject)
	m.Message.SetBody("text/html", body)

	if len(cc) > 0 {
		m.Message.SetHeader("Cc", cc...)
	}

	if len(bcc) > 0 {
		m.Message.SetHeader("Bcc", bcc...)
	}

	config := m.GetConfig()
	dialer := gomailV2.NewDialer(
		config.SmtpHost,
		config.SmtpPort,
		config.Username,
		config.Password,
	)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	logos.NewLogger().Info(fmt.Sprintf("email sent to %v", to))
	return dialer.DialAndSend(m.Message)
}

func (m *mailer) GetConfig() Config {
	return Config{
		SmtpHost: common.GetEnvWithDefault("SMTP_SERVER", "mailcatcher"),
		SmtpPort: common.GetIntEnv("SMTP_PORT", 25),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}
