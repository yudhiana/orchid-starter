package clients

import (
	externalClient "orchid-starter/internal/clients/external"
	internalClient "orchid-starter/internal/clients/internal"
)

type Client struct {
	InternalClient internalClient.InternalClientInterface
	ExternalClient *externalClient.ExternalClient
}

func NewClient() *Client {
	return &Client{
		InternalClient: internalClient.NewInternalClient(),
		ExternalClient: externalClient.ApplyExternalClient(),
	}
}
