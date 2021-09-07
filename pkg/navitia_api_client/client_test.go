package navitia_api_client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// package utilities
func newTestClient(ts *httptest.Server) *NavitiaClient {
	return &NavitiaClient{
		baseURL:    fmt.Sprintf(ts.URL),
		httpClient: ts.Client(),
		cache:      make(map[string]cachedResult),
	}
}

func newTestClientFromFilename(t *testing.T, filename string) (*NavitiaClient, *httptest.Server) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatalf("Could not open test file : %s", err)
		}
		w.Write(page)
	}))
	return newTestClient(ts), ts
}

type testClientCase struct {
	urlPath  string
	filename string
}

func newTestClientFromFilenames(t *testing.T, tcs []testClientCase) (*NavitiaClient, *httptest.Server) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		found := false
		for i := 0; i < len(tcs) && !found; i++ {
			if r.URL.Path+"?"+r.URL.RawQuery == tcs[i].urlPath {
				page, err := ioutil.ReadFile(tcs[i].filename)
				if err != nil {
					t.Fatalf("Could not open test file : %s", err)
				}
				w.Write(page)
				found = true
			}
		}
		if !found {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}))
	return newTestClient(ts), ts
}
