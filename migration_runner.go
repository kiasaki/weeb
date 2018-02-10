package weeb

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Migration represents an instance of a migration the migration runner can run
type Migration struct {
	ID   string
	Up   func(app *App) error
	Down func(app *App) error
}

// MigrationRunner represents an instance of a migration runner with it's
// associated app and registered migrations
type MigrationRunner struct {
	app        *App
	migrations []*Migration
}

// NewMigrationRunner creates a MigrationRunner instance
func NewMigrationRunner(app *App) *MigrationRunner {
	return &MigrationRunner{app: app, migrations: []*Migration{}}
}

// Add adds a new migration definition to the MigrationRunner
func (m *MigrationRunner) Add(id string, upFn, downFn func(app *App) error) {
	m.migrations = append(m.migrations, &Migration{ID: id, Up: upFn, Down: downFn})
}

// EnsureTable ensures the 'migrations' table exists in the database
func (m *MigrationRunner) EnsureTable() error {
	return m.app.DB.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
		  id text,
		  created timestamp NOT NULL,
		  PRIMARY KEY (id)
		)
	`)
}

func (m *MigrationRunner) currentMigrationIndex() (int, error) {
	if err := m.EnsureTable(); err != nil {
		return 0, err
	}

	lastMigrationID := ""
	lastMigrationSQL := `SELECT id FROM migrations ORDER BY id DESC LIMIT 1`
	err := m.app.DB.QueryOne(&lastMigrationID, lastMigrationSQL)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if lastMigrationID == "" {
		return -1, nil
	}

	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].ID < m.migrations[j].ID
	})

	currentIndex := 0
	for i, m := range m.migrations {
		if m.ID > lastMigrationID {
			break
		}
		currentIndex = i
	}

	return currentIndex, nil
}

// RunUp runs the 'up' part of 'n' migrations from the last migration
func (m *MigrationRunner) RunUp(n int) error {
	currentMigrationIndex, err := m.currentMigrationIndex()
	if err != nil {
		return err
	}

	if len(m.migrations) == 0 {
		fmt.Printf("\nThere is 0 migration registered. Nothing to do\n\n")
		return nil
	}

	fmt.Println()

	targetIndex := len(m.migrations)
	if n > 0 && currentMigrationIndex+n < len(m.migrations) {
		targetIndex = currentMigrationIndex + n
	}

	var i int
	for i = currentMigrationIndex + 1; i < targetIndex; i++ {
		if err := m.migrations[i].Up(m.app); err != nil {
			fmt.Printf("Error running up for migration '%s'\n\n%v\n\n", m.migrations[i].ID, err)
			return nil
		}
		insertMigrationSQL := `INSERT INTO migrations (id, created) VALUES ($1, NOW())`
		if err := m.app.DB.Exec(insertMigrationSQL, m.migrations[i].ID); err != nil {
			return err
		}
		fmt.Printf("Ran up for '%s'\n", m.migrations[i].ID)
	}

	if i == currentMigrationIndex+1 {
		fmt.Println("Nothing to do")
	}

	fmt.Println()

	return nil
}

// RunDown runs the 'down' part of 'n' migrations from the last migration
func (m *MigrationRunner) RunDown(n int) error {
	currentMigrationIndex, err := m.currentMigrationIndex()
	if err != nil {
		return err
	}

	if len(m.migrations) == 0 {
		fmt.Printf("\nThere is 0 migration registered. Nothing to do\n\n")
		return nil
	}

	if currentMigrationIndex == -1 {
		fmt.Printf("\nNothing to do\n\n")
		return nil
	}

	if n != 1 {
		return errors.New("RunDown does not support an 'n' value other than '1'")
	}

	fmt.Println()
	migration := m.migrations[currentMigrationIndex]
	if err := migration.Down(m.app); err != nil {
		fmt.Printf("Error running down for migration '%s'\n\n%v\n\n", migration.ID, err)
		return nil
	}
	deleteMigrationSQL := `DELETE FROM migrations WHERE id = $1`
	if err := m.app.DB.Exec(deleteMigrationSQL, migration.ID); err != nil {
		return err
	}
	fmt.Printf("Ran down for '%s'\n", migration.ID)
	fmt.Println()

	return nil
}

// Tasks

func migrationRunnerTask(app *App, args []string) error {
	if len(args) == 0 {
		return migrationRunnerTaskHelp(app)
	} else if args[0] == "help" {
		return migrationRunnerTaskHelp(app)
	} else if args[0] == "list" {
		return migrationRunnerTaskList(app)
	} else if args[0] == "up" {
		return migrationRunnerTaskUp(app)
	} else if args[0] == "down" {
		return migrationRunnerTaskDown(app)
	} else if args[0] == "create" {
		return migrationRunnerTaskCreate(app, args[1:])
	}
	fmt.Printf("Error: unknown sub-task '%s' for task 'migrate'\n\n", args[0])
	return nil
}

func migrationRunnerTaskHelp(app *App) error {
	fmt.Println("'migrate' task usage:")
	fmt.Println()
	fmt.Println("    list    shows all registered migrations")
	fmt.Println("    up      runs the 'up' part for all pending migrations")
	fmt.Println("    down    runs the 'down' part of the latest migration")
	fmt.Println("    create  creates a new migration file in 'migrations/'")
	fmt.Println()
	return nil
}

func migrationRunnerTaskList(app *App) error {
	migrations := app.Migrations.migrations
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].ID < migrations[j].ID
	})
	fmt.Println()
	for _, migration := range migrations {
		fmt.Println("    " + migration.ID)
	}
	fmt.Println()
	return nil
}

func migrationRunnerTaskCreate(app *App, args []string) error {
	if err := os.MkdirAll("migrations", os.ModePerm); err != nil {
		return nil
	}

	migrationsPackageFile, err := os.OpenFile(filepath.Join("migrations", "migrations.go"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer migrationsPackageFile.Close()
	migrationsPackageFile.Write([]byte(`package migrations

