package branch

import (
	"time"

	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/internal/consumer"
	"github.com/CesarDelgadoM/generator-reports/internal/generators"
	"github.com/CesarDelgadoM/generator-reports/internal/repositorys"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
	"github.com/CesarDelgadoM/generator-reports/pkg/stream"
)

const (
	idle_timeout = 30 * time.Second
	suffix       = "-consumer"
)

type BranchConsumer struct {
	config   *config.Config
	email    IEmailBranch
	consumer consumer.IConsumer
	repo     repositorys.IUserRepository
}

func NewBranchConsumer(config *config.Config, rabbitmq *stream.RabbitMQ, email IEmailBranch, repo repositorys.IUserRepository) consumer.IQueueConsumer {
	opts := &consumer.ConsumerOpts{
		ExchangeType: config.Branch.Consumer.ExchangeType,
		ContentType:  config.Branch.Consumer.ContentType,
	}

	return &BranchConsumer{
		config:   config,
		email:    email,
		consumer: consumer.NewConsumer(opts, rabbitmq),
		repo:     repo,
	}
}

func (bc *BranchConsumer) ConsumeQueueName(queuename string) {
	zap.Log.Info(queuename, " Consumer branch queue start")

	timeout := time.NewTimer(idle_timeout)
	consumerName := queuename + suffix

	// Disconnect consumer
	defer bc.consumer.QueueDelete(&stream.QueueDelete{Name: queuename})
	defer bc.consumer.Cancel(consumerName, false)

	queue := bc.consumer.Queue(&stream.QueueOpts{
		Name: queuename,
	})

	msgs := bc.consumer.Consume(&stream.ConsumeOpts{
		Name:     queue.Name,
		Consumer: consumerName,
	})

	var generator generators.IReport

	// Restaurant data
	m := <-msgs
	msg := consumer.UnmarshalMessage(m.Body)

	generator, err := GetGeneratorBranchReport(msg.Format, bc.config)
	if err != nil {
		zap.Log.Info(queuename, " Error: ", err)
		return
	}

	if err = generator.GenerateReport(msg); err != nil {
		zap.Log.Info(queuename, " Error: ", err)
		return
	}

	// Branches
Loop:
	for {
		select {
		case m := <-msgs:
			msg := consumer.UnmarshalMessage(m.Body)

			if err = generator.GenerateReport(msg); err != nil {
				zap.Log.Info(queuename, " Error: ", err)
				return
			}

			m.Ack(false)

			if msg.Status == 0 {
				zap.Log.Info(queuename, " Status indicator value is: ", msg.Status)

				// Close report
				file, err := generator.CloseReport()
				if err != nil {
					zap.Log.Info(queuename, " Error: ", err)
					return
				}

				// Get user email by user id
				email, err := bc.repo.GetEmailById(msg.UserId)
				if err != nil {
					zap.Log.Info(queuename, " Error: ", err)
					return
				}

				// Send email
				bc.email.SendEmail(queuename, file, email)

				break Loop
			}

			timeout.Reset(idle_timeout)

		case <-timeout.C:
			zap.Log.Info(queuename, " Timeout exceeded, disconnecting consumer")
			break Loop
		}
	}

	zap.Log.Info(queuename, " Consumer branch queue finished")
}
