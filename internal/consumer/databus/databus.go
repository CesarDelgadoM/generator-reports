package databus

import (
	"fmt"

	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/internal/consumer"
	"github.com/CesarDelgadoM/generator-reports/internal/generators/branch"
	"github.com/CesarDelgadoM/generator-reports/internal/workerpool"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
	"github.com/CesarDelgadoM/generator-reports/pkg/stream"
)

const (
	queuenames = "queues-names-queue"

	reportTypeBranch = "branch"
)

type IDataBus interface {
	ConsumeQueueNames(config *config.Consumer)
}

type dataBus struct {
	consumer   consumer.IConsumer
	rabbitmq   *stream.RabbitMQ
	workerpool *workerpool.WorkerPool
}

func NewDataBus(config *config.DataBus, rabbitmq *stream.RabbitMQ, workerpool *workerpool.WorkerPool) IDataBus {

	opts := &consumer.ConsumerOpts{
		ExchangeType: config.ExchangeType,
		ContentType:  config.ContentType,
	}

	return &dataBus{
		consumer:   consumer.NewConsumer(opts, rabbitmq),
		rabbitmq:   rabbitmq,
		workerpool: workerpool,
	}
}

func (db *dataBus) ConsumeQueueNames(config *config.Consumer) {

	queue := db.consumer.Queue(&stream.QueueOpts{
		Name:    queuenames,
		Durable: true,
	})

	msgs := db.consumer.Consume(&stream.ConsumeOpts{
		Name:    queue.Name,
		AutoAck: true,
	})

	go func() {
		var task workerpool.Task

		for m := range msgs {
			msg := consumer.UnmarshalMessageQueueNames(m.Body)
			if msg == nil {
				continue
			}

			switch msg.ReportType {

			case reportTypeBranch:
				task = func() {
					consumer := branch.NewBranchConsumer(config.Branch, db.rabbitmq)
					consumer.ConsumeBranchQueue(msg.QueueName)
				}

			default:
				zap.Log.Info("Report type not implemented: ", msg.ReportType)
			}

			// Submit task to the workerpool
			db.workerpool.Submit(task)
		}
		fmt.Println("finished...")
	}()
}
