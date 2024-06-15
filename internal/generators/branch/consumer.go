package branch

import (
	"time"

	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/internal/consumer"
	"github.com/CesarDelgadoM/generator-reports/internal/generators"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
	"github.com/CesarDelgadoM/generator-reports/pkg/stream"
)

const (
	idle_timeout = 30 * time.Second
	suffix       = "-consumer"
)

type BranchConsumer struct {
	config   *config.Config
	consumer consumer.IConsumer
}

func NewBranchConsumer(config *config.Config, rabbitmq *stream.RabbitMQ) consumer.IQueueConsumer {
	opts := &consumer.ConsumerOpts{
		ExchangeType: config.Branch.Consumer.ExchangeType,
		ContentType:  config.Branch.Consumer.ContentType,
	}

	return &BranchConsumer{
		config:   config,
		consumer: consumer.NewConsumer(opts, rabbitmq),
	}
}

func (bc *BranchConsumer) ConsumeQueueName(queueName string) {
	zap.Log.Info("Consumer branch queue start: ", queueName)

	timeout := time.NewTimer(idle_timeout)
	consumerName := queueName + suffix

	// Disconnect consumer
	defer bc.consumer.QueueDelete(&stream.QueueDelete{Name: queueName})
	defer bc.consumer.Cancel(consumerName, false)

	queue := bc.consumer.Queue(&stream.QueueOpts{
		Name: queueName,
	})

	msgs := bc.consumer.Consume(&stream.ConsumeOpts{
		Name:     queue.Name,
		Consumer: consumerName,
		AutoAck:  true,
	})

	var generator generators.IReport

	m := <-msgs
	msg := consumer.UnmarshalMessage(m.Body)

	// Initialize report with restaurant data
	generator, err := GetGeneratorBranchReport(msg.Format, bc.config)
	if err != nil {
		zap.Log.Error(err)
		return
	}
	generator.GenerateReport(msg)

	// Branches
Loop:
	for {
		select {
		case m := <-msgs:
			msg := consumer.UnmarshalMessage(m.Body)

			generator.GenerateReport(msg)

			timeout.Reset(idle_timeout)

		case <-timeout.C:
			zap.Log.Info("Timeout exceeded, disconnecting consumer")
			break Loop
		}
	}

	zap.Log.Info("Consumer branch queue finished: ", queueName)
}