// This is an AUTOGENERATED file, if you edit it, you will loose your changes

import (
	"github.com/kiasaki/weeb"
)

var migrations = []*weeb.Migration{}

func migration(id string, upFn, downFn func(*weeb.App)error) struct{} {
	migrations = append(migrations, &weeb.Migration{ID:id, Up: upFn, Down: downFn})
	return struct{}{}
}

// AddMigrationsToApp adds migrations defined in this package to the given 'App'
func AddMigrationsToApp(app *weeb.App) {
	for _, m := range migrations {
		app.Migrations.Add(m.ID, m.Up, m.Down)
	}
}
`))

	migrationID := time.Now().UTC().Format("20060102150405")
	migrationName := migrationID
	if len(args) >= 1 {
		migrationName += "_" + ToSnakeCase(args[0])
	}
	migrationFileName := migrationName + ".go"
	migrationFile, err := os.OpenFile(filepath.Join("migrations", migrationFileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer migrationFile.Close()
	migrationFile.Write([]byte(fmt.Sprintf(`package migrations

import (
	"github.com/kiasaki/weeb"
)

var _ = migration("%v", up%v, down%v)

func up%v(app *weeb.App) error {
	return app.DB.Exec(`+"`"+`
	CREATE TABLE ...
	`+"`"+`)
}

func down%v(app *weeb.App) error {
	return app.DB.Exec(`+"`"+`
	DROP TABLE ...
	`+"`"+`)
}
`, migrationName, migrationID, migrationID, migrationID, migrationID)))

	fmt.Printf("\nCreated migration file 'migrations/%s'\n\n", migrationFileName)

	return nil
}

func migrationRunnerTaskUp(app *App) error {
	return app.Migrations.RunUp(-1)
}

func migrationRunnerTaskDown(app *App) error {
	return app.Migrations.RunDown(1)
}
