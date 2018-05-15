package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/imdario/mergo"
	yaml "gopkg.in/yaml.v2"
)

// Config is a config :)
type Config struct {
	LogLevel                 string `yaml:"log_level,omitempty"`
	Addr                     string `yaml:"addr,omitempty"`
	TrackXML                 bool   `yaml:"track_xml,omitempty"`
	SirenaClientID           string `yaml:"sirena_client_id,omitempty"`
	SirenaHost               string `yaml:"sirena_host,omitempty"`
	SirenaPort               string `yaml:"sirena_port,omitempty"`
	ClientPublicKey          string `yaml:"client_public_key,omitempty"`
	ClientPrivateKey         string `yaml:"client_private_key,omitempty"`
	ClientPrivateKeyPassword string `yaml:"client_private_key_password,omitempty"`
	ServerPublicKey          string `yaml:"server_public_key,omitempty"`
	RedisHost                string `yaml:"redis_host,omitempty"`
	RedisPort                string `yaml:"redis_port,omitempty"`
}

// CNFG is a Config singletone
var CNFG *Config

func init() {
	loadConfig()
}

// Get returns config
func Get() *Config {
	return CNFG
}

// getEnv return env variable or default value provided
func getEnv(name, defaultVal string) string {
	val := os.Getenv("name")
	if val != "" {
		return val
	}
	return defaultVal
}

// loadConfig loads config from YAML files
func loadConfig() *Config {
	configPath := getEnv("CONFIG_PATH", "config")
	stage := getEnv("STAGE", "development")

	yamlFileList := []string{}
	err := filepath.Walk(configPath, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yaml" {
			return nil
		}
		yamlFileList = append(yamlFileList, path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	loadedConfigs := map[string]Config{}
	for _, yamlFilePath := range yamlFileList {
		yamlFileBytes, err := ioutil.ReadFile(yamlFilePath)
		if err != nil {
			log.Fatal(err)
		}
		fileConfig := map[string]Config{}
		err = yaml.Unmarshal(yamlFileBytes, &fileConfig)
		if err != nil {
			log.Fatal(err)
		}
		mergo.Merge(&loadedConfigs, fileConfig)
	}
	if _, ok := loadedConfigs[stage]; !ok {
		d := loadedConfigs["defaults"]
		return &d
	}
	if defaultConfig, ok := loadedConfigs["defaults"]; ok {
		CNFG = &defaultConfig
	}
	mergo.Merge(CNFG, loadedConfigs[stage], mergo.WithOverride)

	return CNFG
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

// KeyDirs is a list of directories to search Sirena key files in
var KeyDirs = []string{
	os.Getenv("GOPATH"),
	binaryDir() + "/keys",
	pwdDir() + "/sirena-agent-go/keys",
}

// GetKeyFile returns contents of key file
func (config *Config) GetKeyFile(keyFile string) ([]byte, error) {
	for _, keyDir := range KeyDirs {
		exists, err := pathExists(keyDir + "/" + keyFile)
		if err != nil {
			log.Print(err)
		}
		if !exists {
			continue
		}
		return ioutil.ReadFile(keyDir + "/" + keyFile)
	}
	return nil, fmt.Errorf("No key file %s found", keyFile)
}

// binaryDir returns path where binary was run from
func binaryDir() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(ex)
}

// pwdDir returns pwd dir
func pwdDir() string {
	ex, err := os.Getwd()
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
