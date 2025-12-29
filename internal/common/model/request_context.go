package model

import (
	"context"
	"strconv"

	"orchid-starter/constants"
)

// RequestContext contains common request metadata extracted from headers
type RequestContext struct {
	AppOrigin       string
	AppToken        string
	UserID          uint64
	CompanyID       uint64
	RequestID       string
	AppRequestID    string
	TokenState      string
	TokenIdentifier string
	ThirdParty      string
	Partner         string
}

// WithRequestContext adds RequestContext to the given context
func WithRequestContext(ctx context.Context, key any, value any) context.Context {
	return context.WithValue(ctx, value, key)
}

func SetRequestContext(ctx context.Context, reqCtx *RequestContext) context.Context {
	return context.WithValue(ctx, constants.RequestContextKey, reqCtx)
}

// GetRequestContext retrieves RequestContext from the given context
func GetRequestContext(ctx context.Context) (*RequestContext, bool) {
	reqCtx, ok := ctx.Value(constants.RequestContextKey).(*RequestContext)
	return reqCtx, ok
}

// GetUserIDString returns UserID as string
func (rc *RequestContext) GetUserIDString() string {
	if rc.UserID == 0 {
		return ""
	}
	return strconv.FormatUint(rc.UserID, 10)
}

// GetCompanyIDString returns CompanyID as string
func (rc *RequestContext) GetCompanyIDString() string {
	if rc.CompanyID == 0 {
		return ""
	}
	return strconv.FormatUint(rc.CompanyID, 10)
}

// IsEmpty checks if the request context has minimal required data
func (rc *RequestContext) IsEmpty() bool {
	return rc.AppOrigin == "" && rc.AppToken == ""
}
