package http

import (
	"encoding/json"
	"net/http"
	"orchid-starter/internal/common/model"
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
