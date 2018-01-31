package weeb

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/markbates/refresh/refresh"
	refreshweb "github.com/markbates/refresh/refresh/web"
)

// App represents a web application instance
type App struct {
	Log       *Logger
	Config    *Config
	Router    *Router
	Session   *Session
	templates *template.Template

	Cache Cache
	DB    DB

	Migrations *MigrationRunner
	Tasks      *TaskRunner
	Auth       *Auth
}

// NewApp create a new App instance
func NewApp() *App {
	app := &App{}

	setupTasks(app)
	setupConfig(app)
	setupLog(app)
	setupCache(app)
	setupRouter(app)
	setupTemplates(app)
	setupDatabase(app)
	setupMigrations(app)
	setupAuth(app)

	addWeebMigrationsToApp(app)

	return app
}

func setupTasks(app *App) {
	app.Tasks = NewTaskRunner(app)

	app.Tasks.Register("start", func(app *App, _ []string) error {
		app.Start()
		return nil
	})

	app.Tasks.Register("dev", func(app *App, _ []string) error {
		app.Dev()
		return nil
	})

	app.Tasks.Register("generate-session-key", func(app *App, _ []string) error {
		fmt.Println(generateRandomKey(64))
		return nil
	})
}

func setupConfig(app *App) {
	app.Config = NewConfig()

	// Most config comes from APP_ env vars, but let's load
	// PORT and DEBUG without prefixes as convenience
	app.Config.Set("port", os.Getenv("PORT"))
	app.Config.Set("port", app.Config.Get("port", "3000"))
	app.Config.Set("debug", os.Getenv("DEBUG"))
	app.Config.Set("debug", app.Config.Get("debug", "1"))

	app.Config.LoadFromEnv()

	app.Tasks.Register("config", func(app *App, _ []string) error {
		fmt.Println()
		fmt.Print(displayMap(app.Config.Values(), 2, 2))
		fmt.Println()
		return nil
	})
}

func setupLog(app *App) {
	context := L{}
	name := app.Config.Get("name", "")
	if name != "" {
		context["app"] = name
	}
	app.Log = NewLogger().SetContext(context)
}

func setupCache(app *App) {
	app.Cache = NewMemoryCache()
}

func setupRouter(app *App) {
	app.Router = NewRouter(app)

	app.Router.Use(recoverMiddleware)
	app.Router.Use(loggingMiddleware)
	if app.Config.GetBool("dev") {
		app.Router.UseHTTP(refreshweb.ErrorChecker)
	}

	app.Router.Static("/static/", "static")
}

func setupTemplates(app *App) {
	if dirExists("templates") {
		var err error
		app.templates, err = template.ParseGlob("templates/*")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, t := range app.templates.Templates() {
			name := t.Name()
			ext := filepath.Ext(name)
			app.templates, _ = app.templates.AddParseTree(name[:len(name)-len(ext)], t.Tree)
		}
	} else {
		// Make sure tempaltes is a valid instance of *tempalte.Template
		app.templates = template.New("weeb")
	}
}

func setupDatabase(app *App) {
	dbURL := app.Config.Get("databaseUrl", "postgres://postgres:postgres@localhost:5432/app?sslmode=disable")
	app.DB = NewPostgresDB(dbURL, app.Log)
}

func setupMigrations(app *App) {
	app.Migrations = NewMigrationRunner(app)

	app.Tasks.Register("migrate", migrationRunnerTask)
}

func setupAuth(app *App) {
	app.Auth = NewAuth(app)
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

// Dev runs the application in dev mode rebuilding the current directory on file changes
func (app *App) Dev() {
	config := &refresh.Configuration{
		AppRoot:            ".",
		IgnoredFolders:     []string{"tmp", "vendor", "node_modules"},
		IncludedExtensions: []string{".go", ".tmpl", ".html", ".json", ".yml", ".yaml"},
		BuildPath:          "/tmp",
		BuildDelay:         10 * time.Millisecond,
		BinaryName:         "weeb-dev",
		CommandFlags:       []string{"start"},
		CommandEnv:         []string{"APP_DEV=1"},
		EnableColors:       true,
	}
	if err := refresh.New(config).Start(); err != nil {
		app.Log.Error("error starting dev server", L{"err": err.Error()})
		os.Exit(1)
	}
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
