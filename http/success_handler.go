package http

import (
	"encoding/json"
	"net/http"
	"orchid-starter/internal/common/model"

	bunker "github.com/yudhiana/bunker/errors"
)

func SuccessResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{
		Success: true,
		Status:  http.StatusOK,
		Data:    data,
	})
}

func ErrorResponse(w http.ResponseWriter, err error) {
	bunkerErr, ok := err.(*bunker.ApplicationError)
	if !ok {
		bunkerErr = bunker.New(bunker.StatusUnprocessableEntity).SetMessage(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(bunkerErr.HttpStatusCode)
	json.NewEncoder(w).Encode(model.Response{
		Success:      false,
		Status:       bunkerErr.HttpStatusCode,
		ErrorCode:    bunkerErr.ErrorCode,
		ErrorMessage: bunkerErr.Message,
	})
}
