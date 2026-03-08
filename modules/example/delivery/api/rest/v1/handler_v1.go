package v1

import (
	"net/http"
	httpUtil "orchid-starter/internal/bootstrap/server/restful-server/http-util"
	"orchid-starter/internal/common"
	modelResponse "orchid-starter/modules/example/delivery/models/response"
	"orchid-starter/modules/example/usecase"
	openTelemetry "orchid-starter/observability/open-telemetry"

	bunker "github.com/yudhiana/bunker/errors"
)

type exampleHandler struct {
	usecase usecase.ExampleUsecaseInterface
	otel    *openTelemetry.OTel
}

func NewExampleHandler(u usecase.ExampleUsecaseInterface, otel *openTelemetry.OTel) *exampleHandler {
	return &exampleHandler{
		usecase: u,
		otel:    otel,
	}
}

func (base *exampleHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	ctx := common.SetRequestContext(r.Context(), r)

	ctx, span := base.otel.StartSpan(ctx, "handler", openTelemetry.GetFuncName())
	defer span.End()

	welcome := base.usecase.GetWelcome(ctx)
	httpUtil.SuccessResponse(w, modelResponse.WelcomeResponse{
		Message: welcome.Message,
	})
}

func (base *exampleHandler) ErrorResponse(w http.ResponseWriter, r *http.Request) {
	err := bunker.New(bunker.StatusBadRequest).SetMessage("Error occurred")
	httpUtil.ErrorResponse(w, err)
}
