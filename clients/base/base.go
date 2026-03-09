package base

import (
	"net/http"
	"orchid-starter/constants"
	"orchid-starter/internal/common"

	"github.com/go-resty/resty/v2"
)

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

func GetRestyClient() *resty.Client {
	return resty.New().SetTransport(NewTransport())
}
