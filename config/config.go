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
	DataBus  *DataBus
	Branch   *Branch
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
	ExchangeType string
	ContentType  string
}

type DataBus struct {
	Consumer *Consumer
}

type Branch struct {
	Consumer *Consumer
	Pdf      *PDF
}

type PDF struct {
	Path   string
	Suffix string
	Font   string
	Title  string
}

func GetConfig(filename string) *Config {
	load := loadConfig("config-dev.yml")
	return parseConfig(load)
}

func loadConfig(filename string) *viper.Viper {
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

func parseConfig(v *viper.Viper) *Config {
	var conf Config

	if err := v.Unmarshal(&conf); err != nil {
		log.Panic(err)
	}

	return &conf
}
