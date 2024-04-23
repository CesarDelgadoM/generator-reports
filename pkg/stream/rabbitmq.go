package stream

import (
	"fmt"

	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ExchangeOpts struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

type QueueOpts struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type BindOpts struct {
	Name     string
	Key      string
	Exchange string
	NoWait   bool
	Args     amqp.Table
}

type ConsumeOpts struct {
	Name      string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

type RabbitMQ struct {
	*amqp.Connection
}

func ConnectRabbitMQ(config *config.RabbitMQ) *RabbitMQ {
	conn, err := amqp.Dial(fmt.Sprintf(config.URI, config.User, config.Password))
	if err != nil {
		zap.Log.Fatal("Connect to rabbitmq failed: ", err)
	}

	zap.Log.Info("Connection to rabbitmq success")

	return &RabbitMQ{
		conn,
	}
}

func (rmq *RabbitMQ) OpenChannel() *amqp.Channel {
	ch, err := rmq.Channel()
	if err != nil {
		zap.Log.Fatal("Opened channel failed: ", err)
	}

	return ch
}
