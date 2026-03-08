package internalClient

import (
	"net/http"
	"os"
	"strings"

	"orchid-starter/constants"
	"orchid-starter/internal/common"

	"github.com/go-resty/resty/v2"
)

// InternalClient provides internal API client functionality
type InternalClient struct {
	TokenSystem     string
	AppOriginSystem string
	Debug           bool
}

type requestContextTransport struct {
	base http.RoundTripper
}

func NewTransport() *requestContextTransport {
	return &requestContextTransport{
		base: http.DefaultTransport,
	}
}

func (t *requestContextTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if reqCtx := req.Context().Value(constants.RequestContextKey); reqCtx != nil {
		req.Header.Set(constants.HeaderRequestID, common.GetRequestIDFromContext(req.Context()))
	}
	return t.base.RoundTrip(req)
}

func NewInternalClient() InternalClientInterface {
	return &InternalClient{
		TokenSystem:     os.Getenv("SYSTEM_TOKEN"),
		AppOriginSystem: "system",
		Debug:           strings.ToUpper(os.Getenv("LOG_LEVEL")) == "DEBUG",
	}
}

func GetRestyClient() *resty.Client {
	return resty.New().SetTransport(NewTransport())
}
