package branch

import (
	"sync"
	"time"

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

func NewBranchConsumer(consumer consumer.IConsumer) *BranchConsumer {
	return &BranchConsumer{
		consumer: consumer,
	}
}

func (bc *BranchConsumer) ConsumeBranchQueue(queueName string) {
	var timeout *time.Timer = time.NewTimer(idleTimeout)
	var consumerName string = queueName + suffix

	queue := bc.consumer.Queue(&stream.QueueOpts{
		Name:    queueName,
		Durable: true,
	})

	msgs := bc.consumer.Consume(&stream.ConsumeOpts{
		Name:      queue.Name,
		Consumer:  consumerName,
		AutoAck:   true,
		Exclusive: true,
	})

	bc.wg.Add(1)
	go func() {
		defer bc.consumer.Cancel(consumerName, false)
		defer bc.wg.Done()

		zap.Log.Info("Consume branch queue...")

		// Restaurant Data
		select {
		case m := <-msgs:
			msg := utils.UnmarshalMessage(m.Body)
			restaurant := utils.UnmarshalRestaurant(msg.Data)

			zap.Log.Info(restaurant)

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
				if msg.Status == 0 {
					zap.Log.Info("Finished process successfully")
					return
				}
				//branch := utils.UnmarshalBranches(msg.Data)

				// Branches Data
				//zap.Log.Info(branch)

				timeout.Reset(idleTimeout)

			case <-timeout.C:
				zap.Log.Info("Timeout exceeded, the process not finished successfully")
				return
			}
		}

	}()
	bc.wg.Wait()
	zap.Log.Info("Consume branch queue finished")
}
