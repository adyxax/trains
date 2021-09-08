package navitia_api_client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.adyxax.org/adyxax/trains/pkg/model"
)

type StopsResponse struct {
	Pagination struct {
		StartPage    int `json:"start_page"`
		ItemsOnPage  int `json:"items_on_page"`
		ItemsPerPage int `json:"items_per_page"`
		TotalResult  int `json:"total_result"`
	} `json:"pagination"`
	StopAreas []struct {
		Name                 string        `json:"name"`
		ID                   string        `json:"id"`
		Codes                []interface{} `json:"codes"`
		Links                []interface{} `json:"links"`
		Coord                interface{}   `json:"coord"`
		Label                string        `json:"label"`
		Timezone             interface{}   `json:"timezone"`
		AdministrativeRegion interface{}   `json:"administrative_regions"`
	} `json:"stop_areas"`
	Links          []interface{} `json:"links"`
	Disruptions    []interface{} `json:"disruptions"`
	FeedPublishers []interface{} `json:"feed_publishers"`
	Context        interface{}   `json:"context"`
}

func (c *NavitiaClient) GetStops() (trainStops []model.Stop, err error) {
	return getStopsPage(c, 0)
}

func getStopsPage(c *NavitiaClient, i int) (trainStops []model.Stop, err error) {
	request := fmt.Sprintf("%s/coverage/sncf/stop_areas?count=1000&start_page=%d", c.baseURL, i)
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		return nil, newHttpClientError("http.NewRequest error", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, newHttpClientError("httpClient.Do error", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var data StopsResponse
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, newJsonDecodeError("GetStops ", err)
		}
		for i := 0; i < len(data.StopAreas); i++ {
			if data.StopAreas[i].Label != "" {
				trainStops = append(trainStops, model.Stop{data.StopAreas[i].ID, data.StopAreas[i].Label})
			}
		}
		if data.Pagination.ItemsOnPage+data.Pagination.ItemsPerPage*data.Pagination.StartPage < data.Pagination.TotalResult {
			tss, err := getStopsPage(c, i+1)
			if err != nil {
				return nil, err
			}
			trainStops = append(trainStops, tss...)
		}
	} else {
		err = newApiError(resp.StatusCode, "GetStops")
	}
	return
}
