package queue

import (
	"log-flow/internal/infrastructure/config"

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
		SentForRetry(msg amqp.Delivery) error 
	}

	LogQueue interface {
		LogQueueReceiver
		LogQueueSender
	}

	rabbitMqLogFileQueue struct {
		conn     *amqp.Connection
		ch       *amqp.Channel
	}

	LogMessage struct {
		JobID    string `json:"job_id"`
		FileURL  string `json:"file_url"`
		Priority uint8  `json:"priority"`
	}
)

func InitLogQueue() LogQueue {
	logFileQueue, err := NewRabbitMQLogQueue(getRabbitMQConfig())
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	return logFileQueue
}

func InitLiveStatusQueue() LiveStatusQueue {
	liveStatusQueue, err := NewRabbitMqLiveStatusQueue(getRabbitMQConfig())
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ live progress channel: %v", err)
	}

	return liveStatusQueue
}

func getRabbitMQConfig() RabbitMQConfig {
	return RabbitMQConfig{
		Host:     config.Env.RabbitMQConfig.Host,
		Port:     config.Env.RabbitMQConfig.Port,
		User:     config.Env.RabbitMQConfig.User,
		Password: config.Env.RabbitMQConfig.Password,
	}
}
