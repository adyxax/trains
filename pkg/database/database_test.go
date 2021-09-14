package database

import (
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestInitDB(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		dbType        string
		dsn           string
		expectedError error
	}{
		{"Invalid dbType", "non-existant", "test", InitError{}},
		{"Non existant path", "sqlite3", "/non-existant/non-existant", InitError{}},
		{"Working DB", "sqlite3", "file::memory:?_foreign_keys=on", nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := InitDB(tc.dbType, tc.dsn)
			if tc.expectedError != nil {
				require.Error(t, err)
				requireErrorTypeMatch(t, err, tc.expectedError)
				require.Nil(t, db)
			} else {
				require.NoError(t, err)
				require.NotNil(t, db)
			}
		})
	}
}

func TestMigrate(t *testing.T) {
	badMigration := []func(tx *sql.Tx) error{
		func(tx *sql.Tx) (err error) {
			return newMigrationError(0, nil)
		},
	}
	noSchemaVersionMigration := []func(tx *sql.Tx) error{
		func(tx *sql.Tx) (err error) {
			return nil
		},
	}
	onlyFirstMigration := []func(tx *sql.Tx) error{migrations[0]}
	notFromScratchMigration := []func(tx *sql.Tx) error{
		func(tx *sql.Tx) (err error) {
			return MigrationError{version: 1, err: nil}
		},
		func(tx *sql.Tx) (err error) {
			return nil
		},
	}
	os.Remove("testfile_notFromScratch.db")
	notFromScratchDB, err := InitDB("sqlite3", "file:testfile_notFromScratch.db?_foreign_keys=on")
	require.NoError(t, err, "Failed to init testfile.db : %+v", err)
	defer os.Remove("testfile_notFromScratch.db")
	migrations = onlyFirstMigration
	err = notFromScratchDB.Migrate()
	require.NoError(t, err, "Failed to migrate testfile.db to first schema version : %+v", err)

	// Test cases
	testCases := []struct {
		name          string
		dsn           string
		migrs         []func(tx *sql.Tx) error
		expectedError error
	}{
		{"bad migration", "file::memory:?_foreign_keys=on", badMigration, MigrationError{}},
		{"no schema_version migration", "file::memory:?_foreign_keys=on", noSchemaVersionMigration, MigrationError{}},
		{"not from scratch", "file:testfile_notFromScratch.db?_foreign_keys=on", notFromScratchMigration, nil},
		{"from scratch", "file::memory:?_foreign_keys=on", allMigrations, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := InitDB("sqlite3", tc.dsn)
			require.NoError(t, err)
			migrations = tc.migrs
			err = db.Migrate()
			if tc.expectedError != nil {
				require.Error(t, err)
				requireErrorTypeMatch(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMigrateWithSQLMock(t *testing.T) {
	fakeMigration := []func(tx *sql.Tx) error{
		func(tx *sql.Tx) (err error) {
			return nil
		},
	}
	// Transaction begin error
	dbBeginError, _, err := sqlmock.New()
	require.NoError(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer dbBeginError.Close()
	// Transaction commit error
	dbCommitError, mockCommitError, err := sqlmock.New()
	require.NoError(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer dbCommitError.Close()
	mockCommitError.ExpectBegin()
	mockCommitError.ExpectExec(`DELETE FROM schema_version`).WillReturnResult(sqlmock.NewResult(1, 1))
	// Test cases
	testCases := []struct {
		name          string
		db            *DBEnv
		migrs         []func(tx *sql.Tx) error
		expectedError error
	}{
		{"begin transaction error", &DBEnv{db: dbBeginError}, fakeMigration, TransactionError{}},
		{"commit transaction error", &DBEnv{db: dbCommitError}, fakeMigration, TransactionError{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			migrations = tc.migrs
			err = tc.db.Migrate()
			migrations = allMigrations
			if tc.expectedError != nil {
				require.Error(t, err)
				requireErrorTypeMatch(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
