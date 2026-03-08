package gqlHandler

import (
	"net/http"
	"orchid-starter/gql/graph/generated"
	"orchid-starter/gql/graph/resolvers"
	"orchid-starter/internal/bootstrap/container"
	"orchid-starter/internal/common"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

type graphHandler struct {
	di *container.DirectInjection
}

func NewGraphHandler(di *container.DirectInjection) *graphHandler {
	return &graphHandler{
		di: di,
	}
}

func (base *graphHandler) GQLHandler() http.HandlerFunc {

	// directiveHandler := directive.NewDirective(base.di)
	conf := generated.Config{
		Resolvers: &resolvers.Resolver{
			DI: base.di,
		},
	}

	serverGraphql := handler.NewDefaultServer(generated.NewExecutableSchema(conf))
	return func(w http.ResponseWriter, r *http.Request) {
		baseContext := r.Context()
		serverGraphql.ServeHTTP(w, r.WithContext(common.SetRequestContext(baseContext, r)))
	}
}

func PlaygroundHandler() http.HandlerFunc {
	h := playground.Handler("GraphQL Playground", "/gql/query")

	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}
