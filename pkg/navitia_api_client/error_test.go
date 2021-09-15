package navitia_api_client

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireErrorTypeMatch(t *testing.T, err error, expected error) {
	require.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(expected), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(expected))
}

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
