package database

import (
	"fmt"
	"reflect"
	"testing"

	"git.adyxax.org/adyxax/trains/pkg/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountStops(t *testing.T) {
	stops := []model.Stop{
		model.Stop{Id: "id1", Name: "name1"},
		model.Stop{Id: "id2", Name: "name2"},
	}
	// test db setup
	db, err := InitDB("sqlite3", "file::memory:?_foreign_keys=on")
	require.NoError(t, err)
	// check sql error
	i, err := db.CountStops()
	require.Error(t, err)
	assert.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(&QueryError{}), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(&QueryError{}))
	// normal check
	err = db.Migrate()
	require.NoError(t, err)
	err = db.ReplaceAndImportStops(stops)
	i, err = db.CountStops()
	require.NoError(t, err)
	assert.Equal(t, i, len(stops))
}

func TestGetStop(t *testing.T) {
	stops := []model.Stop{
		model.Stop{Id: "id1", Name: "name1"},
		model.Stop{Id: "id2", Name: "name2"},
	}
	// test db setup
	db, err := InitDB("sqlite3", "file::memory:?_foreign_keys=on")
	require.NoError(t, err)
	err = db.Migrate()
	require.NoError(t, err)
	err = db.ReplaceAndImportStops(stops)
	// normal check
	stop, err := db.GetStop("id1")
	require.NoError(t, err)
	assert.Equal(t, stop, &stops[0])
	// error check
	stop, err = db.GetStop("non_existent")
	require.Error(t, err)
	assert.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(&QueryError{}), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(&QueryError{}))
}

func TestGetStops(t *testing.T) {
	stops := []model.Stop{
		model.Stop{Id: "id1", Name: "name1"},
		model.Stop{Id: "id2", Name: "name2"},
	}
	// test db setup
	db, err := InitDB("sqlite3", "file::memory:?_foreign_keys=on")
	require.NoError(t, err)
	// error check
	res, err := db.GetStops()
	require.Error(t, err)
	assert.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(&QueryError{}), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(&QueryError{}))
	// normal check
	err = db.Migrate()
	require.NoError(t, err)
	err = db.ReplaceAndImportStops(stops)
	res, err = db.GetStops()
	require.NoError(t, err)
	assert.Equal(t, res, stops)
}

func TestGetStopsWithSQLMock(t *testing.T) {
	// Transaction commit error
	dbScanError, mockScanError, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mockScanError.ExpectQuery(`SELECT id, name FROM stops;`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(nil, "b").RowError(1, fmt.Errorf("row error")))
	// Test cases
	testCases := []struct {
		name          string
		db            *DBEnv
		expectedError interface{}
	}{
		{"query error when scanning row", &DBEnv{db: dbScanError}, &QueryError{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.db.GetStops()
			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(tc.expectedError), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestReplaceAndImportStops(t *testing.T) {
	// test db setup
	db, err := InitDB("sqlite3", "file::memory:?_foreign_keys=on")
	require.NoError(t, err)
	err = db.Migrate()
	require.NoError(t, err)
	// datasets
	data1 := []model.Stop{
		model.Stop{Id: "first", Name: "firstName"},
		model.Stop{Id: "second", Name: "secondName"},
	}
	data2 := []model.Stop{
		model.Stop{Id: "first", Name: "firstTest"},
		model.Stop{Id: "secondTest", Name: "secondTest"},
		model.Stop{Id: "thirdTest", Name: "thirdTest"},
	}
	testCases := []struct {
		name          string
		input         []model.Stop
		expectedError interface{}
	}{
		{"Normal insert", data1, nil},
		{"Normal insert overwrite", data2, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := db.ReplaceAndImportStops(tc.input)
			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(tc.expectedError), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestReplaceAndImportStopsWithSQLMock(t *testing.T) {
	// datasets
	data1 := []model.Stop{
		model.Stop{Id: "first", Name: "firstName"},
		model.Stop{Id: "second", Name: "secondName"},
	}
	// Transaction begin error
	dbBeginError, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbBeginError.Close()
	// Query error cannot delete from
	dbCannotDeleteFrom, mockCannotDeleteFrom, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbCannotDeleteFrom.Close()
	mockCannotDeleteFrom.ExpectBegin()
	// Transaction commit error
	dbCannotInsertError, mockCannotInsertError, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbCannotInsertError.Close()
	mockCannotInsertError.ExpectBegin()
	mockCannotInsertError.ExpectExec(`DELETE FROM`).WillReturnResult(sqlmock.NewResult(1, 1))
	// Transaction commit error
	dbCommitError, mockCommitError, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbCommitError.Close()
	mockCommitError.ExpectBegin()
	mockCommitError.ExpectExec(`DELETE FROM`).WillReturnResult(sqlmock.NewResult(1, 1))
	mockCommitError.ExpectExec(`INSERT INTO`).WillReturnResult(sqlmock.NewResult(1, 1))
	mockCommitError.ExpectExec(`INSERT INTO`).WillReturnResult(sqlmock.NewResult(1, 1))
	// Test cases
	testCases := []struct {
		name          string
		db            *DBEnv
		expectedError interface{}
	}{
		{"begin transaction error", &DBEnv{db: dbBeginError}, &TransactionError{}},
		{"query error cannot delete from", &DBEnv{db: dbCannotDeleteFrom}, &QueryError{}},
		{"query error cannot insert into", &DBEnv{db: dbCannotInsertError}, &QueryError{}},
		{"commit transaction error", &DBEnv{db: dbCommitError}, &TransactionError{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.db.ReplaceAndImportStops(data1)
			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(tc.expectedError), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
