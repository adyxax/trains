package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireErrorTypeMatch(t *testing.T, err error, expected error) {
	require.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(expected), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(expected))
}

func TestErrorsCoverage(t *testing.T) {
	openErr := OpenError{}
	_ = openErr.Error()
	_ = openErr.Unwrap()
	decodeErr := DecodeError{}
	_ = decodeErr.Error()
	_ = decodeErr.Unwrap()
	invalidAddressErr := InvalidAddressError{}
	_ = invalidAddressErr.Error()
	_ = invalidAddressErr.Unwrap()
	invalidPortErr := InvalidPortError{}
	_ = invalidPortErr.Error()
	_ = invalidPortErr.Unwrap()
	invalidTokenErr := InvalidTokenError{}
	_ = invalidTokenErr.Error()
}
