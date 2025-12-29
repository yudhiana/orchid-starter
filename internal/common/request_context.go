package common

import (
	"context"
	"strconv"
	"strings"

	"orchid-starter/constants"

	modelCommon "orchid-starter/internal/common/model"

	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/sange/v2"
)

// ExtractRequestContext extracts common headers from iris.Context and creates RequestContext
func ExtractRequestContext(irisCtx iris.Context) *modelCommon.RequestContext {
	userID, _ := strconv.ParseUint(irisCtx.GetHeader(constants.HeaderUserID), 10, 64)
	companyID, _ := strconv.ParseUint(irisCtx.GetHeader(constants.HeaderCompanyID), 10, 64)
	authHeader := irisCtx.GetHeader(constants.HeaderAuthorization)
	if strings.HasPrefix(authHeader, sange.OAUTH) {
		authHeader = strings.TrimPrefix(authHeader, sange.OAUTH)
		authHeader = strings.TrimSpace(authHeader)
	}

	return &modelCommon.RequestContext{
		AppOrigin:       irisCtx.GetHeader(constants.HeaderAppOrigin),
		AppToken:        authHeader,
		UserID:          userID,
		CompanyID:       companyID,
		RequestID:       irisCtx.GetHeader(constants.HeaderRequestID),
		AppRequestID:    irisCtx.GetHeader(constants.HeaderAppRequestID),
		TokenState:      irisCtx.GetHeader(constants.HeaderTokenState),
		TokenIdentifier: irisCtx.GetHeader(constants.HeaderTokenIdentifier),
		ThirdParty:      irisCtx.GetHeader(constants.HeaderThirdParty),
		Partner:         irisCtx.GetHeader(constants.HeaderPartner),
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

// GetAppOriginFromContext extracts app origin from context
func GetAppOriginFromContext(ctx context.Context) string {
	if reqCtx, ok := modelCommon.GetRequestContext(ctx); ok {
		return reqCtx.AppOrigin
	}
	return ""
}

// GetAppTokenFromContext extracts app token from context
func GetAppTokenFromContext(ctx context.Context) string {
	if reqCtx, ok := modelCommon.GetRequestContext(ctx); ok {
		return reqCtx.AppToken
	}
	return ""
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) uint64 {
	if reqCtx, ok := modelCommon.GetRequestContext(ctx); ok {
		return reqCtx.UserID
	}
	return 0
}

// GetCompanyIDFromContext extracts company ID from context
func GetCompanyIDFromContext(ctx context.Context) uint64 {
	if reqCtx, ok := modelCommon.GetRequestContext(ctx); ok {
		return reqCtx.CompanyID
	}
	return 0
}

// GetRequestIDFromContext extracts request ID from context for tracing
func GetRequestIDFromContext(ctx context.Context) string {
	if reqCtx, ok := modelCommon.GetRequestContext(ctx); ok {
		return reqCtx.RequestID
	}
	return ""
}

func GetAppRequestIDFromContext(ctx context.Context) string {
	if reqCtx, ok := modelCommon.GetRequestContext(ctx); ok {
		return reqCtx.AppRequestID
	}
	return ""
}

// GetTokenStateFromContext extracts token state from context
func GetTokenStateFromContext(ctx context.Context) string {
	if reqCtx, ok := modelCommon.GetRequestContext(ctx); ok {
		return reqCtx.TokenState
	}
	return ""
}

// GetThirdPartyFromContext extracts third party from context
func GetThirdPartyFromContext(ctx context.Context) string {
	if reqCtx, ok := modelCommon.GetRequestContext(ctx); ok {
		return reqCtx.ThirdParty
	}
	return ""
}

// GetPartnerFromContext extracts partner from context
func GetPartnerFromContext(ctx context.Context) string {
	if reqCtx, ok := modelCommon.GetRequestContext(ctx); ok {
		return reqCtx.Partner
	}
	return ""
}

// GetTokenIdentifierFromContext extracts token identifier from context
func GetTokenIdentifierFromContext(ctx context.Context) string {
	if reqCtx, ok := modelCommon.GetRequestContext(ctx); ok {
		return reqCtx.TokenIdentifier
	}
	return ""
}
