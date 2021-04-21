package webui

import (
	"net/http"
	"net/url"
	"testing"

	"git.adyxax.org/adyxax/trains/pkg/database"
	"git.adyxax.org/adyxax/trains/pkg/model"
	"github.com/stretchr/testify/require"
)

func TestLoginHandler(t *testing.T) {
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
	e := &env{dbEnv: dbEnv}
	// test GET requests
	runHttpTest(t, e, loginHandler, &httpTestCase{
		name: "a simple get should display the login page",
		input: httpTestInput{
			method: http.MethodGet,
			path:   "/login",
		},
		expect: httpTestExpect{
			code:       http.StatusOK,
			bodyString: "<form action=\"/login\"",
		},
	})
	runHttpTest(t, e, loginHandler, &httpTestCase{
		name: "an invalid or expired token should also just display the login page",
		input: httpTestInput{
			method: http.MethodGet,
			path:   "/login",
			cookie: &http.Cookie{Name: sessionCookieName, Value: "graou"},
		},
		expect: httpTestExpect{
			code:       http.StatusOK,
			bodyString: "<form action=\"/login\"",
		},
	})
	runHttpTest(t, e, loginHandler, &httpTestCase{
		name: "if already logged in we should be redirected to /",
		input: httpTestInput{
			method: http.MethodGet,
			path:   "/login",
			cookie: &http.Cookie{Name: sessionCookieName, Value: *token1},
		},
		expect: httpTestExpect{
			code:     http.StatusFound,
			location: "/",
		},
	})
	runHttpTest(t, e, loginHandler, &httpTestCase{
		name: "an invalid path should get a 404",
		input: httpTestInput{
			method: http.MethodGet,
			path:   "/login/non_existent",
		},
		expect: httpTestExpect{
			code: http.StatusNotFound,
			err:  &statusError{http.StatusNotFound, simpleErrorMessage},
		},
	})
	runHttpTest(t, e, loginHandler, &httpTestCase{
		name: "an invalid path should get a 404 even if we are already logged in",
		input: httpTestInput{
			method: http.MethodGet,
			path:   "/login/non_existent",
			cookie: &http.Cookie{Name: sessionCookieName, Value: *token1},
		},
		expect: httpTestExpect{
			code: http.StatusNotFound,
			err:  &statusError{http.StatusNotFound, simpleErrorMessage},
		},
	})
	// Test POST requests
	runHttpTest(t, e, loginHandler, &httpTestCase{
		name: "a valid login attempt should succeed and redirect to /",
		input: httpTestInput{
			method: http.MethodPost,
			path:   "/login",
			data: url.Values{
				"username": []string{"user1"},
				"password": []string{"password1"},
			},
		},
		expect: httpTestExpect{
			code:       http.StatusFound,
			location:   "/",
			setsCookie: true,
		},
	})
	//errorNoUsername := newTestRequest(t, http.MethodPost, "/login", nil)
	//// too many username fields
	//dataWtfUsername := url.Values{"username": []string{"user1", "user2"}}
	//errorWtfUsername, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(dataWtfUsername.Encode()))
	//require.Nil(t, err)
	//errorWtfUsername.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//// Invalid username
	//dataInvalidUsername := url.Values{"username": []string{"%"}}
	//errorInvalidUsername, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(dataInvalidUsername.Encode()))
	//require.Nil(t, err)
	//errorInvalidUsername.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//// no password field
	//dataNoPassword := url.Values{"username": []string{"user1"}}
	//errorNoPassword, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(dataNoPassword.Encode()))
	//require.Nil(t, err)
	//errorNoPassword.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//// too many password fields
	//dataWtfPassword := url.Values{"username": []string{"user1"}, "password": []string{"user1", "user2"}}
	//errorWtfPassword, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(dataWtfPassword.Encode()))
	//require.Nil(t, err)
	//errorWtfPassword.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//// Invalid password
	//dataInvalidPassword := url.Values{"username": []string{"user1"}, "password": []string{""}}
	//errorInvalidPassword, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(dataInvalidPassword.Encode()))
	//require.Nil(t, err)
	//errorInvalidPassword.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//// run the tests
	////    {"error no username", &env{dbEnv: dbEnv}, errorNoUsername, &expected{err: &statusError{code: 500, err: simpleError}}},
	////    {"error wtf username", &env{dbEnv: dbEnv}, errorWtfUsername, &expected{err: &statusError{code: 500, err: simpleError}}},
	////    {"error invalid username", &env{dbEnv: dbEnv}, errorInvalidUsername, &expected{err: &statusError{code: 500, err: simpleError}}},
	////    {"error no password", &env{dbEnv: dbEnv}, errorNoPassword, &expected{err: &statusError{code: 500, err: simpleError}}},
	////    {"error wtf password", &env{dbEnv: dbEnv}, errorWtfPassword, &expected{err: &statusError{code: 500, err: simpleError}}},
	////    {"error invalid password", &env{dbEnv: dbEnv}, errorInvalidPassword, &expected{err: &statusError{code: 500, err: simpleError}}},
	////}
}
