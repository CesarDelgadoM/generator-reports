package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   *ServerConfig
	Worker   *WorkerConfig
	Postgres *PostgresConfig
	RabbitMQ *RabbitMQ
	Consumer *Consumer
}

// Server config
type ServerConfig struct {
	Port string
}

type WorkerConfig struct {
	Pool int
	Idle int
}

// Postgres config
type PostgresConfig struct {
	Driver   string
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
}

// RabbitMQ config
type RabbitMQ struct {
	URI      string
	User     string
	Password string
}

type Consumer struct {
	DataBus *DataBus
	Branch  *Branch
}

type DataBus struct {
	ExchangeType string
	ContentType  string
}

type Branch struct {
	ExchangeType string
	ContentType  string
}

func LoadConfig(filename string) *viper.Viper {
	v := viper.New()

	v.SetConfigName(filename)
	v.SetConfigType("yaml")
	v.AddConfigPath("../config")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	return v
}

func ParseConfig(v *viper.Viper) *Config {
	var conf Config

	if err := v.Unmarshal(&conf); err != nil {
		log.Panic(err)
	}

	return &conf
}
