package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Concurrency int `yaml:"concurrency"`
	Timeouts    struct {
		CommandSeconds int `yaml:"command_seconds"`
	} `yaml:"timeouts"`
	Output struct {
		JSON bool `yaml:"json"`
	} `yaml:"output"`
	Paths struct {
		Python   string `yaml:"python"`
		Nmap     string `yaml:"nmap"`
		Nslookup string `yaml:"nslookup"`
		Whois    string `yaml:"whois"`
	} `yaml:"paths"`
}

func Default() Config {
	cfg := Config{Concurrency: 4}
	cfg.Timeouts.CommandSeconds = 120
	cfg.Output.JSON = false
	cfg.Paths.Python = "python3"
	cfg.Paths.Nmap = "nmap"
	cfg.Paths.Nslookup = "nslookup"
	cfg.Paths.Whois = "whois"
	return cfg
}

func Load(path string) (Config, error) {
	cfg := Default()
	resolved, err := resolvePath(path)
	if err != nil {
		return cfg, err
	}
	if resolved == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(resolved)
	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

func resolvePath(path string) (string, error) {
	if path != "" {
		return path, nil
	}
	if env := os.Getenv("CT_CONFIG"); env != "" {
		return env, nil
	}

	local := filepath.Join("configs", "default.yaml")
	if _, err := os.Stat(local); err == nil {
		return local, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return "", nil
}
