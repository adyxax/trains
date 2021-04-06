package config

import "fmt"

type ErrorType int

// file open configuration file error
type OpenError struct {
	path string
	err  error
}

func (e *OpenError) Error() string {
	return fmt.Sprintf("Failed to open configuration file : %s", e.path)
}
func (e *OpenError) Unwrap() error { return e.err }

func newOpenError(path string, err error) error {
	return &OpenError{
		path: path,
		err:  err,
	}
}

// Yaml configuration file decoding error
type DecodeError struct {
	path string
	err  error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("Failed to decode configuration file : %s", e.path)
}
func (e *DecodeError) Unwrap() error { return e.err }

func newDecodeError(path string, err error) error {
	return &DecodeError{
		path: path,
		err:  err,
	}
}

// Invalid address field error
type InvalidAddressError struct {
	address string
	err     error
}

func (e *InvalidAddressError) Error() string {
	return fmt.Sprintf("Invalid address %s : it must be a valid ipv4 address, ipv6 address, or resolvable name", e.address)
}
func (e *InvalidAddressError) Unwrap() error { return e.err }

func newInvalidAddressError(address string, err error) error {
	return &InvalidAddressError{
		address: address,
		err:     err,
	}
}

// Invalid port field error
type InvalidPortError struct {
	port string
	err  error
}

func (e *InvalidPortError) Error() string {
	return fmt.Sprintf("Invalid port %s : it must be a valid port number or tcp service name", e.port)
}
func (e *InvalidPortError) Unwrap() error { return e.err }

func newInvalidPortError(port string, err error) error {
	return &InvalidPortError{
		port: port,
		err:  err,
	}
}

// Invalid token field error
type InvalidTokenError struct {
	token string
}

func (e *InvalidTokenError) Error() string {
	return fmt.Sprintf("Invalid token %s : it must be an hexadecimal string that lookslike 12345678-9abc-def0-1234-56789abcdef0", e.token)
}

func newInvalidTokenError(token string) error {
	return &InvalidTokenError{
		token: token,
	}
}

// Invalid trainStop field error
type InvalidTrainStopError struct {
	trainStop string
}

func (e *InvalidTrainStopError) Error() string {
	return fmt.Sprintf("Invalid trainStop %s : it must be a string that lookslike \"stop_area:SNCF:87723502\" (make sure to quote the string because of the colon characters)", e.trainStop)
}

func newInvalidTrainStopError(trainStop string) error {
	return &InvalidTrainStopError{
		trainStop: trainStop,
	}
}
