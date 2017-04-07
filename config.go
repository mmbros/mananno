package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/naoina/toml"
)

type serverConfig struct {
	Host string
	Port int
}
type transmissionConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}
type configuration struct {
	Server       serverConfig
	Transmission transmissionConfig
}

func unmarshalConfig(data []byte) (*configuration, error) {
	var cfg configuration
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// loadConfigFromFile returns the configuration from a configuration file
func loadConfigFromFile(path string) (*configuration, error) {
	// open config file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// read config file
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return unmarshalConfig(buf)
}

func (t transmissionConfig) Address() string {
	return fmt.Sprintf("%s:%d", t.Host, t.Port)
}
func (t transmissionConfig) Web() string {
	return fmt.Sprintf("http://%s/transmission/web/", t.Address())
}

func (s serverConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
