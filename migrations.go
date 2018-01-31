package weeb

func addWeebMigrationsToApp(app *App) {
	app.Migrations.Add("0001_users_table", migrate0001UsersTableUp, migrate0001UsersTableDown)
}

func migrate0001UsersTableUp(app *App) error {
	return app.DB.Exec(`CREATE TABLE weeb_users (
	  id bigint PRIMARY KEY,
	  username text NOT NULL,
	  password text NOT NULL,
	  roles text[] NOT NULL,
	  created timestamp NOT NULL without time zone default (now() at time zone 'utc'),
	  updated timestamp NOT NULL without time zone default (now() at time zone 'utc'),
	)`)
}

func migrate0001UsersTableDown(app *App) error {
	return app.DB.Exec(`DROP TABLE weeb_users`)
}
