package bootstrap

import (
    "log"
    "os"
)

type Application struct {
    Env *Env
    // Don't rely on sqlx if you aren't using MySQL
}

func App() Application {
    app := &Application{}
    app.Env = NewEnv()

    // DEBUG: Print what we are reading
    mongoURI := os.Getenv("MONGO_URI") 
    log.Printf("DEBUG: The MONGO_URI variable is: '%s'", mongoURI)

    if mongoURI == "" {
        log.Fatal("FATAL: MONGO_URI is empty! Check your Railway Environment Variables.")
    }

    return *app
}