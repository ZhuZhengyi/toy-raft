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
	ID               int      `yaml:"Id"`
	ServerListenPort int      `yaml:"ListenPort"`  //
	PeerTcpPort      int      `yaml:"PeerTcpPort"` //
	Peers            []string `yaml:"Peers"`
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
			ServerListenPort: CONF_DEF_LISTEN_PORT,
			PeerTcpPort:      CONF_DEF_PEER_TCP_PORT,
		},
	}
}

//LoadConfig load config from configPath
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()
	confYaml, err := ioutil.ReadFile(*configFile)
	if err != nil {
		logger.Fatal("config file %v error: %v", *configFile, err)
		return nil, err
	}

	err = yaml.Unmarshal(confYaml, config)
	if err != nil {
		logger.Fatal("config file %v error: %v", *configFile, err)
		return nil, err
	}

	logger.SetLogLevel(config.Logger.LogLevel)

	return config, nil
}
