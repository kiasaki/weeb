package weeb

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// DB
type DB interface {
	Connect() error
	QueryOne(interface{}, string, ...interface{}) error
	QueryAll(interface{}, string, ...interface{}) error
	Exec(string, ...interface{}) error
	ExecNamed(string, interface{}) error
	ExecWithResult(string, ...interface{}) (sql.Result, error)
}

// PostgresDB
type PostgresDB struct {
	dbURL  string
	db     *sqlx.DB
	logger *Logger
}

func NewPostgresDB(dbURL string, logger *Logger) *PostgresDB {
	return &PostgresDB{dbURL: dbURL, logger: logger}
}

func (db *PostgresDB) Connect() error {
	if db.db != nil {
		return nil
	}

	dataSourceName, err := pq.ParseURL(db.dbURL)
	if err != nil {
		return err
	}
	db.db, err = sqlx.Connect("postgres", dataSourceName)
	return err
}

func (db *PostgresDB) QueryOne(dest interface{}, query string, args ...interface{}) error {
	if err := db.Connect(); err != nil {
		return err
	}
	db.logger.Debug("sql", L{"query": query, "args": args})
	return db.db.Get(dest, query, args...)
}

func (db *PostgresDB) QueryAll(dest interface{}, query string, args ...interface{}) error {
	if err := db.Connect(); err != nil {
		return err
	}
	db.logger.Debug("sql", L{"query": query, "args": args})
	return db.db.Select(dest, query, args...)
}

func (db *PostgresDB) QueryRow(dest []interface{}, query string, args ...interface{}) error {
	if err := db.Connect(); err != nil {
		return err
	}

	db.logger.Debug("sql", L{"query": query, "args": args})
	row := db.db.QueryRow(query, args...)
	return row.Scan(dest...)
}

func (db *PostgresDB) Exec(query string, args ...interface{}) error {
	_, err := db.ExecWithResult(query, args...)
	return err
}

func (db *PostgresDB) ExecNamed(query string, arg interface{}) error {
	if err := db.Connect(); err != nil {
		return err
	}

	db.logger.Debug("sql", L{"query": query, "arg": arg})
	_, err := db.db.NamedExec(query, arg)
	return err
}

func (db *PostgresDB) ExecWithResult(query string, args ...interface{}) (sql.Result, error) {
	if err := db.Connect(); err != nil {
		return nil, err
	}

	db.logger.Debug("sql", L{"query": query, "args": args})
	return db.db.Exec(query, args...)
}

type DBStringArray pq.StringArray
