package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

var (
	ErrCtxCancelled = fmt.Errorf("Context cancelled")
)

type (
	LiveStatusQueue interface {
		StartQueue(jobID string) (*LiveStatusQueueSession, error)
		WaitAndRecieveProgressMsgsQueue(ctx context.Context, jobID string) (<-chan amqp.Delivery, error)
	}

	RabbitMqLiveStatusQueue struct {
		Ch *amqp.Channel
	}
)

func NewRabbitMqLiveStatusQueue(rabbitConfig RabbitMQConfig) (LiveStatusQueue, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		rabbitConfig.User, rabbitConfig.Password, rabbitConfig.Host, rabbitConfig.Port))
	if err != nil {
		return nil, fmt.Errorf("RabbitMQ Connection Error: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("RabbitMQ Channel Error: %v", err)
	}

	return &RabbitMqLiveStatusQueue{Ch: ch}, nil
}

func (rpm *RabbitMqLiveStatusQueue) StartQueue(jobID string) (*LiveStatusQueueSession, error) {
	queueName := fmt.Sprintf("result_queue_%s", jobID)
	// log.Debug("Creating queue:", queueName)
	_, err := rpm.Ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("RabbitMQ Queue Declare Error: %v", err)
	}

	return &LiveStatusQueueSession{
		ch:        rpm.Ch,
		queueName: queueName,
	}, nil
}

func (rpm *RabbitMqLiveStatusQueue) WaitAndRecieveProgressMsgsQueue(ctx context.Context, jobID string) (<-chan amqp.Delivery, error) {

	queueName := fmt.Sprintf("result_queue_%s", jobID)

	// Poll for queue existence with a timeoutr
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done(): // Context cancelled (e.g., WebSocket closed)
			fmt.Println("Context cancelled: stopping queue wait for ", queueName)

			return nil, ErrCtxCancelled

		case <-ticker.C: // Periodically check for queue existence
			_, err := rpm.Ch.QueueInspect(queueName)
			if err == nil {
				// Queue exists, start consuming
				msgs, err := rpm.Ch.Consume(queueName, "", true, false, false, false, nil)
				if err != nil {
					return nil, fmt.Errorf("RabbitMQ Consume Error: %v", err)
				}

				return msgs, nil
			}
		}
	}
}
