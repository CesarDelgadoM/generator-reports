package databus

import (
	"sync"

	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/internal/consumer"
	"github.com/CesarDelgadoM/generator-reports/internal/workerpool"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
	"github.com/CesarDelgadoM/generator-reports/pkg/stream"
)

const (
	queues_names_queue = "queues-names-queue"

	reportTypeBranch = "branch"
)

type IDataBus interface {
	StartDataBusConsumer()
}

type dataBus struct {
	queuenames consumer.IConsumer
	workerpool *workerpool.WorkerPool
	wg         sync.WaitGroup
	branches   consumer.IQueueConsumer
}

func NewDataBusConsumer(config *config.Config, rabbitmq *stream.RabbitMQ, workerpool *workerpool.WorkerPool, branches consumer.IQueueConsumer) IDataBus {

	opts := &consumer.ConsumerOpts{
		ExchangeType: config.DataBus.Consumer.ExchangeType,
		ContentType:  config.DataBus.Consumer.ContentType,
	}

	return &dataBus{
		queuenames: consumer.NewConsumer(opts, rabbitmq),
		branches:   branches,
		workerpool: workerpool,
	}
}

func (db *dataBus) StartDataBusConsumer() {

	queue := db.queuenames.Queue(&stream.QueueOpts{
		Name:    queues_names_queue,
		Durable: true,
	})

	msgs := db.queuenames.Consume(&stream.ConsumeOpts{
		Name:    queue.Name,
		AutoAck: true,
	})

	db.wg.Add(1)
	go func() {
		defer db.wg.Done()

		var task workerpool.Task

		for m := range msgs {
			msg := consumer.UnmarshalMessageQueueNames(m.Body)
			if msg == nil {
				continue
			}

			switch msg.ReportType {

			case reportTypeBranch:
				task = func() {
					db.branches.ConsumeQueueName(msg.QueueName)
				}

			default:
				zap.Log.Info("Report type not implemented: ", msg.ReportType)
			}

			// Submit task to the workerpool
			db.workerpool.Submit(task)
		}
	}()
	db.wg.Wait()
}
