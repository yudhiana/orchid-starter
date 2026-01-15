package gqlHandler

import (
	"orchid-starter/gql/graph/directive"
	"orchid-starter/gql/graph/generated"
	"orchid-starter/gql/graph/resolvers"
	"orchid-starter/http"
	"orchid-starter/internal/bootstrap"
	"orchid-starter/internal/common"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kataras/iris/v12"
	promHttp "github.com/prometheus/client_golang/prometheus/promhttp"
)

type graphHandler struct {
	di *bootstrap.DirectInjection
}

func NewGraphHandler(di *bootstrap.DirectInjection) *graphHandler {
	return &graphHandler{
		di: di,
	}
}

func NewDefaultGQLHandler(app iris.Party, di *bootstrap.DirectInjection) {

	// defaultRepository := repository.NewDefaultRepository(di.GetMySQL(), di.GetElasticsearch())

	// Get the comprehensive client for all API operations
	// client := di.GetClient()

	// Initialize usecase with client access
	// defaultUseCase := usecase.NewDefaultUsecase(di.GetMySQL(), defaultRepository, client)
	// defaultV2 := v2.NewDefaultHandler(defaultUseCase)

	app.Get("/metrics", iris.FromStd(promHttp.Handler()))
	// app.Get("/", defaultV2.Welcome)
	// app.Get("/health-check", http.HealthCheckHandler)
	app.OnErrorCode(iris.StatusNotFound, http.NotFoundHandler)
}

func (base *graphHandler) GQLHandler() iris.Handler {

	directiveHandler := directive.NewDirective(base.di)
	conf := generated.Config{
		Resolvers: &resolvers.Resolver{
			DI: base.di,
		},
		Directives: generated.DirectiveRoot{
			AuthToken:      directiveHandler.AuthToken,
			AllowedOrigin:  directiveHandler.AllowedOrigin,
			HasSellerRoles: directiveHandler.HasSellerRoles,
			HasBuyerRoles:  directiveHandler.HasBuyerRoles,
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
