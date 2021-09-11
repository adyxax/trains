package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFile(t *testing.T) {
	// Minimal yaml file
	minimalConfig := Config{
		Address: "127.0.0.1",
		Port:    "8080",
		Token:   "12345678-9abc-def0-1234-56789abcdef0",
	}

	// Minimal yaml file with hostname resolving
	minimalConfigWithResolving := Config{
		Address: "localhost",
		Port:    "www",
		Token:   "12345678-9abc-def0-1234-56789abcdef0",
	}

	// Complete yaml file
	completeConfig := Config{
		Address: "127.0.0.2",
		Port:    "8082",
		Token:   "12345678-9abc-def0-1234-56789abcdef0",
	}

	// Test cases
	testCases := []struct {
		name          string
		input         string
		expected      *Config
		expectedError interface{}
	}{
		{"Non existant file", "test_data/non-existant", nil, &OpenError{}},
		{"Invalid file content", "test_data/invalid.yaml", nil, &DecodeError{}},
		{"Invalid address should fail to load", "test_data/invalid_address.yaml", nil, &InvalidAddressError{}},
		{"Unresolvable address should fail to load", "test_data/invalid_address_unresolvable.yaml", nil, &InvalidAddressError{}},
		{"Invalid port should fail to load", "test_data/invalid_port.yaml", nil, &InvalidPortError{}},
		{"Invalid token should fail to load", "test_data/invalid_token.yaml", nil, &InvalidTokenError{}},
		{"Minimal config", "test_data/minimal.yaml", &minimalConfig, nil},
		{"Minimal config with resolving", "test_data/minimal_with_hostname.yaml", &minimalConfigWithResolving, nil},
		{"Complete config", "test_data/complete.yaml", &completeConfig, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := LoadFile(tc.input)
			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equalf(t, reflect.TypeOf(err), reflect.TypeOf(tc.expectedError), "Invalid error type. Got %s but expected %s", reflect.TypeOf(err), reflect.TypeOf(tc.expectedError))
				require.Nil(t, valid)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.expected, valid, "Invalid value")
		})
	}
}
