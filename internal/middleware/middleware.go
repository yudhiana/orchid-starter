package middleware

import "net/http"

type Module func(http.Handler) http.Handler

type Base struct {
	Modules []Module
}
