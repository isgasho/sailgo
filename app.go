package sailgo

import (
	"net/http"
)

type App struct {
	controllerRegister *Mux
	server             *http.Server
}

func NewApp() *App {
	return &App{
		controllerRegister: NewMux(),
		server:             &http.Server{},
	}
}

func (app *App) Add(path string, controller interface{}) {
	app.controllerRegister.RegisterRouter(path, controller)
}

func (app *App) Run(add string) {
	Info("app start to run")
	app.server.Addr = add
	app.server.Handler = app.controllerRegister
	app.server.ListenAndServe()
}
