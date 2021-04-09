package database

import "database/sql"

// allMigrations is the list of migrations to perform to get an up to date database
// Order is important. Add new migrations at the end of the list.
var allMigrations = []func(tx *sql.Tx) error{
	func(tx *sql.Tx) (err error) {
		sql := `
			CREATE TABLE schema_version (
				version INTEGER NOT NULL
			);`
		_, err = tx.Exec(sql)
		return err
	},
}

// This variable exists so that tests can override it
var migrations = allMigrations
