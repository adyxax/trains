package navitia_api_client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// package utilities
func NewTestClient(ts *httptest.Server) *Client {
	return &Client{
		baseURL:    fmt.Sprintf(ts.URL),
		httpClient: ts.Client(),
	}
}

func NewTestClientFromFilename(t *testing.T, filename string) (*Client, *httptest.Server) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatalf("Could not open test file : %s", err)
		}
		w.Write(page)
	}))
	return NewTestClient(ts), ts
}

// tests
func TestNewClient(t *testing.T) {
	client := NewClient("test")
	want := "https://test@api.sncf.com/v1"
	if client.baseURL != want {
		t.Fatal("Invalid new client")
	}
}
