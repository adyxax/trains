package config

import (
	"reflect"
	"testing"
)

func TestLoadFile(t *testing.T) {
	// Non existant file
	_, err := LoadFile("test_data/non-existant")
	if err == nil {
		t.Fatal("non-existant config file failed without error")
	}

	// Invalid yaml file
	_, err = LoadFile("test_data/invalid_yaml")
	if err == nil {
		t.Fatal("invalid_yaml config file failed without error")
	}

	// Invalid address
	if _, err = LoadFile("test_data/invalid_address.yaml"); err == nil {
		t.Fatal("Invalid address should fail to load")
	}

	// Invalid address unreasolvable
	if _, err = LoadFile("test_data/invalid_address_unresolvable.yaml"); err == nil {
		t.Fatal("Unresolvable address should fail to load")
	}

	// Invalid port
	if _, err = LoadFile("test_data/invalid_port.yaml"); err == nil {
		t.Fatal("Invalid port should fail to load")
	}

	// Invalid token
	if _, err = LoadFile("test_data/invalid_token.yaml"); err == nil {
		t.Fatal("Invalid token should fail to load")
	}

	// Minimal yaml file
	want := Config{
		Address: "127.0.0.2",
		Port:    "8082",
		Token:   "12345678-9abc-def0-1234-56789abcdef0",
	}
	config, err := LoadFile("test_data/minimal.yaml")
	if err != nil {
		t.Fatalf("minimal example failed with error: %v", err)
	}
	if config != nil && !reflect.DeepEqual(want, *config) {
		t.Fatalf("minimal example failed:\nwant:%+v\ngot: %+v", want, *config)
	}

	// Minimal yaml file with hostname resolving
	want = Config{
		Address: "localhost",
		Port:    "8082",
		Token:   "12345678-9abc-def0-1234-56789abcdef0",
	}
	config, err = LoadFile("test_data/minimal_with_hostname.yaml")
	if err != nil {
		t.Fatalf("minimal example failed with error: %v", err)
	}
	if config != nil && !reflect.DeepEqual(want, *config) {
		t.Fatalf("minimal example failed:\nwant:%+v\ngot: %+v", want, *config)
	}
}
