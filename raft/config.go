// config.go

package raft

import (
	"flag"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

//LoggerConfig logger config
type LoggerConfig struct {
	LogFilePath string `yaml:"LogFilePath"` //
	LogLevel    string `yaml:"LogLevel"`    //
}

//RaftConfig raft config
type RaftConfig struct {
	Id          int        `yaml:"Id"`
	PeerTcpPort int        `yaml:"PeerTcpPort"` //
	Peers       []RaftPeer `yaml:"Peers"`
}

//Config
type Config struct {
	Raft   RaftConfig   `yaml:"raft"`
	Logger LoggerConfig `yaml:"logger"`
}

var (
	logger     = stdLogger
	configFile = flag.String("c", "", "toy raft config file")
)

//DefaultConfig get config with default value
func DefaultConfig() *Config {
	return &Config{
		Raft: RaftConfig{
			PeerTcpPort: DEF_PEER_TCP_PORT,
		},
	}
}

//LoadConfig load config from configPath
func LoadConfig(configPath string) (*Config, error) {
	confYaml, err := ioutil.ReadFile(*configFile)
	if err != nil {
		logger.Fatal("config file %v error: %v", *configFile, err)
		return nil, err
	}

	config := DefaultConfig()
	err = yaml.Unmarshal(confYaml, config)
	if err != nil {
		logger.Fatal("config file %v error: %v", *configFile, err)
		return nil, err
	}

	logger.SetLogLevel(config.Logger.LogLevel)

	return config, nil
}
