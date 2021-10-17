// config.go

package raft

type Config struct {
	ServerListenPort int    `json:"ListenPort" yaml:"ListenPort"`   //
	PeerTcpPort      int    `json:"PeerTcpPort" yaml:"PeerTcpPort"` //
	LogFilePath      string `json:"LogFilePath" yaml:"LogFilePath"` //
	LogLevel         string `json:"LogLevel" yaml:"LogLevel"`       //
}

func DefaultConfig() *Config {
	return &Config{
		ServerListenPort: CONF_DEF_LISTEN_PORT,
		PeerTcpPort:      CONF_DEF_PEER_TCP_PORT,
	}
}

func ParseConfig(configPath string) *Config {
	config := DefaultConfig()

	return config
}
