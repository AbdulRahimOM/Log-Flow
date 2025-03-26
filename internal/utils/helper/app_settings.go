package helper

import "github.com/gofiber/fiber/v2/log"

func SetFiberLogLevel(logLevel string) {
	switch logLevel {
	case "debug", "DEBUG":
		log.SetLevel(log.LevelDebug)
	case "info", "INFO":
		log.SetLevel(log.LevelInfo)
	case "error", "ERROR":
		log.SetLevel(log.LevelError)
	default:
		log.SetLevel(log.LevelInfo)
		log.Info("Log level not found, setting to default: INFO")
	}
}
