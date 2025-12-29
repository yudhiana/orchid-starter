package sentry

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"orchid-starter/internal/common"

	"github.com/getsentry/sentry-go"
	"github.com/kataras/iris/v12"
	"github.com/mataharibiz/sange/v2"
	"github.com/mataharibiz/ward/logging"
)

func InitSentry() {
	// Init sentry DSN if sentry is enabled
	if strings.ToUpper(os.Getenv("SENTRY_ENABLE_orchid-starter")) == "TRUE" {
		dsn := os.Getenv("SENTRY_DSN_orchid-starter")
		appEnv := os.Getenv("APP_ENV")

		// Sentry Init V2
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              dsn,
			Environment:      appEnv,
			AttachStacktrace: true,
			Debug:            strings.ToUpper(os.Getenv("SENTRY_DEBUG")) == "TRUE",
		})

		if err != nil {
			logging.NewLogger().Error(err.Error())
		}

		logging.NewLogger().Info("Initializing Sentry", "dsn", dsn, "app_env", appEnv)
		defer sentry.Flush(2 * time.Second)
	}
}

// SentryLogger logs error to Sentry
func SentryLogger(err error, args ...any) {
	if skippedSentSentry(err) {
		return
	}

	logging.NewLogger().Error("sentry error", "error", common.GetChainError(err))

	hub := sentry.CurrentHub().Clone()

	//set extra context to Sentry
	hub.WithScope(func(scope *sentry.Scope) {
		// Set Sentry Error Level
		scope.SetLevel(sentry.LevelError)
		for _, arg := range args {
			switch v := arg.(type) {
			case map[string]any:
				for idx, row := range arg.(map[string]any) {
					scope.SetExtra(idx, row)
				}
			case iris.Context:
				scope.SetRequest(v.Request())
			case *http.Request:
				scope.SetRequest(v)
			case context.Context:
				if requestID := common.GetRequestIDFromContext(v); requestID != "" {
					scope.SetContext("app_context", map[string]any{
						"request_id":     requestID,
						"app_request_id": common.GetAppRequestIDFromContext(v),
						"app_origin":     common.GetAppOriginFromContext(v),
						"user_id":        common.GetUserIDFromContext(v),
						"company_id":     common.GetCompanyIDFromContext(v),
					})
				}
			default:
				switch reflect.TypeOf(arg).Kind() {
				case reflect.Struct:
					scope.SetExtra("data", arg)
				}
			}
		}

		hub.CaptureException(err)
		scope.SetFingerprint([]string{sanitizeErrorMessage(err.Error())})
		if ok, sangeErr := common.IsSangeError(err); ok {
			if sangeErr.ExtraData != nil {
				for idx, row := range sangeErr.ExtraData {
					scope.SetExtra(idx, row)
				}
			}
		}

	})
}

// regex number and UUID cleanup
var cleanupRegex = regexp.MustCompile(`(?i)\b\d+\b|[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}|[a-f0-9]{24}`)

func sanitizeErrorMessage(msg string) string {
	// Remove numbers and UUIDs
	clean := cleanupRegex.ReplaceAllString(msg, "")

	// Regex to match arrays of JSON objects
	arrayRegex := regexp.MustCompile(`\[(\{[^}]+}(?:,\{[^}]+})+)]`)
	clean = arrayRegex.ReplaceAllStringFunc(clean, func(match string) string {
		// Remove brackets and split objects
		content := match[1 : len(match)-1]
		objects := strings.Split(content, "},{")
		for i := range objects {
			if i > 0 {
				objects[i] = "{" + objects[i]
			}
			if i < len(objects)-1 {
				objects[i] = objects[i] + "}"
			}
		}
		// Check if all objects are the same
		first := objects[0]
		for _, obj := range objects[1:] {
			if obj != first {
				return match // Not all the same, return original
			}
		}
		return "[" + first + "]"
	})

	return strings.Join(strings.Fields(clean), " ")
}

func skippedSentSentry(err error) bool {
	if ok, sangeAppError := common.IsSangeError(err); ok {
		code := sangeAppError.Code
		for _, whiteCode := range whiteListCode() {
			if sange.ErrorIs(whiteCode, code) {
				return true
			}
		}
	}
	return false
}

func whiteListCode() []sange.ErrorStatus {
	return []sange.ErrorStatus{
		sange.NotAuthorized,
		sange.Forbidden,
		sange.NotFound,
	}
}
