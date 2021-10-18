// config.go

package raft

import (
	"flag"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ServerListenPort int    `yaml:"ListenPort"`  //
	PeerTcpPort      int    `yaml:"PeerTcpPort"` //
	LogFilePath      string `yaml:"LogFilePath"` //
	LogLevel         string `yaml:"LogLevel"`    //
}

var (
	logger     = stdLogger
	configFile = flag.String("c", "", "toy raft config file")
)

func DefaultConfig() *Config {
	return &Config{
		ServerListenPort: CONF_DEF_LISTEN_PORT,
		PeerTcpPort:      CONF_DEF_PEER_TCP_PORT,
	}
}

func LoadConfig(configPath string) *Config {
	config := DefaultConfig()
	confYaml, err := ioutil.ReadFile(*configFile)
	if err != nil {
		logger.Fatal("config file %v error: %v", *configFile, err)
	}

	err = yaml.Unmarshal(confYaml, config)
	if err != nil {
		logger.Fatal("config file %v error: %v", *configFile, err)
	}

	logger.SetLogLevel(config.LogLevel)

	return config
}
