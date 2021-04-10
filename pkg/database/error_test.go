package database

import "testing"

func TestErrorsCoverage(t *testing.T) {
	initErr := InitError{}
	_ = initErr.Error()
	_ = initErr.Unwrap()
	migrationErr := MigrationError{}
	_ = migrationErr.Error()
	_ = migrationErr.Unwrap()
	passwordError := PasswordError{}
	_ = passwordError.Error()
	_ = passwordError.Unwrap()
	queryErr := QueryError{}
	_ = queryErr.Error()
	_ = queryErr.Unwrap()
	transactionErr := TransactionError{}
	_ = transactionErr.Error()
	_ = transactionErr.Unwrap()
}
