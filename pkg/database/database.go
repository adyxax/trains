package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DBEnv is the struct that holds this package together
type DBEnv struct {
	db *sql.DB
}

// InitDB initializes database access and the connection pool
func InitDB(dbType string, dsn string) (*DBEnv, error) {
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		return nil, newInitError(dsn, err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close() // TODO think about handling this error case?
		return nil, newInitError(dsn, err)
	}

	return &DBEnv{db: db}, nil
}

// Migrate performs the migrations of the database to the latest schema_version
func (env *DBEnv) Migrate() error {
	var currentVersion int
	if err := env.db.QueryRow(`SELECT version FROM schema_version`).Scan(&currentVersion); err != nil {
		currentVersion = 0
	}
	for version := currentVersion; version < len(migrations); version++ {
		newVersion := version + 1
		tx, err := env.db.Begin()
		if err != nil {
			return newTransactionError("Could not begin transaction", err)
		}
		if err := migrations[version](tx); err != nil {
			tx.Rollback()
			return newMigrationError(newVersion, err)
		}
		if _, err := tx.Exec(`DELETE FROM schema_version; INSERT INTO schema_version (version) VALUES ($1)`, newVersion); err != nil {
			tx.Rollback()
			return newMigrationError(newVersion, err)
		}
		if err := tx.Commit(); err != nil {
			return newTransactionError("Could not commit transaction", err)
		}
	}
	return nil
}
