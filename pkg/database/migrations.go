package database

import "database/sql"

// allMigrations is the list of migrations to perform to get an up to date database
// Order is important. Add new migrations at the end of the list.
var allMigrations = []func(tx *sql.Tx) error{
	func(tx *sql.Tx) (err error) {
		sql := `
			CREATE TABLE schema_version (
				version INTEGER NOT NULL
			);
			CREATE TABLE users (
				id INTEGER PRIMARY KEY,
				username TEXT NOT NULL UNIQUE,
				hash TEXT,
				email TEXT,
				created_at DATE DEFAULT (datetime('now')),
				last_login_at DATE DEFAULT NULL
			);
			CREATE TABLE sessions (
				token TEXT NOT NULL UNIQUE,
				user_id INTEGER NOT NULL,
				created_at DATE DEFAULT (datetime('now')),
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
			);`
		_, err = tx.Exec(sql)
		return err
	},
}

// This variable exists so that tests can override it
var migrations = allMigrations
