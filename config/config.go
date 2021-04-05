package config

import (
	"net"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

var validToken = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type Config struct {
	// Address is the hostname or ip the web server will listen to
	Address string `yaml:"address",default:"127.0.0.1"`
	Port    string `yaml:"port",default:"8080"`
	// Token is the sncf api token
	Token string `yaml:"token"`
}

func (c *Config) validate() error {
	// address
	if ip := net.ParseIP(c.Address); ip == nil {
		if _, err := net.LookupIP(c.Address); err != nil {
			return newInvalidAddressError(c.Address, err)
		}
	}
	// port
	if _, err := net.LookupPort("tcp", c.Port); err != nil {
		return newInvalidPortError(c.Port, err)
	}
	// token
	if ok := validToken.MatchString(c.Token); !ok {
		return newInvalidTokenError(c.Token)
	}
	return nil
}

// LoadFile loads the c from a given file
func LoadFile(path string) (*Config, error) {
	var c *Config
	f, errOpen := os.Open(path)
	if errOpen != nil {
		return nil, newOpenError(path, errOpen)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&c); err != nil {
		return nil, newDecodeError(path, err)
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c, nil
}
