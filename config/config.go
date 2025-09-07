package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config holds application configuration. 支持从 env 或 env.yaml 加载。
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

// envNames lists supported environment variable names (preferred in QingLong).
// 使用示例：
// FAW_MOBILE, FAW_PASSWORD, FAW_WEBDID, FAW_BARK_PUSH_SERVER_URL, FAW_SECURITY_CODE, FAW_DID, FAW_EXEC_FREQ
func LoadConfig() (*Config, error) {
	// If environment variables are provided (common in 青龙), prefer them.
	// Detect by presence of FAW_MOBILE or FAW_PASSWORD or explicit FAW_USE_ENV=1
	useEnv := false
	if strings.TrimSpace(os.Getenv("FAW_USE_ENV")) == "1" {
		useEnv = true
	}
	if os.Getenv("FAW_MOBILE") != "" || os.Getenv("FAW_PASSWORD") != "" {
		useEnv = true
	}

	if useEnv {
		cfg := &Config{
			BarkPushServerURL: os.Getenv("FAW_BARK_PUSH_SERVER_URL"),
			Mobile:            os.Getenv("FAW_MOBILE"),
			Password:          os.Getenv("FAW_PASSWORD"),
			WebDid:            os.Getenv("FAW_WEBDID"),
			SecurityCode:      os.Getenv("FAW_SECURITY_CODE"),
			Did:               os.Getenv("FAW_DID"),
			ExecFreq:          os.Getenv("FAW_EXEC_FREQ"),
		}
		return cfg, nil
	}

	// Otherwise load from env.yaml (or file pointed by CONFIG_FILE)
	filename := os.Getenv("CONFIG_FILE")
	if filename == "" {
		filename = "env.yaml"
	}

	data, err := os.ReadFile(filename)
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
