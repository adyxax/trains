package navitia_api_client

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"git.adyxax.org/adyxax/trains/pkg/model"
)

type Client interface {
	GetDepartures(trainStop string) (departures []model.Departure, err error)
}

type NavitiaClient struct {
	baseURL    string
	httpClient *http.Client

	mutex sync.Mutex
	cache map[string]cachedResult
}

type cachedResult struct {
	ts     time.Time
	result interface{}
}

func NewClient(token string) Client {
	return &NavitiaClient{
		baseURL: fmt.Sprintf("https://%s@api.sncf.com/v1", token),
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
		cache: make(map[string]cachedResult),
	}
}
