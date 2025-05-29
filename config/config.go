package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// Bark推送的URL
	BarkPushServerURL string `yaml:"BarkPushServerURL"`

	// 登录
	Mobile string `yaml:"mobile"`
	// 登录
	Password string `yaml:"password"`
	// Web端Did
	WebDid string `yaml:"WebDid"`

	SecurityCode string `yaml:"securityCode"`
	Did          string `yaml:"did"`
	ExecFreq     string `yaml:"ExecFreq"`
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
