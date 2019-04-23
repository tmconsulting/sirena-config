package sirenaConfig

import (
	"io/ioutil"
	"os"
)

type SirenaConfig struct {
	ClientID                 string `yaml:"client_id,omitempty"`
	Host                     string `yaml:"host,omitempty"`
	Port                     string `yaml:"port,omitempty"`
	ClientPublicKey          string `yaml:"client_public_key,omitempty"`
	ClientPrivateKey         string `yaml:"client_private_key,omitempty"`
	ClientPrivateKeyPassword string `yaml:"client_private_key_password,omitempty"`
	ServerPublicKey          string `yaml:"server_public_key,omitempty"`
	KeysPath                 string `yaml:"key_path"`
}

// GetSirenaAddr return sirena address to connect client to
func (config *SirenaConfig) GetSirenaAddr() string {
	if config == nil {
		return ""
	}
	if config.Port == "" {
		return config.Host
	}
	return config.Host + ":" + config.Port
}

// GetKeyFile returns contents of key file
func (config *SirenaConfig) GetKeyFile(keyFile string) ([]byte, error) {
	KeyPath := config.KeysPath + "/" + keyFile
	if _, err := os.Stat(KeyPath); os.IsNotExist(err) {
		return nil, err
	}

	return ioutil.ReadFile(KeyPath)
}
