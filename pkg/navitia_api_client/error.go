package navitia_api_client

import "fmt"

// navitia api query error
type ApiError struct {
	code    int
	request string
}

func (e ApiError) Error() string {
	return fmt.Sprintf("Navitia Api error return code %d - %s", e.code, e.request)
}

func newApiError(code int, request string) error {
	return ApiError{
		code:    code,
		request: request,
	}
}

// http client error
type HttpClientError struct {
	msg string
	err error
}

func (e HttpClientError) Error() string { return fmt.Sprintf("Navitia HttpClient error %s", e.msg) }
func (e HttpClientError) Unwrap() error { return e.err }

func newHttpClientError(msg string, err error) error {
	return HttpClientError{
		msg: msg,
		err: err,
	}
}

// json decoding error
type JsonDecodeError struct {
	msg string
	err error
}

func (e JsonDecodeError) Error() string { return fmt.Sprintf("Navitia JsonDecode error %s", e.msg) }
func (e JsonDecodeError) Unwrap() error { return e.err }

func newJsonDecodeError(msg string, err error) error {
	return JsonDecodeError{
		msg: msg,
		err: err,
	}
}

// date parsing error
type DateParsingError struct {
	date string
	err  error
}

func (e DateParsingError) Error() string {
	return fmt.Sprintf("Navitia date parsing error %s", e.date)
}
func (e DateParsingError) Unwrap() error { return e.err }

func newDateParsingError(date string, err error) error {
	return DateParsingError{
		date: date,
		err:  err,
	}
}
