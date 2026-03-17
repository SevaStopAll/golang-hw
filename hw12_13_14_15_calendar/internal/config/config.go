package config

import (
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Logger    LoggerConfig    `yaml:"logger"`
	Storage   StorageConfig   `yaml:"storage"`
	Kafka     KafkaConfig     `yaml:"kafka"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type StorageConfig struct {
	Type string `yaml:"type"`
	DSN  string `yaml:"dsn"`
}

type KafkaConfig struct {
	Brokers      []string      `yaml:"brokers"`
	Topic        string        `yaml:"topic"`
	GroupID      string        `yaml:"group_id"`
	ClientID     string        `yaml:"client_id"`
	MaxAttempts  int           `yaml:"max_attempts"`
	RetryBackoff time.Duration `yaml:"retry_backoff"`
}

type SchedulerConfig struct {
	Interval         time.Duration `yaml:"interval"`
	CleanupOlderThan time.Duration `yaml:"cleanup_older_than"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
