package navitia_api_client

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"git.adyxax.org/adyxax/trains/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDepartures(t *testing.T) {
	// Simple Test cases
	testCases := []struct {
		name               string
		inputNewCLient     string
		inputGetDepartures string
		expected           []model.Departure
		expectedError      interface{}
	}{
		{"invalid characters in token should fail", "}", "test", nil, &HttpClientError{}},
		{"unreachable server should fail", "https://", "test", nil, &HttpClientError{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClient(tc.inputNewCLient)
			valid, err := client.GetDepartures(tc.inputGetDepartures)
			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(tc.expectedError), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(tc.expectedError))
				assert.Equal(t, tc.expected, valid)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, valid)
			}
		})
	}
	// Test cases with a filename
	testCasesFilename := []struct {
		name               string
		inputFilename      string
		inputGetDepartures string
		expected           []model.Departure
		expectedError      interface{}
	}{
		{"invalid json should fail", "test_data/invalid.json", "test", nil, &JsonDecodeError{}},
		{"invalid date should fail", "test_data/invalid_date.json", "test", nil, &DateParsingError{}},
	}
	for _, tc := range testCasesFilename {
		t.Run(tc.name, func(t *testing.T) {
			client, ts := newTestClientFromFilename(t, tc.inputFilename)
			defer ts.Close()
			valid, err := client.GetDepartures(tc.inputGetDepartures)
			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(tc.expectedError), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(tc.expectedError))
				assert.Equal(t, tc.expected, valid)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, valid)
			}
		})
	}
	// http error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	client := newTestClient(ts)
	_, err := client.GetDepartures("test")
	if err == nil {
		t.Fatalf("404 should raise an error")
	}
	// normal working request
	client, ts = newTestClientFromFilename(t, "test_data/normal-crepieux.json")
	defer ts.Close()
	departures, err := client.GetDepartures("test")
	if err != nil {
		t.Fatalf("could not get normal-crepieux departures : %s", err)
	}
	if len(departures) != 10 {
		t.Fatalf("did not decode normal-crepieux departures properly, got %d departures when expected 10", len(departures))
	}
	// test the cache (assuming the test takes less than 60 seconds (and it really should) it will be accurate)
	ts.Close()
	departures, err = client.GetDepartures("test")
	if err != nil {
		t.Fatalf("could not get normal-crepieux departures : %s", err)
	}
	if len(departures) != 10 {
		t.Fatalf("did not decode normal-crepieux departures properly, got %d departures when expected 10", len(departures))
	}
}
