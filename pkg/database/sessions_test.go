package database

import (
	"testing"

	"git.adyxax.org/adyxax/trains/pkg/model"
	"github.com/DATA-DOG/go-sqlmock"
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
		expectedError error
	}{
		{"Normal user", user1, nil},
		{"A normal user can request multiple tokens", user1, nil},
		{"a non existant user id triggers an error", &user2, QueryError{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := db.CreateSession(tc.input)
			if tc.expectedError != nil {
				require.Error(t, err)
				requireErrorTypeMatch(t, err, tc.expectedError)
				require.Nil(t, valid)
			} else {
				require.NoError(t, err)
				require.NotNil(t, valid)
			}
		})
	}
}

func TestCreateSessionWithSQLMock(t *testing.T) {
	// Transaction begin error
	dbBeginError, _, err := sqlmock.New()
	require.NoError(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer dbBeginError.Close()
	// Transaction commit error
	dbCommitError, mockCommitError, err := sqlmock.New()
	require.NoError(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer dbCommitError.Close()
	mockCommitError.ExpectBegin()
	mockCommitError.ExpectExec(`INSERT INTO`).WillReturnResult(sqlmock.NewResult(1, 1))
	// Test cases
	testCases := []struct {
		name          string
		db            *DBEnv
		expectedError error
	}{
		{"begin transaction error", &DBEnv{db: dbBeginError}, TransactionError{}},
		{"commit transaction error", &DBEnv{db: dbCommitError}, TransactionError{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := tc.db.CreateSession(&model.User{})
			if tc.expectedError != nil {
				require.Error(t, err)
				requireErrorTypeMatch(t, err, tc.expectedError)
				require.Nil(t, valid)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestResumeSession(t *testing.T) {
	// test db setup : of the three users only two have session tokens
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
	token1, err := db.CreateSession(user1)
	require.NoError(t, err)
	token1bis, err := db.CreateSession(user1)
	require.NoError(t, err)
	userReg2 := model.UserRegistration{
		Username: "user2",
		Password: "user2_pass",
		Email:    "user2",
	}
	user2, err := db.CreateUser(&userReg2)
	require.NoError(t, err)
	token2, err := db.CreateSession(user2)
	require.NoError(t, err)
	userReg3 := model.UserRegistration{
		Username: "user3",
		Password: "user3_pass",
		Email:    "user3",
	}
	_, err = db.CreateUser(&userReg3)
	require.NoError(t, err)
	// Test cases
	testCases := []struct {
		name          string
		input         string
		expected      *model.User
		expectedError error
	}{
		{"Normal user resume", *token1, user1, nil},
		{"Normal user resume 1bis", *token1bis, user1, nil},
		{"Normal user resume 2", *token2, user2, nil},
		{"a non existant user token triggers an error", "XXX", nil, QueryError{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := db.ResumeSession(tc.input)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Nil(t, valid)
				requireErrorTypeMatch(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, valid)
				require.Equal(t, valid.Id, tc.expected.Id)
			}
		})
	}
}
