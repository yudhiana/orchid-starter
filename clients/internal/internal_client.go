package internalClient

import (
	authService "orchid-starter/clients/internal/auth-service"
)

func ApplyInternalClient() *InternalClientService {
	return &InternalClientService{
		AuthService: authService.NewAuthService(),
	}
}
