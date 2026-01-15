package applications

import (
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mataharibiz/ward/logging"
)

type chiApplication struct {
	app *chi.Mux
}

var chiAppInstance *chiApplication
var onceApp sync.Once

// // GetIrisApplication get iris Application instance
func GetChiApplication() *chi.Mux {
	onceIrisApp.Do(func() {
		log := logging.NewLogger()
		log.Info("Initialize chi application instance...")

		app := chi.NewRouter()
		// app.Use(middleware.RequestID)
		app.Use(middleware.Logger)

		chiAppInstance = &chiApplication{
			app: app,
		}
	})

	return chiAppInstance.app
}
