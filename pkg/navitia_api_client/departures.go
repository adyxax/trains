package navitia_api_client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"git.adyxax.org/adyxax/trains/pkg/model"
)

type DeparturesResponse struct {
	Disruptions []interface{} `json:"disruptions"`
	Notes       []interface{} `json:"notes"`
	Departures  []struct {
		DisplayInformations struct {
			Direction      string        `json:"direction"`
			Code           string        `json:"code"`
			Network        string        `json:"network"`
			Links          []interface{} `json:"links"`
			Color          string        `json:"color"`
			Name           string        `json:"name"`
			PhysicalMode   string        `json:"physical_mode"`
			Headsign       string        `json:"headsign"`
			Label          string        `json:"label"`
			Equipments     []interface{} `json:"equipments"`
			TextColor      string        `json:"text_color"`
			TripShortName  string        `json:"trip_short_name"`
			CommercialMode string        `json:"commercial_mode"`
			Description    string        `json:"description"`
		} `json:"display_informations"`
		StopDateTime struct {
			Links                  []interface{} `json:"links"`
			ArrivalDateTime        string        `json:"arrival_date_time"`
			AdditionalInformations []interface{} `json:"additional_informations"`
			DepartureDateTime      string        `json:"departure_date_time"`
			BaseArrivalDateTime    string        `json:"base_arrival_date_time"`
			BaseDepartureDateTime  string        `json:"base_departure_date_time"`
			DataFreshness          string        `json:"data_freshness"`
		} `json:"stop_date_time"`
	} `json:"departures"`
	Context struct {
		Timezone        string `json:"timezone"`
		CurrentDatetime string `json:"current_datetime"`
	} `json:"context"`
}

func (c *NavitiaClient) GetDepartures(trainStop string) (departures []model.Departure, err error) {
	request := fmt.Sprintf("%s/coverage/sncf/stop_areas/%s/departures", c.baseURL, trainStop)
	start := time.Now()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if cachedResult, ok := c.cache[request]; ok {
		if start.Sub(cachedResult.ts) < 60*1000*1000*1000 {
			return cachedResult.result.([]model.Departure), nil
		}
	}
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
		var data DeparturesResponse
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, newJsonDecodeError("GetDepartures "+trainStop, err)
		}
		// TODO test for no json error
		for i := 0; i < len(data.Departures); i++ {
			t, err := time.Parse("20060102T150405", data.Departures[i].StopDateTime.ArrivalDateTime)
			if err != nil {
				return nil, newDateParsingError(data.Departures[i].StopDateTime.ArrivalDateTime, err)
			}
			departures = append(departures, model.Departure{data.Departures[i].DisplayInformations.Direction, t.Format("Mon, 02 Jan 2006 15:04:05")})
		}
		c.cache[request] = cachedResult{
			ts:     start,
			result: departures,
		}
	} else {
		err = newApiError(resp.StatusCode, "GetDepartures "+trainStop)
	}
	return
}
