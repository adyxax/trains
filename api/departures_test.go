package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDepartures(t *testing.T) {
	// invalid characters in token
	client := NewClient("}")
	_, err := client.GetDepartures()
	if err == nil {
		t.Fatalf("invalid characters in token should raise an error")
	}
	// unreachable server
	client = NewClient("https://")
	_, err = client.GetDepartures()
	if err == nil {
		t.Fatalf("unreachable server should raise an error")
	}
	// invalid json
	client, ts := NewTestClientFromFilename(t, "invalid.json")
	defer ts.Close()
	_, err = client.GetDepartures()
	if err == nil {
		t.Fatalf("invalid json should raise an error")
	}
	// http error
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	client = NewTestClient(ts)
	_, err = client.GetDepartures()
	if err == nil {
		t.Fatalf("404 should raise an error")
	}
	// normal working request
	client, ts = NewTestClientFromFilename(t, "normal-crepieux.json")
	defer ts.Close()
	departures, err := client.GetDepartures()
	if err != nil {
		t.Fatalf("could not get normal-crepieux departures : %s", err)
	}
	if len(departures.Departures) != 10 {
		t.Fatalf("did not decode normal-crepieux departures properly, got %d departures when expected 10", len(departures.Departures))
	}
}
