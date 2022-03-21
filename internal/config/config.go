package config

import (
	"fmt"
	"os"
	"strconv"

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

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid value (%s) for port", os.Getenv("PORT"))
	}

	amqpPort, err := strconv.Atoi(os.Getenv("AMQP_PORT"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid value (%s) for AMQP port", os.Getenv("AMQP_PORT"))
	}

	dbPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid value (%s) for postgres port", os.Getenv("POSTGRES_PORT"))
	}

	return Config{
		Port: port,
		AMQPConfig: AMQPConfig{
			Host:     os.Getenv("AQMP_HOST"),
			Port:     amqpPort,
			Username: os.Getenv("AMQP_USERNAME"),
			Password: os.Getenv("AMQP_PASSWORD"),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     dbPort,
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DATABASE"),
		},
	}, nil
}
