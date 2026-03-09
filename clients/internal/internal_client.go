package internalClient

import (
	authService "orchid-starter/clients/internal/auth-service"
)

func ApplyInternalClient() *InternalClientService {
	return &InternalClientService{
		// TODO : add others service
		AuthService: authService.NewAuthService(),
	}
}
