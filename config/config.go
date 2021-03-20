package config

import (
	"net"
	"os"
	"regexp"

	"github.com/pkg/errors"
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
			return errors.New("Invalid address " + c.Address + ", it must be a valid ipv4 address, ipv6 address, or resolvable name.")
		}
	}
	// port
	if _, err := net.LookupPort("tcp", c.Port); err != nil {
		return errors.New("Invalid port " + c.Port + ", it must be a valid port number or tcp service name. Got error : " + err.Error())
	}
	// token
	if ok := validToken.MatchString(c.Token); !ok {
		return errors.New("Invalid token, must be an hexadecimal string that lookslike 12345678-9abc-def0-1234-56789abcdef0, got " + c.Token + " instead.")
	}
	return nil
}

// LoadFile loads the c from a given file
func LoadFile(path string) (*Config, error) {
	var c *Config
	f, errOpen := os.Open(path)
	if errOpen != nil {
		return nil, errors.Wrapf(errOpen, "Failed to open configuration file %s", path)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&c); err != nil {
		return nil, errors.Wrap(err, "Failed to decode configuration file")
	}
	if err := c.validate(); err != nil {
		return nil, errors.Wrap(err, "Failed to validate configuration")
	}
	return c, nil
}
