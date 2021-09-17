package webui

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"git.adyxax.org/adyxax/trains/pkg/model"
	"github.com/stretchr/testify/require"
)

func requireErrorTypeMatch(t *testing.T, err error, expected error) {
	require.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(expected), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(expected))
}

type NavitiaMockClient struct {
	departures []model.Departure
	stops      []model.Stop
	err        error
}

func (c *NavitiaMockClient) GetDepartures(stop string) (departures []model.Departure, err error) {
	return c.departures, c.err
}

func (c *NavitiaMockClient) GetStops() (stops []model.Stop, err error) {
	return c.stops, c.err
}

var simpleErrorMessage = fmt.Errorf("")

type httpTestCase struct {
	name   string
	input  httpTestInput
	expect httpTestExpect
}
type httpTestInput struct {
	method string
	path   string
	cookie *http.Cookie
	data   url.Values
}
type httpTestExpect struct {
	code       int
	bodyString string
	location   string
	setsCookie bool
	err        error
}

func runHttpTest(t *testing.T, e *env, h func(e *env, w http.ResponseWriter, r *http.Request) error, tc *httpTestCase) {
	req, err := http.NewRequest(tc.input.method, tc.input.path, nil)
	require.Nil(t, err)
	if tc.input.data != nil {
		req, err = http.NewRequest(tc.input.method, tc.input.path, strings.NewReader(tc.input.data.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	if tc.input.cookie != nil {
		req.AddCookie(tc.input.cookie)
	}
	t.Run(tc.name, func(t *testing.T) {
		rr := httptest.NewRecorder()
		err := h(e, rr, req)
		if tc.expect.err != nil {
			require.Error(t, err)
			requireErrorTypeMatch(t, err, tc.expect.err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tc.expect.code, rr.Code)
			if tc.expect.bodyString != "" {
				require.Contains(t, rr.Body.String(), tc.expect.bodyString)
			}
			if tc.expect.location != "" {
				require.Contains(t, rr.HeaderMap, "Location")
				require.Len(t, rr.HeaderMap["Location"], 1)
				require.Equal(t, rr.HeaderMap["Location"][0], tc.expect.location)
			}
			if tc.expect.setsCookie {
				require.Contains(t, rr.HeaderMap, "Set-Cookie")
			} else {
				require.NotContains(t, rr.HeaderMap, "Set-Cookie")
			}
		}
	})
}
