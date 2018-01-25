package migrations

import (
	"github.com/kiasaki/weeb"
)

var _ = migration("20180125020047", up20180125020047, down20180125020047)

func up20180125020047(app *weeb.App) error {
	return app.DB.Exec(`
	SELECT 1
	`)
}

func down20180125020047(app *weeb.App) error {
	return app.DB.Exec(`
	SELECT 2
	`)
}
