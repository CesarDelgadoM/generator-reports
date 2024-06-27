package server

import (
	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/internal/consumer/databus"
	"github.com/CesarDelgadoM/generator-reports/internal/generators/branch"
	"github.com/CesarDelgadoM/generator-reports/internal/repositorys"
	"github.com/CesarDelgadoM/generator-reports/internal/workerpool"
	"github.com/CesarDelgadoM/generator-reports/pkg/database"
	"github.com/CesarDelgadoM/generator-reports/pkg/email"
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

	// Databases
	postgresdb := database.ConnectPostgresDB(s.config.Postgres)

	// SMTP
	smtp := email.NewEmailSMTP(s.config)

	// Repositorys
	userRepo := repositorys.NewUserRepository(postgresdb)

	// Emails
	email := branch.NewEmailBranch(s.config, smtp)

	// Consumers
	branchConsumer := branch.NewBranchConsumer(s.config, rabbitmq, email, userRepo)

	// DataBus
	databus := databus.NewDataBusConsumer(s.config, rabbitmq, workerpool, branchConsumer)

	// Launch main consumer
	databus.StartDataBusConsumer()

	// App
	app := fiber.New()

	// Launch
	app.Listen(s.config.Server.Port)
}
