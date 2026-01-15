package common

import (
	"context"

	"orchid-starter/constants"

	modelCommon "orchid-starter/internal/common/model"

	"github.com/kataras/iris/v12"
)

// ExtractRequestContext extracts common headers from iris.Context and creates RequestContext
func ExtractRequestContext(irisCtx iris.Context) *modelCommon.RequestContext {
	return &modelCommon.RequestContext{
		RequestID: irisCtx.GetHeader(constants.HeaderRequestID),
	}
}

func SetRequestContext(ctx context.Context, irisCtx iris.Context) context.Context {
	reqCtx := ExtractRequestContext(irisCtx)
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
