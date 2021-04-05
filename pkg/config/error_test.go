package config

import "testing"

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
