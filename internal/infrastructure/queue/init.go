package queue

import (
	"log-flow/internal/infrastructure/config"

	"github.com/gofiber/fiber/v2/log"
)

const (
	logFilesExchange   = "log_files_exchange"
	logProcessingQueue = "log_processing_queue"
)

func InitLogQueue() LogQueue {
	logFileQueue, err := NewRabbitMQLogQueue(getRabbitMQConfig(), logFilesExchange, logProcessingQueue)
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
