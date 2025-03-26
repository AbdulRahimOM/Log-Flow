package queue

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

type (
	LiveStatusQueueSession struct {
		ch        *amqp.Channel
		queueName string
	}
)

func (q *LiveStatusQueueSession) SendIntermediateResult(result string) {
	err := q.ch.Publish("", q.queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(result),
	})
	if err != nil {
		fmt.Println("RabbitMQ Publish Error:", err)
	}

	log.Debug("✅ Sent intermediate result to RabbitMQ:", result)
}

func (q *LiveStatusQueueSession) Delete() {
	time.Sleep(3 * time.Second) // Wait for 3 seconds before deleting the queue, so that the client can consume all messages
	_, err := q.ch.QueueDelete(q.queueName, false, false, false)
	if err != nil {
		log.Error("❌ Failed to delete queue: %v", err)
		return
	}
}
