package gqlHandler

import (
	"orchid-starter/gql/graph/generated"
	"orchid-starter/gql/graph/resolvers"
	"orchid-starter/internal/bootstrap"
	"orchid-starter/internal/common"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/kataras/iris/v12"
)

type graphHandler struct {
	di *bootstrap.DirectInjection
}

func NewGraphHandler(di *bootstrap.DirectInjection) *graphHandler {
	return &graphHandler{
		di: di,
	}
}

func NewDefaultGQLHandler(app chi.Router, di *bootstrap.DirectInjection) {

}

func (base *graphHandler) GQLHandler() iris.Handler {

	// directiveHandler := directive.NewDirective(base.di)
	conf := generated.Config{
		Resolvers: &resolvers.Resolver{
			DI: base.di,
		},
	}

	serverGraphql := handler.NewDefaultServer(generated.NewExecutableSchema(conf))
	return func(ctx iris.Context) {
		baseContext := ctx.Request().Context()
		serverGraphql.ServeHTTP(ctx.ResponseWriter(), ctx.Request().WithContext(common.SetRequestContext(baseContext, ctx)))
	}
}

func PlaygroundHandler() iris.Handler {
	h := playground.Handler("GraphQL Playground", "/gql/query")

	return func(ctx iris.Context) {
		h.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
	}
}
