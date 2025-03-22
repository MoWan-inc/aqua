package config

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

type ApiConfig struct {
	Addr        string `json:"addr" validate:"required,url里的版本prefix"`
	EnablePProf bool   `json:"enable_pprof,omitempty"`
	// 优雅退出
	GracefullyShutDownSeconds int `json:"gracefully_shutdown_seconds,omitempty"`
	// url里的版本prefix
	Prefix string `json:"prefix,omitempty"`
	// swagger启动，用于生成网页api和生成client代码
	EnableSwagger bool `json:"enable_swagger,omitempty"`
	// tokens，配置里允许的内部token
	Tokens []string `json:"tokens,omitempty"`
}

func DefaultApiConfig() *ApiConfig {
	return &ApiConfig{
		Addr:                      "0.0.0.0:8080",
		GracefullyShutDownSeconds: 10,
		Prefix:                    "v1",
		Tokens:                    []string{""},
	}
}

func (c *ApiConfig) Validate() error {
	return validator.New().Struct(c)
}

func (c *ApiConfig) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (c *ApiConfig) Set(s string) error {
	content, err := getConfigContent(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, c)
}

func (c *ApiConfig) Type() string {
	return "apiConfig"
}
