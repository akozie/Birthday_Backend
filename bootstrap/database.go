package bootstrap

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func NewMySQLDatabase(env *Env) *sqlx.DB {
	if env == nil {
		log.Warn("MySQL init skipped: env is nil")
		return nil
	}

	if env.DBHost == "" || env.DBPort == "" || env.DBUser == "" || env.DBName == "" {
		log.Warn("MySQL config is incomplete; skipping database initialization")
		return nil
	}

	dsn := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?parseTime=true",
		env.DBUser,
		env.DBPass,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Errorf("failed to connect to MySQL: %v", err)
		return nil
	}

	log.Info("Successfully connected to MySQL")
	return db
}

func CloseMySqlConnection(client *sqlx.DB) {
	if client == nil {
		return
	}

	if err := client.Close(); err != nil {
		log.Errorf("failed to close MySQL connection: %v", err)
		return
	}

	log.Info("Connection to MySQL closed")
}
