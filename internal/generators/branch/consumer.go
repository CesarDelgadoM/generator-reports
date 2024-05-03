package branch

import (
	"sync"
	"time"

	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/internal/consumer"
	"github.com/CesarDelgadoM/generator-reports/internal/utils"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
	"github.com/CesarDelgadoM/generator-reports/pkg/stream"
)

const (
	idleTimeout = 60 * time.Second
	suffix      = "-consumer"
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

	timeout := time.NewTimer(idleTimeout)
	consumerName := queueName + suffix

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
		defer bc.consumer.Close()
		defer bc.consumer.QueueDelete(&stream.QueueDelete{
			Name: queueName,
		})
		defer bc.consumer.Cancel(consumerName, false)

		zap.Log.Info("Consume branch queue...")

		// Restaurant Data
		select {
		case m := <-msgs:
			msg := utils.UnmarshalMessage(m.Body)

			restaurant := utils.UnmarshalRestaurant(msg.Data)
			if restaurant != nil {
				zap.Log.Info(restaurant)
			}

			timeout.Reset(idleTimeout)

		case <-timeout.C:
			zap.Log.Info("Timeout exceeded, the process not finished succesfully")
			return
		}

		// Branches
		for {
			select {
			case m := <-msgs:
				msg := utils.UnmarshalMessage(m.Body)

				branches := utils.UnmarshalBranches(msg.Data)
				if branches != nil {
					zap.Log.Info(branches)

					if msg.Status == 0 {
						zap.Log.Info("Finished process successfully")
						return
					}
				}

				timeout.Reset(idleTimeout)

			case <-timeout.C:
				zap.Log.Info("Timeout exceeded, the process not finished successfully")
				return
			}
		}

	}()
	bc.wg.Wait()
	zap.Log.Info("Consumer branch queue finished: ", queueName)
}
