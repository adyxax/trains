package database

import "testing"

func TestErrorsCoverage(t *testing.T) {
	initErr := InitError{}
	_ = initErr.Error()
	_ = initErr.Unwrap()
	migrationErr := MigrationError{}
	_ = migrationErr.Error()
	_ = migrationErr.Unwrap()
	transactionErr := TransactionError{}
	_ = transactionErr.Error()
	_ = transactionErr.Unwrap()
}
