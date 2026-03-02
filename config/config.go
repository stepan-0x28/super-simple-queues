package config

import "gopkg.in/ini.v1"

type Config struct {
	TCPPort  int
	HTTPPort int
}

func LoadConfig(path string) (Config, error) {
	cfgFile, err := ini.Load(path)

	if err != nil {
		return Config{}, err
	}

	serverSection := cfgFile.Section("servers")

	tcpPort, err := serverSection.Key("tcp_port").Int()

	if err != nil {
		return Config{}, err
	}

	httpPort, err := serverSection.Key("http_port").Int()

	if err != nil {
		return Config{}, err
	}

	return Config{tcpPort, httpPort}, nil
}
