package database

import (
	"reflect"
	"testing"

	"git.adyxax.org/adyxax/trains/pkg/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	// test db setup
	db, err := InitDB("sqlite3", "file::memory:?_foreign_keys=on")
	require.NoError(t, err)
	err = db.Migrate()
	require.NoError(t, err)
	// a normal user
	normalUser := model.UserRegistration{
		Username: "testUsername",
		Password: "testPassword",
		Email:    "testEmail",
	}
	// Test another normal user
	normalUser2 := model.UserRegistration{
		Username: "testUsername2",
		Password: "yei*j2Ien2xa?g6bieh~i=asoo5ii7Bi",
		Email:    "testEmail",
	}
	// Test cases
	testCases := []struct {
		name          string
		db            *DBEnv
		input         *model.UserRegistration
		expected      int
		expectedError interface{}
	}{
		{"Normal user", db, &normalUser, 1, nil},
		{"Duplicate user", db, &normalUser, 0, &QueryError{}},
		{"Normal user 2", db, &normalUser2, 2, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := tc.db.CreateUser(tc.input)
			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(tc.expectedError), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(tc.expectedError))
				require.Nil(t, valid)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, valid.Id)
			}
		})
	}
	// Test for bad password
	passwordFunction = func(password []byte, cost int) ([]byte, error) { return nil, newPasswordError(nil) }
	valid, err := db.CreateUser(&normalUser)
	passwordFunction = bcrypt.GenerateFromPassword
	require.Error(t, err)
	require.Nil(t, valid)
	assert.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(&PasswordError{}), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(&PasswordError{}))
}

func TestCreateUserWithSQLMock(t *testing.T) {
	// Transaction begin error
	dbBeginError, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbBeginError.Close()
	// Transaction LastInsertId not supported
	dbLastInsertIdError, mockLastInsertIdError, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbLastInsertIdError.Close()
	mockLastInsertIdError.ExpectBegin()
	mockLastInsertIdError.ExpectExec(`INSERT INTO`).WillReturnResult(sqlmock.NewErrorResult(&TransactionError{"test", nil}))
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
		{"last insert id transaction error", &DBEnv{db: dbLastInsertIdError}, &TransactionError{}},
		{"commit transaction error", &DBEnv{db: dbCommitError}, &TransactionError{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := tc.db.CreateUser(&model.UserRegistration{})
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
