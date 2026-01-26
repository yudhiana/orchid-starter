package common

import (
	"context"
	"net/http"

	"orchid-starter/constants"

	modelCommon "orchid-starter/internal/common/model"
)

// ExtractRequestContext extracts common headers from http.Request and creates RequestContext
func ExtractRequestContext(r *http.Request) *modelCommon.RequestContext {
	return &modelCommon.RequestContext{
		RequestID: r.Header.Get(constants.HeaderRequestID),
	}
}

func SetRequestContext(ctx context.Context, r *http.Request) context.Context {
	reqCtx := ExtractRequestContext(r)
	return modelCommon.SetRequestContext(ctx, reqCtx)
}

// GetRequestContext creates a new context with Context RequestID attached
func GetRequestContext(ctx context.Context) context.Context {
	requestId := GetRequestIDFromContext(ctx)
	return modelCommon.WithRequestContext(ctx, constants.HeaderRequestID, requestId)
}

// GetRequestIDFromContext extracts request ID from context for tracing
func GetRequestIDFromContext(ctx context.Context) string {
	if reqCtx, ok := modelCommon.GetRequestContext(ctx); ok {
		return reqCtx.RequestID
	}
	return ""
}
