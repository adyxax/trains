package database

import (
	"reflect"
	"testing"

	"git.adyxax.org/adyxax/trains/pkg/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	// test db setup
	db, err := InitDB("sqlite3", "file::memory:?_foreign_keys=on")
	require.NoError(t, err)
	err = db.Migrate()
	require.NoError(t, err)
	userReg1 := model.UserRegistration{
		Username: "user1",
		Password: "user1_pass",
		Email:    "user1",
	}
	user1, err := db.CreateUser(&userReg1)
	require.NoError(t, err)
	user2 := *user1
	user2.Id++ // we want a token request for an invalid user id
	// Test cases
	testCases := []struct {
		name          string
		input         *model.User
		expectedError interface{}
	}{
		{"Normal user", user1, nil},
		{"A normal user can request multiple tokens", user1, nil},
		{"a non existant user id triggers an error", &user2, &QueryError{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := db.CreateSession(tc.input)
			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(tc.expectedError), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(tc.expectedError))
				require.Nil(t, valid)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, valid)
			}
		})
	}
}

func TestCreateSessionWithSQLMock(t *testing.T) {
	// Transaction begin error
	dbBeginError, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbBeginError.Close()
	// Transaction commit error
	dbCommitError, mockCommitError, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbCommitError.Close()
	mockCommitError.ExpectBegin()
	mockCommitError.ExpectExec(`INSERT INTO`).WillReturnResult(sqlmock.NewResult(1, 1))
	// Test cases
	testCases := []struct {
		name          string
		db            *DBEnv
		expectedError interface{}
	}{
		{"begin transaction error", &DBEnv{db: dbBeginError}, &TransactionError{}},
		{"commit transaction error", &DBEnv{db: dbCommitError}, &TransactionError{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := tc.db.CreateSession(&model.User{})
			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(tc.expectedError), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(tc.expectedError))
				require.Nil(t, valid)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
