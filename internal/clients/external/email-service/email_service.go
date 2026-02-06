package emailService

import (
	"orchid-starter/pkg/gomail"
)

type EmailServiceInterface interface {
	SendEmail(subject string, body string, to, cc, bcc []string) error
}

func GetEmailService() EmailServiceInterface {
	return gomail.NewService()
}
