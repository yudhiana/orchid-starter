package v2

import (
	"net/http"
	"orchid-starter/modules/default/usecase"
)

type defaultHandler struct {
	usecase usecase.DefaultUsecaseInterface
}

func NewDefaultHandler(u usecase.DefaultUsecaseInterface) *defaultHandler {
	return &defaultHandler{
		usecase: u,
	}
}

func (base *defaultHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	welcome := base.usecase.WelcomeUsecase(ctx)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(welcome))
}
