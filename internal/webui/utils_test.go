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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	err        interface{}
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
			assert.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(tc.expect.err), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(tc.expect.err))
		} else {
			require.NoError(t, err)
			assert.Equal(t, tc.expect.code, rr.Code)
			if tc.expect.bodyString != "" {
				assert.Contains(t, rr.Body.String(), tc.expect.bodyString)
			}
			if tc.expect.location != "" {
				assert.Contains(t, rr.HeaderMap, "Location")
				assert.Len(t, rr.HeaderMap["Location"], 1)
				assert.Equal(t, rr.HeaderMap["Location"][0], tc.expect.location)
			}
			if tc.expect.setsCookie {
				assert.Contains(t, rr.HeaderMap, "Set-Cookie")
			} else {
				assert.NotContains(t, rr.HeaderMap, "Set-Cookie")
			}
		}
	})
}
