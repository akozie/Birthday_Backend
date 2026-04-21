package bootstrap

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPass                 string `mapstructure:"DB_PASS"`
	DBName                 string `mapstructure:"DB_NAME"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	GoogleClientID         string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret     string `mapstructure:"GOOGLE_CLIENT_SECRET"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Allow environment variables to override or replace missing .env values.
	// This matters in deployment environments where no .env file exists.
	for _, key := range []string{
		"APP_ENV",
		"SERVER_ADDRESS",
		"CONTEXT_TIMEOUT",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASS",
		"DB_NAME",
		"ACCESS_TOKEN_EXPIRY_HOUR",
		"REFRESH_TOKEN_EXPIRY_HOUR",
		"ACCESS_TOKEN_SECRET",
		"REFRESH_TOKEN_SECRET",
		"GOOGLE_CLIENT_ID",
		"GOOGLE_CLIENT_SECRET",
	} {
		_ = viper.BindEnv(key)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No .env file found")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	// Fall back to raw environment variables for any values Viper did not load.
	// This keeps deployment configs working even if the config file is missing.
	if env.AppEnv == "" {
		env.AppEnv = os.Getenv("APP_ENV")
	}
	if env.ServerAddress == "" {
		env.ServerAddress = os.Getenv("SERVER_ADDRESS")
	}
	if env.ContextTimeout == 0 {
		if timeout, err := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT")); err == nil {
			env.ContextTimeout = timeout
		}
	}
	if env.DBHost == "" {
		env.DBHost = os.Getenv("DB_HOST")
	}
	if env.DBPort == "" {
		env.DBPort = os.Getenv("DB_PORT")
	}
	if env.DBUser == "" {
		env.DBUser = os.Getenv("DB_USER")
	}
	if env.DBPass == "" {
		env.DBPass = os.Getenv("DB_PASS")
	}
	if env.DBName == "" {
		env.DBName = os.Getenv("DB_NAME")
	}
	if env.AccessTokenSecret == "" {
		env.AccessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
	}
	if env.RefreshTokenSecret == "" {
		env.RefreshTokenSecret = os.Getenv("REFRESH_TOKEN_SECRET")
	}
	if env.GoogleClientID == "" {
		env.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	}
	if env.GoogleClientSecret == "" {
		env.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	}

	if env.AppEnv == "development" {
		log.Info("The App is running in development env")
	}

	return &env
}
