package navitia_api_client

import "testing"

func TestErrorsCoverage(t *testing.T) {
	apiErr := ApiError{}
	_ = apiErr.Error()
	httpClientErr := HttpClientError{}
	_ = httpClientErr.Error()
	_ = httpClientErr.Unwrap()
	jsonDecodeErr := JsonDecodeError{}
	_ = jsonDecodeErr.Error()
	_ = jsonDecodeErr.Unwrap()
	dateParsingErr := DateParsingError{}
	_ = dateParsingErr.Error()
	_ = dateParsingErr.Unwrap()
}
