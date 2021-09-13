package database

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireErrorTypeMatch(t *testing.T, err error, expected error) {
	require.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(expected), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(expected))
}

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
