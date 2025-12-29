package directive

import (
	"context"
	"orchid-starter/constants"
	"orchid-starter/internal/common"

	"github.com/99designs/gqlgen/graphql"
	gqlError "github.com/vektah/gqlparser/v2/gqlerror"
)

func AuthToken(ctx context.Context, obj any, next graphql.Resolver) (res any, err error) {

	var (
		appToken   = common.GetAppTokenFromContext(ctx)
		appOrigin  = common.GetAppOriginFromContext(ctx)
		tokenState = common.GetTokenStateFromContext(ctx)
		userID     = common.GetUserIDFromContext(ctx)
		companyID  = common.GetCompanyIDFromContext(ctx)
	)

	if isThirdParty := common.GetThirdPartyFromContext(ctx); isThirdParty == "1" {
		if partner := common.GetPartnerFromContext(ctx); partner != "" {

			// handle third party authentication
			// TODO: implement third party auth validation
		}
	}

	if appOrigin == "" || appToken == "" {
		return nil, &gqlError.Error{
			Message: "app origin Or app token header is empty",
			Extensions: map[string]interface{}{
				"code": "INVALID_TOKEN",
			},
		}
	}

	// Bypass non-user token include (system,admin,worker)
	if constants.IsNonUser(appOrigin) {
		return next(ctx)
	}

	// Validate user token
	if userID == 0 || (companyID == 0 && tokenState == "PRE-LOGIN") {
		return nil, &gqlError.Error{
			Message: "user id Or company id is invalid",
			Extensions: map[string]interface{}{
				"code": "INVALID_TOKEN",
			},
		}
	}

	// valid session token
	// TODO: implement session token validation

	return next(ctx)
}
