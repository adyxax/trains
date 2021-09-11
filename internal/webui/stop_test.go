package webui

import (
	"net/http"
	"testing"

	"git.adyxax.org/adyxax/trains/pkg/config"
	"git.adyxax.org/adyxax/trains/pkg/database"
	"git.adyxax.org/adyxax/trains/pkg/model"
	"github.com/stretchr/testify/require"
)

func TestStopHandler(t *testing.T) {
	// test environment setup
	dbEnv, err := database.InitDB("sqlite3", "file::memory:?_foreign_keys=on")
	require.Nil(t, err)
	err = dbEnv.Migrate()
	require.Nil(t, err)
	user1, err := dbEnv.CreateUser(&model.UserRegistration{Username: "user1", Password: "password1", Email: "julien@adyxax.org"})
	require.Nil(t, err)
	_, err = dbEnv.Login(&model.UserLogin{Username: "user1", Password: "password1"})
	require.Nil(t, err)
	token1, err := dbEnv.CreateSession(user1)
	require.Nil(t, err)
	err = dbEnv.ReplaceAndImportStops([]model.Stop{model.Stop{Id: "stop_area:test:01", Name: "test"}})
	require.Nil(t, err)
	e := env{
		dbEnv: dbEnv,
		conf:  &config.Config{},
	}
	departures1 := []model.Departure{
		model.Departure{
			Direction: "test direction",
			Arrival:   "20210503T150405",
		},
	}
	e.navitia = &NavitiaMockClient{departures: departures1, err: nil}
	// test GET requests
	runHttpTest(t, &e, stopHandler, &httpTestCase{
		name: "a simple get when not logged in should redirect to the login page",
		input: httpTestInput{
			method: http.MethodGet,
			path:   "/stop",
		},
		expect: httpTestExpect{
			code:     http.StatusFound,
			location: "/login",
		},
	})
	runHttpTest(t, &e, stopHandler, &httpTestCase{
		name: "a simple get when logged in should display the stops list",
		input: httpTestInput{
			method: http.MethodGet,
			path:   "/stop",
			cookie: &http.Cookie{Name: sessionCookieName, Value: *token1},
		},
		expect: httpTestExpect{
			code:       http.StatusOK,
			bodyString: "stop_area:test:01",
		},
	})
}
