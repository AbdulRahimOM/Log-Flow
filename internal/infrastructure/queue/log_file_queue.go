package queue

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

type (
	RabbitMQConfig struct {
		Host     string
		Port     string
		User     string
		Password string
	}

	LogQueueSender interface {
		SendToQueue(logMsg LogMessage) error
		GetQueueStatus() (map[string]any, error)
	}

	LogQueueReceiver interface {
		RecieveLogFileDetails() (<-chan amqp.Delivery, error)
	}

	LogQueue interface {
		LogQueueReceiver
		LogQueueSender
	}

	rabbitMqLogFileQueue struct {
		conn     *amqp.Connection
		ch       *amqp.Channel
		queue    string
		exchange string
	}

	LogMessage struct {
		JobID    string `json:"job_id"`
		FileURL  string `json:"file_url"`
		Priority uint8  `json:"priority"`
	}
)

func NewRabbitMQLogQueue(rabbitConfig RabbitMQConfig, exchange, queue string) (*rabbitMqLogFileQueue, error) {

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		rabbitConfig.User, rabbitConfig.Password, rabbitConfig.Host, rabbitConfig.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	err = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	_, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-max-priority": 10,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	err = ch.QueueBind(
		queue,
		queue, // Routing key same as queue name for direct exchange
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}

	return &rabbitMqLogFileQueue{
		conn:  conn,
		ch:    ch,
		queue: queue,
	}, nil
}

func (rq *rabbitMqLogFileQueue) SendToQueue(logMsg LogMessage) error {
	msgBody, err := json.Marshal(logMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	err = rq.ch.Publish(
		rq.exchange,
		rq.queue,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgBody,
			Priority:    logMsg.Priority,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	log.Trace("✅ Sent message to RabbitMQ: %s\n", msgBody)
	return nil
}

func (rq *rabbitMqLogFileQueue) RecieveLogFileDetails() (<-chan amqp.Delivery, error) {
	return rq.ch.Consume(rq.queue, "", true, false, false, false, nil)
}

func (rq *rabbitMqLogFileQueue) GetQueueStatus() (map[string]any, error) {
	queueInfo, err := rq.ch.QueueInspect(rq.queue)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect queue: %v", err)
	}

	return map[string]any{
		"name":           queueInfo.Name, //queue name
		"message_count":  queueInfo.Messages, //no. of msgs
		"consumer_count": queueInfo.Consumers, //no. of consumers
	}, nil
}
