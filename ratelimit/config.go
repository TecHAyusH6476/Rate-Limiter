package ratelimit

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RateLimitConfig struct {
	RateLimits []RateLimitRule `yaml:"rate_limits"`
}

type RateLimitRule struct {
	Domain      string            `yaml:"domain"`
	Descriptors map[string]string `yaml:"descriptors"`
	RateLimit   struct {
		Unit            string `yaml:"unit"`
		RequestsPerUnit int    `yaml:"requests_per_unit"`
	} `yaml:"rate_limit"`
}

func LoadConfig(path string) (*RateLimitConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg RateLimitConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
