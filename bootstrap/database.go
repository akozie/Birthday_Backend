package bootstrap

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func NewMySQLDatabase(env *Env) *sqlx.DB {
	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass
	dbName := env.DBName

	log.Infof("MySQL config: host='%s' port='%s' user='%s' name='%s'", dbHost, dbPort, dbUser, dbName)

	// example connection string: "test:test@(localhost:3306)/test"

	// If no DB host or credentials are provided assume MySQL is not used in this deployment.
	if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" {
		log.Info("MySQL config not provided, skipping MySQL initialization")
		return nil
	}

	db, err := sqlx.Connect("mysql", dbUser+":"+dbPass+"@("+dbHost+":"+dbPort+")/"+dbName+"?parseTime=true")
	if err != nil {
		log.Error("failed to connect to MySQL: ", err)
		return nil
	}

	return db
}

func CloseMySqlConnection(client *sqlx.DB) {
	if client == nil {
		package bootstrap

		import (
			_ "github.com/go-sql-driver/mysql"
			"github.com/jmoiron/sqlx"
			log "github.com/sirupsen/logrus"
		)

		func NewMySQLDatabase(env *Env) *sqlx.DB {
			dbHost := env.DBHost
			dbPort := env.DBPort
			dbUser := env.DBUser
			dbPass := env.DBPass
			dbName := env.DBName

			// example connection string: "test:test@(localhost:3306)/test"

			db, err := sqlx.Connect("mysql", dbUser+":"+dbPass+"@("+dbHost+":"+dbPort+")/"+dbName+"?parseTime=true")
			if err != nil {
				log.Fatal(err)
			}

			return db
		}

		func CloseMySqlConnection(client *sqlx.DB) {
			if client == nil {
				return
			}

			err := client.Close()
			if err != nil {
				log.Fatal(err)
			}

			log.Info("Connection to MySQL closed.")
		}
