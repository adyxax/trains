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

func TestGetStops(t *testing.T) {
	// Simple Test cases
	testCases := []struct {
		name           string
		inputNewCLient string
		expected       []model.Stop
		expectedError  interface{}
	}{
		{"invalid characters in token should fail", "}", nil, &HttpClientError{}},
		{"unreachable server should fail", "https://", nil, &HttpClientError{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClient(tc.inputNewCLient)
			valid, err := client.GetStops()
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
		name          string
		inputFilename string
		expected      []model.Stop
		expectedError interface{}
	}{
		{"invalid json should fail", "test_data/invalid.json", nil, &JsonDecodeError{}},
	}
	for _, tc := range testCasesFilename {
		t.Run(tc.name, func(t *testing.T) {
			client, ts := newTestClientFromFilename(t, tc.inputFilename)
			defer ts.Close()
			valid, err := client.GetStops()
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
	_, err := client.GetStops()
	if err == nil {
		t.Fatalf("404 should raise an error")
	}
	// normal working request
	client, ts = newTestClientFromFilename(t, "test_data/4-train-stops.json")
	defer ts.Close()
	stops, err := client.GetStops()
	if err != nil {
		t.Fatalf("could not get train stops : %s", err)
	}
	// 4 records but one is empty (navitia api quirk)
	if len(stops) != 3 {
		t.Fatalf("did not decode train stops properly, got %d train stops when expected 4", len(stops))
	}
	// normal request in multiple pages
	client, ts = newTestClientFromFilenames(t, []testClientCase{
		testClientCase{"/coverage/sncf/stop_areas?count=1000&start_page=0", "test_data/4-train-stops-page-0.json"},
		testClientCase{"/coverage/sncf/stop_areas?count=1000&start_page=1", "test_data/4-train-stops-page-1.json"},
		testClientCase{"/coverage/sncf/stop_areas?count=1000&start_page=2", "test_data/4-train-stops-page-2.json"},
	})
	defer ts.Close()
	stops, err = client.GetStops()
	if err != nil {
		t.Fatalf("could not get train stops : %+v", err)
	}
	// 12 records but one is empty (navitia api quirk)
	if len(stops) != 11 {
		t.Fatalf("did not decode train stops properly, got %d train stops when expected 4", len(stops))
	}
	// failing request in multiple pages with last one missing
	client, ts = newTestClientFromFilenames(t, []testClientCase{
		testClientCase{"/coverage/sncf/stop_areas?count=1000&start_page=0", "test_data/4-train-stops-page-0.json"},
		testClientCase{"/coverage/sncf/stop_areas?count=1000&start_page=1", "test_data/4-train-stops-page-1.json"},
	})
	defer ts.Close()
	stops, err = client.GetStops()
	if err == nil {
		t.Fatalf("should not be able to get train stops : %+v", stops)
	}
}
