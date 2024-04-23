package main

import (
	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/server"
)

func main() {
	// Config
	load := config.LoadConfig("config-dev.yml")
	config := config.ParseConfig(load)

	// Server
	server := server.NewServer(config)

	server.Run()
}
