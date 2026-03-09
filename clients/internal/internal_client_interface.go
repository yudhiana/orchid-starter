package internalClient

import authService "orchid-starter/clients/internal/auth-service"

type InternalClientService struct {
	// TODO : add others service
	AuthService authService.AuthServiceInterface
}
