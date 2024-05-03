package consumer

import (
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
	"github.com/CesarDelgadoM/generator-reports/pkg/stream"
	amqp "github.com/rabbitmq/amqp091-go"
)

type IConsumer interface {
	Exchange(opts *stream.ExchangeOpts)
	BindQueue(opts *stream.BindOpts)
	Queue(opts *stream.QueueOpts) amqp.Queue
	QueueDelete(opts *stream.QueueDelete)
	Consume(opts *stream.ConsumeOpts) <-chan amqp.Delivery
	Cancel(consumer string, noWait bool)
	Close()
}

type ConsumerOpts struct {
	ExchangeName string
	ExchangeType string
	ContentType  string
}

type consumer struct {
	opts *ConsumerOpts
	ch   *amqp.Channel
}

func NewConsumer(opts *ConsumerOpts, rabbit *stream.RabbitMQ) IConsumer {
	c := &consumer{
		opts: opts,
		ch:   rabbit.OpenChannel(),
	}
	return c
}

func (c *consumer) Exchange(opts *stream.ExchangeOpts) {
	err := c.ch.ExchangeDeclare(
		opts.Name,
		opts.Kind,
		opts.Durable,
		opts.AutoDelete,
		opts.Internal,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		zap.Log.Error("Failed to create exchange: ", err)
	}
}

func (c *consumer) BindQueue(opts *stream.BindOpts) {
	err := c.ch.QueueBind(
		opts.Name,
		opts.Key,
		c.opts.ExchangeName,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		zap.Log.Info("Failed to create queue bind: ", err)
	}
}

func (c *consumer) Queue(opts *stream.QueueOpts) amqp.Queue {
	queue, err := c.ch.QueueDeclare(
		opts.Name,
		opts.Durable,
		opts.AutoDelete,
		opts.Exclusive,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		zap.Log.Error("Failed to create queue: ", err)
	}

	return queue
}

func (c *consumer) QueueDelete(opts *stream.QueueDelete) {
	_, err := c.ch.QueueDelete(opts.Name, opts.IfUnused, opts.IfEmpty, opts.NoWait)
	if err != nil {
		zap.Log.Error("Failed to remove queue: ", err)
	}
}

func (c *consumer) Consume(opts *stream.ConsumeOpts) <-chan amqp.Delivery {
	msgs, err := c.ch.Consume(
		opts.Name,
		opts.Consumer,
		opts.AutoAck,
		opts.Exclusive,
		opts.NoLocal,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		zap.Log.Error("Failed to initialize consume: ", err)
	}

	return msgs
}

func (c *consumer) Cancel(consumer string, noWait bool) {
	if err := c.ch.Cancel(consumer, noWait); err != nil {
		zap.Log.Error("Failed to cancel channel: ", err)
	}
}

func (c *consumer) Close() {
	if err := c.ch.Close(); err != nil {
		zap.Log.Error("Failed to close channel: ", err)
	}
}
