package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// 定义您的配置字段
	Port   int    `yaml:"port"`
	Schema string `yaml:"schema"`

	// Bark推送的URL
	BarkPushServerURL string `yaml:"BarkPushServerURL"`

	// 登录
	Mobile string `yaml:"mobile"`
	// 登录
	Password string `yaml:"password"`
	// 登录，签到请求中使用
	SecurityCode string `yaml:"securityCode"`
	Did          string
}

func LoadConfig() (*Config, error) {
	// var env string
	// // 如果是 debug 模式（通常在本地开发）
	// if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
	// 	env = "dev"
	// } else {
	// 	env = "prod"
	// }

	filename := fmt.Sprintf("env.yaml")
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
