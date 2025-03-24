package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

var (
	ErrCtxCancelled = fmt.Errorf("Context cancelled")
)

type (
	ProgressMessenger interface {
		StartQueue(jobID string) (*progressMsgQueue, error)
		WaitAndRecieveProgressMsgsQueue(ctx context.Context, jobID string) (<-chan amqp.Delivery, error)
	}

	RabbitMqProgressMessenger struct {
		Ch *amqp.Channel
	}

	progressMsgQueue struct {
		ch        *amqp.Channel
		queueName string
	}
)

func NewResultChan(rabbitConfig RabbitMQConfig) (ProgressMessenger, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		rabbitConfig.User, rabbitConfig.Password, rabbitConfig.Host, rabbitConfig.Port))
	if err != nil {
		return nil, fmt.Errorf("RabbitMQ Connection Error: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("RabbitMQ Channel Error: %v", err)
	}

	return &RabbitMqProgressMessenger{Ch: ch}, nil
}

func (rpm *RabbitMqProgressMessenger) StartQueue(jobID string) (*progressMsgQueue, error) {
	queueName := fmt.Sprintf("result_queue_%s", jobID)
	log.Debug("Creating queue:", queueName)
	_, err := rpm.Ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("RabbitMQ Queue Declare Error: %v", err)
	}

	return &progressMsgQueue{
		ch:        rpm.Ch,
		queueName: queueName,
	}, nil
}

func (q *progressMsgQueue) SendIntermediateResult(result string) {
	err := q.ch.Publish("", q.queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(result),
	})
	if err != nil {
		fmt.Println("RabbitMQ Publish Error:", err)
	}

	fmt.Println("✅ Sent intermediate result to RabbitMQ:", result)
}

func (q *progressMsgQueue) Delete() {
	_, err := q.ch.QueueDelete(q.queueName, false, false, false)
	if err != nil {
		fmt.Println("RabbitMQ Queue Delete Error:", err)
		return
	}
}
func (rpm *RabbitMqProgressMessenger) WaitAndRecieveProgressMsgsQueue(ctx context.Context, jobID string) (<-chan amqp.Delivery, error) {

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
