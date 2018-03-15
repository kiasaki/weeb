package weeb

import (
	"fmt"
	"strings"
)

// Entity represents a entity that the DBHelper can use. It expose needed
// information about how the entity maps to the database like the table name
// it lives in and fields it has persisted
type Entity interface {
	Table() string
	Fields() []string
}

// DBHelper represents a database helper providing utility functions around
// a database's basic Query and Exec methods
type DBHelper struct {
	db DB
}

// NewDBHelper create a new instance of a database helper associated to the
// given database
func NewDBHelper(db DB) *DBHelper {
	return &DBHelper{db: db}
}

// Insert inserts given entity in the database
func (h *DBHelper) Insert(e Entity) error {
	insertSQL := "INSERT INTO %s (%s) VALUES (%s)"
	fields := strings.Join(e.Fields(), ", ")
	placeholders := ":" + strings.Join(e.Fields(), ", :")
	insertSQL = fmt.Sprintf(insertSQL, e.Table(), fields, placeholders)
	return h.db.ExecNamed(insertSQL, e)
}
