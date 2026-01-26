package model

type Response struct {
	Status    int    `json:"status"`
	Success   bool   `json:"success"`
	Data      any    `json:"data,omitempty"`
	Error     error  `json:"error,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
}
