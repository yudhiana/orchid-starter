package clients

import (
	externalClient "orchid-starter/clients/external"
	internalClient "orchid-starter/clients/internal"
)

type Client struct {
	InternalClient *internalClient.InternalClientService
	ExternalClient *externalClient.ExternalClientService
}

func NewClient() *Client {
	return &Client{
		InternalClient: internalClient.ApplyInternalClient(),
		ExternalClient: externalClient.ApplyExternalClient(),
	}
}
