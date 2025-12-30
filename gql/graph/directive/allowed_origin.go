package directive

import (
	"context"
	"orchid-starter/gql/graph/model"
	"orchid-starter/internal/common"
	"slices"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	gqlError "github.com/vektah/gqlparser/v2/gqlerror"
)

func (d *Directive) AllowedOrigin(ctx context.Context, obj any, next graphql.Resolver, origin []model.Origin) (res any, err error) {
	appOrigin := common.GetAppOriginFromContext(ctx)
	if appOrigin == "" {
		return nil, &gqlError.Error{
			Message: "unauthorized",
			Extensions: map[string]interface{}{
				"code": "INVALID_ORIGIN",
			},
		}
	}

	exist := slices.ContainsFunc(origin, func(item model.Origin) bool {
		return item.String() == strings.ToUpper(appOrigin)
	})
	if !exist {
		return nil, &gqlError.Error{
			Message: "you're not allowed perform this action",
			Extensions: map[string]interface{}{
				"code": "INVALID_ORIGIN",
			},
		}
	}

	return
}
