package database

import "fmt"

// database init error
type InitError struct {
	path string
	err  error
}

func (e *InitError) Error() string {
	return fmt.Sprintf("Failed to open database : %s", e.path)
}
func (e *InitError) Unwrap() error { return e.err }

func newInitError(path string, err error) error {
	return &InitError{
		path: path,
		err:  err,
	}
}

// database migration error
type MigrationError struct {
	version int
	err     error
}

func (e *MigrationError) Error() string {
	return fmt.Sprintf("Failed to migrate database to version %d : %s", e.version, e.err)
}
func (e *MigrationError) Unwrap() error { return e.err }

func newMigrationError(version int, err error) error {
	return &MigrationError{
		version: version,
		err:     err,
	}
}

// database transaction error
type TransactionError struct {
	msg string
	err error
}

func (e *TransactionError) Error() string {
	return fmt.Sprintf("Failed to perform transaction : %s", e.msg)
}
func (e *TransactionError) Unwrap() error { return e.err }

func newTransactionError(msg string, err error) error {
	return &TransactionError{
		msg: msg,
		err: err,
	}
}
