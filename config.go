package config

import (
	"log"
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
