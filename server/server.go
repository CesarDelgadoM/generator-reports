package server

import (
	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/internal/consumer/databus"
	"github.com/CesarDelgadoM/generator-reports/internal/workerpool"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
	"github.com/CesarDelgadoM/generator-reports/pkg/stream"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Run() {
	// Logger
	zap.InitLogger(s.config)

	// Stream
	rabbitmq := stream.ConnectRabbitMQ(s.config.RabbitMQ)
	defer rabbitmq.Close()

	// Workerpool
	workerpool := workerpool.NewWorkerPool(s.config.Worker)

	// DataBus
	databus := databus.NewDataBus(s.config.Consumer.DataBus, rabbitmq, workerpool)

	databus.ConsumeQueueNames()

	// App
	app := fiber.New()

	// Launch
	app.Listen(s.config.Server.Port)
}
