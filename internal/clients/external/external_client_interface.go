package externalClient

import (
	emailService "orchid-starter/internal/clients/external/email-service"
)

type ExternalClient struct {
	// TODO : add others service
	EmailService emailService.EmailServiceInterface
}

func ApplyExternalClient() *ExternalClient {
	return &ExternalClient{
		EmailService: emailService.GetEmailService(),
	}
}
