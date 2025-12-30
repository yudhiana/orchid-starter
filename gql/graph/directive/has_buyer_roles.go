package directive

import (
	"context"
	"fmt"
	"orchid-starter/constants"
	"orchid-starter/gql/graph/model"
	"orchid-starter/internal/bootstrap"
	"orchid-starter/internal/common"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	mbizUtil "github.com/mataharibiz/sange/v2/utils"
	gqlError "github.com/vektah/gqlparser/v2/gqlerror"
)

type HasBuyerRolesInput struct {
	Injection *bootstrap.DirectInjection
	Ctx       context.Context
}

func (d *Directive) HasBuyerRoles(ctx context.Context, obj any, next graphql.Resolver, roles []model.BuyerRoles) (res any, err error) {
	appOrigin := common.GetAppOriginFromContext(ctx)
	if appOrigin == "" {
		return nil, &gqlError.Error{
			Message: "unauthorized",
			Extensions: map[string]interface{}{
				"code": "INVALID_ORIGIN",
			},
		}
	}

	// Bypass if appOrigin is non-user including (system, tarrasque, worker)
	if constants.IsNonUser(appOrigin) {
		return next(ctx)
	}

	checkRole := &HasBuyerRolesInput{
		Injection: d.DI,
		Ctx:       ctx,
	}

	if exist := checkRole.CheckCacheRoles(roles); !exist {
		return nil, &gqlError.Error{
			Message: "you're not allowed perform this action",
			Extensions: map[string]interface{}{
				"code": "INVALID_SELLER_ROLES",
			},
		}
	}

	return next(ctx)
}

func (i *HasBuyerRolesInput) CheckCacheRoles(roles []model.BuyerRoles) bool {
	var (
		userID = common.GetUserIDFromContext(i.Ctx)
		compId = common.GetCompanyIDFromContext(i.Ctx)
	)

	rdsUtil, err := mbizUtil.NewRedisUtil(constants.RedisTypeRole)
	if err != nil {
		i.Injection.Log.Warn("failed connect to redis", "error", err)
		return false
	}

	redisKey := fmt.Sprintf("COMPANY:%v:USER:ROLES:%v", compId, userID)
	sellerRolesString, errGet := rdsUtil.HGet(redisKey, "buyer-roles")
	if errGet != nil || sellerRolesString == "" {
		i.Injection.Log.Warn("role buyer on redis is empty")
		return false
	}

	if _, ok := sellerRolesString.(string); !ok {
		i.Injection.Log.Warn("role buyer on redis is not string")
		return false
	}

	sellerRoles := strings.Split(sellerRolesString.(string), ",")
	baseRoles := make(map[string]bool)

	for _, role := range sellerRoles {
		baseRoles[strings.ToUpper(role)] = true
	}

	for _, role := range roles {
		if _, ok := baseRoles[strings.ToUpper(role.String())]; !ok {
			return false
		}
	}
	return true
}
