package v2

import (
	"net/http"
	response "orchid-starter/http"
	"orchid-starter/modules/default/usecase"

	bunker "github.com/yudhiana/bunker/errors"
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
	response.SuccessResponse(w, welcome)
}

func (base *defaultHandler) ErrorResponse(w http.ResponseWriter, r *http.Request) {
	err := bunker.New(bunker.StatusBadRequest).SetMessage("Error occurred")
	response.ErrorResponse(w, err)
}
