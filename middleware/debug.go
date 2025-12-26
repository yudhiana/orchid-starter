package middleware

import (
	"bytes"
	"io"
	"strings"

	"orchid-starter/internal/common"

	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/ward/logging"

	jsonSanitizer "github.com/mataharibiz/ward/json_sanitizer"
)

func Debug(ctx iris.Context) {
	if strings.ToUpper(common.GetEnvWithDefault("LOG_LEVEL", "DEBUG")) == "DEBUG" {
		method := ctx.Method()
		path := ctx.Path()
		uri := ctx.FullRequestURI()
		userAgent := ctx.GetHeader("User-Agent")
		ip := ctx.RemoteAddr()
		xRequestId := ctx.GetHeader("X-Request-ID")
		appRequestId := ctx.GetHeader("App-Request-ID")
		dmpOrigin := ctx.GetHeader("Dmp-Origin")
		dmpUserId := ctx.GetHeader("Dmp-User-ID")
		dmpUserCompanyId := ctx.GetHeader("Dmp-User-Company-ID")
		body := getBody(ctx)
		token := ctx.GetHeader("Authorization")
		logging.NewLogger().Debug("Incoming request",
			"action", "request",
			"dmp_origin", dmpOrigin,
			"dmp_user_id", dmpUserId,
			"dmp_user_company_id", dmpUserCompanyId,
			"method", method,
			"path", path,
			"uri", uri,
			"ip", ip,
			"user_agent", userAgent,
			"x_request_id", xRequestId,
			"app_request_id", appRequestId,
			"body", body,
			"token", token,
		)
	}
	ctx.Next()
}

func getBody(ctx iris.Context) (body any) {
	if ctx.Method() == iris.MethodPost || ctx.Method() == iris.MethodPut {

		r := ctx.Request()
		bodyReq, err := ctx.GetBody()
		if err != nil {
			logging.NewLogger().Error("Failed to read request body", "error", err)
			return nil
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyReq))

		defer func() {
			if r := recover(); r != nil {
				logging.NewLogger().Error("Failed to read request body", "error", r)
				body = string(bodyReq)
			}
		}()

		bodyRequest := common.CleanString(string(bodyReq))
		sanitizer := jsonSanitizer.NewJsonSanitizer()
		return sanitizer.Sanitize(bodyRequest)

	}
	return
}
