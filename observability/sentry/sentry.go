package sentry

import (
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"orchid-starter/internal/common"

	"github.com/getsentry/sentry-go"
	"github.com/yudhiana/logos"
)

func InitSentry() {
	// Init sentry DSN if sentry is enabled
	if strings.ToUpper(os.Getenv("SENTRY_ENABLE_ORCHID-STARTER")) == "TRUE" {
		dsn := os.Getenv("SENTRY_DSN_ORCHID-STARTER")
		appEnv := os.Getenv("APP_ENV")

		// Sentry Init V2
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              dsn,
			Environment:      appEnv,
			AttachStacktrace: true,
			Debug:            strings.ToUpper(os.Getenv("SENTRY_DEBUG")) == "TRUE",
		})

		if err != nil {
			logos.NewLogger().Error(err.Error())
		}

		logos.NewLogger().Info("Initializing Sentry", "dsn", dsn, "app_env", appEnv)
		defer sentry.Flush(2 * time.Second)
	}
}

// SentryLogger logs error to Sentry
func SentryLogger(err error, args ...any) {
	logos.NewLogger().Error("sentry error", "error", common.GetChainError(err))

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

			case *http.Request:
				scope.SetRequest(v)

			default:
				switch reflect.TypeOf(arg).Kind() {
				case reflect.Struct:
					scope.SetExtra("data", arg)
				}
			}
		}

		hub.CaptureException(err)
	})
}
