package bootstrap

import (
    "log"
    "os"
)

type Application struct {
    Env *Env
    // We remove the hard-coded DB field for now so we don't force a crash
}

func App() Application {
    app := &Application{}
    
    // 1. Load Env
    app.Env = NewEnv()
    if app.Env == nil {
        log.Println("WARNING: Env could not be loaded.")
    }

    // 2. Debugging: Print to the logs so we can see what's happening in Railway
    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        log.Println("CRITICAL: MONGO_URI is missing from environment variables!")
    } else {
        log.Printf("DEBUG: MONGO_URI is set, attempting connection logic...")
        // If you know the actual function name, you can call it here.
        // If not, leave this block commented out to break the crash loop.
    }

    return *app
}

func (app *Application) CloseDBConnection() {
    // No-op for now to prevent panic
    log.Println("CloseDBConnection called (noop)")
}