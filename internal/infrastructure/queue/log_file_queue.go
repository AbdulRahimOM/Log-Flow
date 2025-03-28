package queue

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

const (
	logFilesExchange   = "log_files_exchange"
	logProcessingQueue = "log_processing_queue"
)

const (
	// dlxExchange     = "dlx_exchange"
	dlxQueue        = "dlx_log_processing_queue"
	retryRoutingKey = "retry"
	dlxTTL          = 10000 // TTL in milliseconds (10 seconds)
)

func NewRabbitMQLogQueue(rabbitConfig RabbitMQConfig) (*rabbitMqLogFileQueue, error) {

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		rabbitConfig.User, rabbitConfig.Password, rabbitConfig.Host, rabbitConfig.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	err = ch.ExchangeDeclare(logFilesExchange, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	_, err = ch.QueueDeclare(
		logProcessingQueue,
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
		logProcessingQueue,
		logProcessingQueue, // Routing key same as queue name for direct exchange
		logFilesExchange,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}

	// Declare DLX Queue
	_, err = ch.QueueDeclare(
		dlxQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl":             dlxTTL,
			"x-dead-letter-exchange":    logFilesExchange,
			"x-dead-letter-routing-key": logProcessingQueue,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare DLX queue: %v", err)
	}

	// Bind DLX Queue to Log Files Exchange (not DLX Exchange)
	err = ch.QueueBind(dlxQueue, retryRoutingKey, logFilesExchange, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to bind DLX queue: %v", err)
	}

	return &rabbitMqLogFileQueue{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rq *rabbitMqLogFileQueue) SendToQueue(logMsg LogMessage) error {
	msgBody, err := json.Marshal(logMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	err = rq.ch.Publish(
		logFilesExchange,
		logProcessingQueue,
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
	return rq.ch.Consume(logProcessingQueue, "", true, false, false, false, nil)
}

func (rq *rabbitMqLogFileQueue) GetQueueStatus() (map[string]any, error) {
	queueInfo, err := rq.ch.QueueInspect(logProcessingQueue)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect queue: %v", err)
	}

	return map[string]any{
		"name":           queueInfo.Name,      //queue name
		"message_count":  queueInfo.Messages,  //no. of msgs
		"consumer_count": queueInfo.Consumers, //no. of consumers
	}, nil
}

func (rq *rabbitMqLogFileQueue) SentForRetry(msg amqp.Delivery) error {
	fmt.Println("🔄 Sending message to DLX for retry")
	retryCount := 0

	if msg.Headers != nil {
		if val, ok := msg.Headers["x-retry-count"].(int32); ok {
			retryCount = int(val)
		}
	} else {
		msg.Headers = amqp.Table{}
	}

	if retryCount >= 3 {
		log.Error("❌ Retry count exceeded for message")
		return nil
	}

	msg.Headers["x-retry-count"] = retryCount + 1

	err := rq.ch.Publish(
		logFilesExchange,
		retryRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         msg.Body,
			Headers:      msg.Headers,
			DeliveryMode: amqp.Persistent,
		})
	if err != nil {
		return fmt.Errorf("failed to send message to DLX: %v", err)
	}

	log.Debug("🔄 Message sent to DLX for retry")
	return nil
}
