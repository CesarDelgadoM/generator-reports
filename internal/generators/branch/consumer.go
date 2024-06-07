package branch

import (
	"sync"
	"time"

	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/internal/consumer"
	"github.com/CesarDelgadoM/generator-reports/internal/generators/restaurant"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
	"github.com/CesarDelgadoM/generator-reports/pkg/stream"
)

const (
	IDLE_TIMEOUT = 60 * time.Second
	SUFFIX       = "-consumer"
)

type BranchConsumer struct {
	wg       sync.WaitGroup
	consumer consumer.IConsumer
}

func NewBranchConsumer(config *config.Branch, rabbitmq *stream.RabbitMQ) *BranchConsumer {
	opts := &consumer.ConsumerOpts{
		ExchangeType: config.ExchangeType,
		ContentType:  config.ContentType,
	}

	return &BranchConsumer{
		consumer: consumer.NewConsumer(opts, rabbitmq),
	}
}

func (bc *BranchConsumer) ConsumeBranchQueue(queueName string) {
	zap.Log.Info("Consumer branch queue start: ", queueName)

	timeout := time.NewTimer(IDLE_TIMEOUT)
	consumerName := queueName + SUFFIX

	queue := bc.consumer.Queue(&stream.QueueOpts{
		Name: queueName,
	})

	msgs := bc.consumer.Consume(&stream.ConsumeOpts{
		Name:     queue.Name,
		Consumer: consumerName,
		AutoAck:  true,
	})

	bc.wg.Add(1)
	go func() {
		defer bc.wg.Done()
		// Disconnect consumer
		defer bc.consumer.Close()
		defer bc.consumer.QueueDelete(&stream.QueueDelete{Name: queueName})
		defer bc.consumer.Cancel(consumerName, false)

		var format string

		// Initialize report wiht restaurant data
		select {
		case m := <-msgs:
			timeout.Reset(IDLE_TIMEOUT)

			msg := consumer.UnmarshalMessage(m.Body)

			format = msg.Format

			restaurantGenerator, err := restaurant.GetGeneratorRestaurantReport(format)
			if err != nil {
				return
			}
			restaurantGenerator.GenerateReport(msg)

		case <-timeout.C:
			zap.Log.Info("Timeout exceeded, disconnecting consumer")
			return
		}

		// Branches
		branchGenerator, err := GetGeneratorBranchReport(format)
	Loop:
		for {
			select {
			case m := <-msgs:
				timeout.Reset(IDLE_TIMEOUT)

				msg := consumer.UnmarshalMessage(m.Body)

				branchGenerator.GenerateReport(msg)
				if err != nil {
					return
				}

			case <-timeout.C:
				zap.Log.Info("Timeout exceeded, disconnecting consumer")
				break Loop
			}
		}
	}()
	bc.wg.Wait()
	zap.Log.Info("Consumer branch queue finished: ", queueName)
}
