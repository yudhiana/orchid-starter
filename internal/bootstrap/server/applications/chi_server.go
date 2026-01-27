package applications

import (
	"sync"

	http "orchid-starter/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yudhiana/logos"
)

type chiApplication struct {
	app *chi.Mux
}

var chiAppInstance *chiApplication
var onceApp sync.Once

func GetChiApplication() *chi.Mux {
	onceApp.Do(func() {
		log := logos.NewLogger()
		log.Info("Initialize chi application instance...")

		app := chi.NewRouter()
		app.Use(middleware.RequestID)
		app.Use(middleware.Logger)
		app.NotFound(http.NotFoundHandler)
		chiAppInstance = &chiApplication{
			app: app,
		}
	})

	return chiAppInstance.app
}
