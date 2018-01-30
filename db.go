package weeb

import (
	"database/sql"

	// Import postgres db drivers here
	_ "github.com/lib/pq"
)

// DB
type DB interface {
	Connect() error
	QueryOne(interface{}, string, ...interface{}) error
	Query([]interface{}, string, ...interface{}) error
	QueryRow([]interface{}, string, ...interface{}) error
	QueryRows(string, ...interface{}) (*sql.Rows, error)
	Exec(string, ...interface{}) error
	ExecWithResult(string, ...interface{}) (sql.Result, error)
}

// PostgresDB
type PostgresDB struct {
	dbURL  string
	db     *sql.DB
	logger *Logger
}

func NewPostgresDB(dbURL string, logger *Logger) *PostgresDB {
	return &PostgresDB{dbURL: dbURL, logger: logger}
}

func (db *PostgresDB) Connect() error {
	if db.db != nil {
		return nil
	}

	database, err := sql.Open("postgres", db.dbURL)
	if err != nil {
		return err
	}
	db.db = database
	err = db.db.Ping()
	if err != nil {
		db.db.Close()
		return err
	}
	return nil
}

func (db *PostgresDB) QueryOne(dest interface{}, query string, args ...interface{}) error {
	destList := []interface{}{dest}
	return db.Query(destList, query, args...)
}

func (db *PostgresDB) Query(dest []interface{}, query string, args ...interface{}) error {
	if err := db.Connect(); err != nil {
		return err
	}

	db.logger.Debug("sql", L{"query": query, "args": args})
	rows, err := db.db.Query(query, args...)
	if err != nil {
		return nil
	}

	index := 0
	for rows.Next() {
		scan(dest[index], rows)
		index++
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}

func (db *PostgresDB) QueryRow(dest []interface{}, query string, args ...interface{}) error {
	if err := db.Connect(); err != nil {
		return err
	}

	db.logger.Debug("sql", L{"query": query, "args": args})
	row := db.db.QueryRow(query, args...)
	return row.Scan(dest...)
}

func (db *PostgresDB) QueryRows(query string, args ...interface{}) (*sql.Rows, error) {
	if err := db.Connect(); err != nil {
		return nil, err
	}

	db.logger.Debug("sql", L{"query": query, "args": args})
	return db.db.Query(query, args...)
}

func (db *PostgresDB) Exec(query string, args ...interface{}) error {
	_, err := db.ExecWithResult(query, args...)
	return err
}

func (db *PostgresDB) ExecWithResult(query string, args ...interface{}) (sql.Result, error) {
	if err := db.Connect(); err != nil {
		return nil, err
	}

	db.logger.Debug("sql", L{"query": query, "args": args})
	return db.db.Exec(query, args...)
}
