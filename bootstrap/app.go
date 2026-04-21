package bootstrap

import log "github.com/sirupsen/logrus"

type Application struct {
	Env *Env
}

func App() Application {
	app := Application{}

	app.Env = NewEnv()
	if app.Env == nil {
		log.Warn("env could not be loaded")
	}

	return app
}

func (app *Application) CloseDBConnection() {
}
