package config

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/jinzhu/configor"
)

// Config is a config :)
type Config struct {
	LogLevel                 string `yaml:"log_level" env:"LOG_LEVEL" default:"debug"` // log everything by default
	Addr                     string `yaml:"addr" env:"SSM_ADDR" required:"true"`
	TrackXML                 bool   `yaml:"track_xml" env:"TRACK_XML"`
	SirenaClientID           string `yaml:"sirena_client_id" env:"SIRENA_CLIENT_ID" required:"true"`
	SirenaHost               string `yaml:"sirena_host" env:"SIRENA_HOST" required:"true"`
	SirenaPort               string `yaml:"sirena_port" env:"SIRENA_PORT" required:"true"`
	ClientPublicKey          string `yaml:"client_public_key" env:"CLIENT_PUBLIC_KEY" required:"true"`
	ClientPrivateKey         string `yaml:"client_private_key" env:"CLIENT_PRIVATE_KEY" required:"true"`
	ClientPrivateKeyPassword string `yaml:"client_private_key_password" env:"CLIENT_PRIVATE_KEY_PASSWORD"`
	ServerPublicKey          string `yaml:"server_public_key" env:"CLIENT_PUBLIC_KEY" required:"true"`
	EnvType                  string
}

var config = &Config{}

// Singleton guard
var once sync.Once

// Get reads config from environment or JSON
func Get() *Config {
	once.Do(func() {
		if err := configor.New(&configor.Config{Debug: true}).Load(config, "config/config.yaml"); err != nil {
			log.Fatal(err)
		}
	})
	return config
}

// GetSirenaAddr return sirena address to connect client to
func (config *Config) GetSirenaAddr() string {
	if config == nil {
		return ""
	}
	if config.SirenaPort == "" {
		return config.SirenaHost
	}
	return config.SirenaHost + ":" + config.SirenaPort
}

// GetKeyFile returns contents of key file
func (config *Config) GetKeyFile(keyFile string) ([]byte, error) {
	keyDirs := []string{
		os.Getenv("GOPATH"),
		binaryDir() + "/keys",
	}
	for _, keyDir := range keyDirs {
		exists, err := pathExists(keyDir + "/" + keyFile)
		if err != nil {
			log.Print(err)
		}
		if !exists {
			continue
		}
		return ioutil.ReadFile(keyDir + "/" + keyFile)
	}
	return nil, errors.New("No key files found")
}

// binaryDir returns path where binary was run from
func binaryDir() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(ex)
}

// pathExists checks if file or dir exist
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
