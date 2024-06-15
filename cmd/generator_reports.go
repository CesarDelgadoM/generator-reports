package main

import (
	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/server"
)

func main() {
	// Config
	config := config.GetConfig("config-dev.yml")

	// Server
	server := server.NewServer(config)

	server.Run()
}
