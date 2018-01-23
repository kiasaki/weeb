package weeb

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// App represents a web application instance
type App struct {
	Log    *Logger
	Cache  Cache
	Config *Config
	Router *mux.Router
	DB     DB
	Tasks  *Tasks
}

// NewApp create a new App instance
func NewApp() *App {
	app := &App{}

	setupConfig(app)
	setupLog(app)
	setupCache(app)
	setupRouter(app)
	setupDatabase(app)
	setupTasks(app)

	return app
}

func setupLog(app *App) {
	context := L{}
	name := app.Config.Get("name", "")
	if name != "" {
		context["appName"] = name
	}
	app.Log = NewLogger().SetContext(context)
}

func setupCache(app *App) {
	app.Cache = NewMemoryCache()
}

func setupConfig(app *App) {
	app.Config = NewConfig()
	app.Config.LoadFromEnv()
}

func setupRouter(app *App) {
	app.Router = mux.NewRouter()

	app.Router.Use(handlers.RecoveryHandler(handlers.RecoveryLogger(app.Log)))
	app.Router.Use(app.loggingMiddleware)

	staticFilesHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	app.Router.PathPrefix("/static/").Handler(staticFilesHandler)
}

func setupDatabase(app *App) {
	dbUrl := app.Config.Get("databaseUrl", "postgres://weeb:weeb@localhost:5432/weeb?sslmode=disable")
	app.DB = NewPostgresDB(dbUrl, app.Log)
}

func setupTasks(app *App) {
	app.Tasks = NewTasks(app)

	app.Tasks.Register("start", func(app *App, _ []string) error {
		app.Start()
		return nil
	})
}

// Start starts the application
func (app *App) Start() {
	port := app.Config.Get("port", "3000")
	server := &http.Server{
		Addr:         "0.0.0.0:" + port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      app.Router,
	}

	go func() {
		app.Log.Info("started", L{"port": port})
		if err := server.ListenAndServe(); err != nil {
			app.Log.Fatal("failed to start", L{})
			os.Exit(1)
		}
	}()

	// Wait for Ctrl-C / SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Shutdown server or just exit after waiting 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	server.Shutdown(ctx)
	app.Log.Info("shutting down", L{})
	os.Exit(0)
}

// Run runs the application tasks. It looks at command line arguments to know
// which task to run
func (app *App) Run() {
	args := os.Args
	if len(args) <= 1 {
		args = []string{"", "start"}
	}
	app.Tasks.Run(args[1], args[2:])
}
