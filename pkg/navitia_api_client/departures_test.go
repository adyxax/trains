package navitia_api_client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDepartures(t *testing.T) {
	// invalid characters in token
	client := NewClient("}")
	_, err := client.GetDepartures("test")
	if err == nil {
		t.Fatalf("invalid characters in token should raise an error because the url is invalid")
	}
	// unreachable server
	client = NewClient("https://")
	_, err = client.GetDepartures("test")
	if err == nil {
		t.Fatalf("unreachable server should raise an error")
	}
	// invalid json
	client, ts := newTestClientFromFilename(t, "test_data/invalid.json")
	defer ts.Close()
	_, err = client.GetDepartures("test")
	if err == nil {
		t.Fatalf("invalid json should raise an error")
	}
	// http error
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	client = newTestClient(ts)
	_, err = client.GetDepartures("test")
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
	if len(departures.Departures) != 10 {
		t.Fatalf("did not decode normal-crepieux departures properly, got %d departures when expected 10", len(departures.Departures))
	}
	// test the cache (assuming the test takes less than 60 seconds (and it really should) it will be accurate)
	ts.Close()
	departures, err = client.GetDepartures("test")
	if err != nil {
		t.Fatalf("could not get normal-crepieux departures : %s", err)
	}
	if len(departures.Departures) != 10 {
		t.Fatalf("did not decode normal-crepieux departures properly, got %d departures when expected 10", len(departures.Departures))
	}
}
