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

// Update updates given entity in the database
func (h *DBHelper) Update(e Entity) error {
	updateSQL := "UPDATE %s SET %s WHERE %s = :%s"
	idField := e.Fields()[0]

	fields := []string{}
	for _, field := range e.Fields() {
		fields = append(fields, fmt.Sprintf("%s = :%s", field, field))
	}
	fieldsSQL := strings.Join(fields, ", ")

	updateSQL = fmt.Sprintf(updateSQL, e.Table(), fieldsSQL, idField, idField)
	return h.db.ExecNamed(updateSQL, e)
}

// FindParams specifies what exatly it is we want to find
type FindParams struct {
	Limit   int64
	OrderBy []string
	Where   map[string]interface{}
}

func (h *DBHelper) findSQLFor(e Entity, params FindParams) (string, []interface{}) {
	sql := "SELECT %s FROM %s"
	values := []interface{}{}
	fields := strings.Join(e.Fields(), ", ")
	sql = fmt.Sprintf(sql, fields, e.Table())

	i := 1

	if len(params.Where) > 0 {
		whereSqls := []string{}
		for whereField, whereValue := range params.Where {
			whereSqls = append(whereSqls, fmt.Sprintf("%s = $%d", ToSnakeCase(whereField), i))
			values = append(values, whereValue)
			i++
		}
		if len(whereSqls) > 0 {
			sql += " WHERE " + strings.Join(whereSqls, " AND ")
		}
	}
	if len(params.OrderBy) > 0 {
		for _, order := range params.OrderBy {
			sql += fmt.Sprintf(" ORDER BY $%d", i)
			i++
			if order[0] == '-' {
				order = order[1:]
				sql += " DESC"
			} else {
				sql += " ASC"
			}
			values = append(values, ToSnakeCase(order))
		}
	}
	if params.Limit > 0 {
		sql += fmt.Sprintf(" LIMIT $%d", i)
		i++
		values = append(values, params.Limit)
	}

	return sql, values
}

// Find finds one entity in the database based on the provided filters,
// limits and sort orders
func (h *DBHelper) Find(e Entity, params FindParams) error {
	sql, values := h.findSQLFor(e, params)
	return h.db.QueryOne(e, sql, values...)
}

// FindAll finds all entities in the database that match the provided filters,
// limits and sort orders
func (h *DBHelper) FindAll(e Entity, result interface{}, params FindParams) error {
	sql, values := h.findSQLFor(e, params)
	return h.db.QueryAll(result, sql, values...)
}
