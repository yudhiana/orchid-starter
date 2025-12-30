package middleware

import (
	"context"
	"fmt"
	"orchid-starter/constants"
	"orchid-starter/internal/common"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/sange/v2"
	mbizUtil "github.com/mataharibiz/sange/v2/utils"
	"github.com/mataharibiz/ward/logging"
)

type Roles struct {
	AllowedBuyerRoles  []string
	AllowedSellerRoles []string
}

func RolesAllowed() *Roles {
	return &Roles{}
}

func (role *Roles) SetBuyer(roleName ...string) *Roles {
	role.AllowedBuyerRoles = roleName
	return role
}

func (role *Roles) SetSeller(roleName ...string) *Roles {
	role.AllowedSellerRoles = roleName
	return role
}

func (role *Roles) Validate(irisCtx iris.Context) *sange.Error {
	ctx := common.SetRequestContext(irisCtx.Request().Context(), irisCtx)

	// Bypass if appOrigin is non-user including (system, tarrasque, worker)
	if constants.IsNonUser(common.GetAppOriginFromContext(ctx)) {
		return nil
	}

	if exist := role.checkCacheRoles(ctx); !exist {
		return sange.SetError(sange.Forbidden, nil, "user role not allowed")
	}

	return nil
}

func (roles *Roles) checkCacheRoles(ctx context.Context) bool {
	var (
		userID = common.GetUserIDFromContext(ctx)
		compId = common.GetCompanyIDFromContext(ctx)
	)

	rdsUtil, err := mbizUtil.NewRedisUtil(constants.RedisTypeRole)
	if err != nil {
		logging.NewLogger().Warn("failed connect to redis", "error", err)
		return false
	}

	redisKey := fmt.Sprintf("COMPANY:%v:USER:ROLES:%v", compId, userID)
	if len(roles.AllowedBuyerRoles) > 0 {
		buyerRolesString, errGet := rdsUtil.HGet(redisKey, "buyer-roles")
		if errGet != nil || buyerRolesString == "" {
			logging.NewLogger().Warn("role buyer on redis is empty")
			return false
		}

		if _, ok := buyerRolesString.(string); !ok {
			logging.NewLogger().Warn("role buyer on redis is not string")
			return false
		}

		if exist := roles.checkCacheRoleBuyer(strings.Split(buyerRolesString.(string), ",")); !exist {
			return false
		}
	}

	if len(roles.AllowedSellerRoles) > 0 {
		sellerRolesString, errGet := rdsUtil.HGet(redisKey, "seller-roles")
		if errGet != nil || sellerRolesString == "" {
			logging.NewLogger().Warn("role buyer on redis is empty")
			return false
		}

		if _, ok := sellerRolesString.(string); !ok {
			logging.NewLogger().Warn("role seller on redis is not string")
			return false
		}

		if exist := roles.checkCacheRoleBuyer(strings.Split(sellerRolesString.(string), ",")); !exist {
			return false
		}
	}

	return true

}

func (roles *Roles) checkCacheRoleBuyer(r []string) bool {
	baseRoles := make(map[string]bool)
	for _, role := range r {
		baseRoles[strings.ToUpper(role)] = true
	}

	for _, role := range roles.AllowedBuyerRoles {
		if _, ok := baseRoles[strings.ToUpper(role)]; !ok {
			return false
		}
	}
	return true
}

func (roles *Roles) checkCacheRoleSeller(r []string) bool {
	baseRoles := make(map[string]bool)
	for _, role := range r {
		baseRoles[strings.ToUpper(role)] = true
	}

	for _, role := range roles.AllowedSellerRoles {
		if _, ok := baseRoles[strings.ToUpper(role)]; !ok {
			return false
		}
	}
	return true
}
