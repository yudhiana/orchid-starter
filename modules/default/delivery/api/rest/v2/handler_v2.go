package v2

import (
	"net/http"
	response "orchid-starter/http"
	"orchid-starter/internal/common"
	modelResponse "orchid-starter/modules/default/delivery/models/response"
	"orchid-starter/modules/default/usecase"
	openTelemetri "orchid-starter/observability/open-telemetri"

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
	ctx := common.SetRequestContext(r.Context(), r)

	tp := openTelemetri.GetTraceProvider(ctx)
	ctx, span := tp.Tracer.Start(ctx, "defaultHandler.Welcome")
	defer span.End()

	welcome := base.usecase.GetWelcome(ctx)
	response.SuccessResponse(w, modelResponse.WelcomeResponse{
		Message: welcome.Message,
	})
}

func (base *defaultHandler) ErrorResponse(w http.ResponseWriter, r *http.Request) {
	err := bunker.New(bunker.StatusBadRequest).SetMessage("Error occurred")
	response.ErrorResponse(w, err)
}
