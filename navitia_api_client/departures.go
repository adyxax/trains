package navitia_api_client

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (c *Client) GetDepartures() (departures *DeparturesResponse, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/coverage/sncf/stop_areas/stop_area:SNCF:87723502/departures", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&departures); err != nil {
		return nil, err
	}
	return
}
