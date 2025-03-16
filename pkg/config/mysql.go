package config

import (
	"encoding/json"
	"time"
)

type ConnectionOption struct {
	MaxIdleConns    int           `json:"max_idle_conns,omitempty"`
	MaxOpenConns    int           `json:"max_open_conns,omitempty"`
	ConnMaxLifeTime time.Duration `json:"conn_max_life_time,omitempty"`
}

type MysqlConfig struct {
	DSN        string            `json:"dsn"`
	ConnOption *ConnectionOption `json:"conn_option"`
}

func DefaultMysqlConfig() *MysqlConfig {
	return &MysqlConfig{}
}

func (c *MysqlConfig) Set(s string) error {
	content, err := getConfigContent(s)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *MysqlConfig) String() string {
	content, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func (c *MysqlConfig) Type() string {
	return "MysqlConfig"
}

func (c *MysqlConfig) Validate() error {
	return nil
}
