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
	runHttpTest(t, e, loginHandler, &httpTestCase{
		name: "a login attempt without username should error",
		input: httpTestInput{
			method: http.MethodPost,
			path:   "/login",
			data: url.Values{
				"password": []string{"password1"},
			},
		},
		expect: httpTestExpect{
			code: http.StatusBadRequest,
			err:  &statusError{http.StatusNotFound, simpleErrorMessage},
		},
	})
	runHttpTest(t, e, loginHandler, &httpTestCase{
		name: "a login attempt with multiple usernames should error",
		input: httpTestInput{
			method: http.MethodPost,
			path:   "/login",
			data: url.Values{
				"username": []string{"user1", "user2"},
				"password": []string{"password1"},
			},
		},
		expect: httpTestExpect{
			code: http.StatusBadRequest,
			err:  &statusError{http.StatusNotFound, simpleErrorMessage},
		},
	})
	runHttpTest(t, e, loginHandler, &httpTestCase{
		name: "a login attempt with an invalid username should error",
		input: httpTestInput{
			method: http.MethodPost,
			path:   "/login",
			data: url.Values{
				"username": []string{"%"},
				"password": []string{"password1"},
			},
		},
		expect: httpTestExpect{
			code: http.StatusBadRequest,
			err:  &statusError{http.StatusNotFound, simpleErrorMessage},
		},
	})
	runHttpTest(t, e, loginHandler, &httpTestCase{
		name: "a login attempt without password should error",
		input: httpTestInput{
			method: http.MethodPost,
			path:   "/login",
			data: url.Values{
				"username": []string{"user1"},
			},
		},
		expect: httpTestExpect{
			code: http.StatusBadRequest,
			err:  &statusError{http.StatusNotFound, simpleErrorMessage},
		},
	})
	runHttpTest(t, e, loginHandler, &httpTestCase{
		name: "a login attempt with multiple passwords should error",
		input: httpTestInput{
			method: http.MethodPost,
			path:   "/login",
			data: url.Values{
				"username": []string{"user1"},
				"password": []string{"password1", "password2"},
			},
		},
		expect: httpTestExpect{
			code: http.StatusBadRequest,
			err:  &statusError{http.StatusNotFound, simpleErrorMessage},
		},
	})
	runHttpTest(t, e, loginHandler, &httpTestCase{
		name: "a login attempt with an empty password should error",
		input: httpTestInput{
			method: http.MethodPost,
			path:   "/login",
			data: url.Values{
				"username": []string{"user1"},
				"password": []string{""},
			},
		},
		expect: httpTestExpect{
			code: http.StatusBadRequest,
			err:  &statusError{http.StatusNotFound, simpleErrorMessage},
		},
	})
	// Test other request types
	methods := []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPut,
		http.MethodTrace,
	}
	for _, method := range methods {
		runHttpTest(t, e, loginHandler, &httpTestCase{
			name: "a login attempt with an invalid method should error",
			input: httpTestInput{
				method: method,
				path:   "/login",
			},
			expect: httpTestExpect{
				code: http.StatusMethodNotAllowed,
				err:  &statusError{http.StatusMethodNotAllowed, simpleErrorMessage},
			},
		})
	}
}
