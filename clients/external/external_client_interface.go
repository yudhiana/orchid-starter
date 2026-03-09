package externalClient

import (
	emailService "orchid-starter/clients/external/email-service"
)

type ExternalClientService struct {
	// TODO : add others service
	EmailService emailService.EmailServiceInterface
}

func ApplyExternalClient() *ExternalClientService {
	return &ExternalClientService{
		EmailService: emailService.GetEmailService(),
	}
}
