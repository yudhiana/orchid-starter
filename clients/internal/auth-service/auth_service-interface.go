package authService

type AuthServiceInterface interface {
	ValidateToken(token string) (bool, error)
}
