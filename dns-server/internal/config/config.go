package config

import (
	"errors"
	"net"
	"os"

	"github.com/goccy/go-yaml"
)


type Config struct {
	Listen string	`yaml:"listen"`
	TTL		uint32	`yaml:"ttl"`
	Upstream []string	`yaml:"upstream"`
	Records map[string]string `yaml:"records"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Listen == "" {
		cfg.Listen = ":53"
	}
	if cfg.TTL == 0 {
		cfg.TTL = 60
	}

	for name, ip := range cfg.Records {
		if net.ParseIP(ip) == nil {
			return nil, errors.New("invalid IP for" + name + ": " + ip)
		}
	}

	return &cfg, nil
}