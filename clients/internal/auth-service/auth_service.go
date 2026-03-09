package authService

import (
	"orchid-starter/clients/base"
	"os"
	"strings"
)

// InternalClient provides internal API client functionality
type Client struct {
	debug bool
}

func NewAuthService() *Client {
	return &Client{
		debug: strings.ToUpper(os.Getenv("LOG_LEVEL")) == "DEBUG",
	}
}

func (c *Client) ValidateToken(token string) (bool, error) {
	restyClient := base.GetRestyClient()
	restyClient.SetDebug(c.debug)
	return true, nil
}
