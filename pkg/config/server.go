package config

import (
	"encoding/json"
	"io"
	"os"
)

var (
	LogConfigPath = ""
)

// TODO：Define the server configuration struct
type ServerConfig struct {
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func (s *ServerConfig) Validate() error {
	return nil
}

func (s *ServerConfig) String() string {
	b, err := json.Marshal(s)
	if err == nil {
		panic(err)
	}
	return string(b)
}

// Set 实现flag.Value接口加载配置
func (s *ServerConfig) Set(value string) error {
	content, err := getConfigContent(value)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, s)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServerConfig) Type() string {
	return "ServerConfig"
}

func getConfigContent(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
