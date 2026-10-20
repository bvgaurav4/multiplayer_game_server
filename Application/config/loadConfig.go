package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Database DatabaseConfig `toml:"database"`
}

type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Name     string `toml:"name"`
	SSLMode  string `toml:"sslmode"`

	Pool PoolConfig `toml:"pool"`
}

type PoolConfig struct {
	MaxOpenConns           int `toml:"max_open_conns"`
	MaxIdleConns           int `toml:"max_idle_conns"`
	ConnMaxIdleTimeMinutes int `toml:"conn_max_idle_time_minutes"`
	ConnMaxLifetimeMinutes int `toml:"conn_max_lifetime_minutes"`
}

func Load(path string) (*Config, error) {
	var cfg Config

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return &cfg, nil
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

func LoadFromEnv() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config.toml"
	}

	return Load(path)
}
