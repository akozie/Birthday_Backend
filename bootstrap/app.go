package bootstrap

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Application struct {
	Env   *Env
	MySql *sqlx.DB
}

func App() Application {
	app := Application{}

	app.Env = NewEnv()
	if app.Env == nil {
		log.Warn("env could not be loaded")
		return app
	}

	app.MySql = NewMySQLDatabase(app.Env)
	if app.MySql == nil {
		log.Warn("MySQL connection is not available; continuing without a database client")
	}

	return app
}

func (app *Application) CloseDBConnection() {
	if app == nil {
		return
	}

	CloseMySqlConnection(app.MySql)
}
