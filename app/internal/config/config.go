package config

import (
	"github.com/allokate-ai/environment"
	"github.com/joho/godotenv"
)

type AMQPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type Config struct {
	Port       int
	AMQPConfig AMQPConfig
	Database   DatabaseConfig
}

func Get() (Config, error) {
	godotenv.Load()

	return Config{
		Port: int(environment.GetIntOrDefault("PORT", 8094)),
		AMQPConfig: AMQPConfig{
			Host:     environment.GetValueOrDefault("AMQP_HOST", "localhost"),
			Port:     int(environment.GetIntOrDefault("AMQP_PORT", 5672)),
			Username: environment.GetValueOrDefault("AMQP_USERNAME", "guest"),
			Password: environment.GetValueOrDefault("AMQP_PASSWORD", "guest"),
		},
		Database: DatabaseConfig{
			Host:     environment.GetValueOrDefault("POSTGRES_HOST", "localhost"),
			Port:     int(environment.GetIntOrDefault("POSTGRES_PORT", 5432)),
			User:     environment.GetValueOrDefault("POSTGRES_USER", "root"),
			Password: environment.GetValueOrDefault("POSTGRES_PASSWORD", "example"),
			Database: environment.GetValueOrDefault("POSTGRES_DATABASE", "allokate"),
		},
	}, nil
}
