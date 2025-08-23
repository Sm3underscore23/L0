package config

import (
	"fmt"
	"net"
	"os"

	"gopkg.in/yaml.v3"
	customerrors "taskL0/internal/entity/custom_errors"
)

type (
	Config struct {
		HTTP     HTTP     `yaml:"http"`
		Cache    Cache    `yaml:"cache"`
		Log      Log      `yaml:"logger"`
		PG       PG       `yaml:"postgres"`
		Kafka    Kafka    `yaml:"kafka"`
	}

	HTTP struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Cache struct {
		Limit        int `yaml:"limit"`
		RecoverLimit int `yaml:"recover_limit"`
	}

	Log struct {
		Level string `yaml:"log_level"`
	}

	PG struct {
		Host          string `yaml:"host"`
		Port          string `yaml:"port"`
		Name          string `yaml:"db_name"`
		User          string `yaml:"user"`
		password      string
		SSL           string `yaml:"sslmode"`
		PoolMax       int32  `yaml:"pool_max"`
		MigrationsDir string `yaml:"migrations_dir"`
	}

	Kafka struct {
		ConsumerGroup string   `yaml:"consumer_group"`
		BrokerList    []string `yaml:"broker_list"`
		Topic         string   `yaml:"topic"`
		WorkersNum    int      `yaml:"workers_num"`
	}
)

// NewConfig returns app config.
func NewConfig(configPath string) (Config, error) {
	if _, err := os.Stat(configPath); err != nil {
		return Config{}, err
	}

	rowConfig, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	err = yaml.Unmarshal(rowConfig, &cfg)
	if err != nil {
		return Config{}, err
	}

	cfg.PG.password = os.Getenv("DB_PASSWORD")
	if cfg.PG.password == "" {
		return Config{}, customerrors.ErrDbPasswordEpt
	}

	return cfg, nil
}

// ServerAddress returns server address in "host:port" format.
func (c *Config) ServerAddress() string {
	return net.JoinHostPort(c.HTTP.Host, c.HTTP.Port)
}

// PostgresURl returns postgres url.
func (c *Config) PostgresURl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		c.PG.User,
		c.PG.password,
		c.PG.Host,
		c.PG.Port,
		c.PG.Name,
	)
}
