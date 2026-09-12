package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration.
type Config struct {
	Port         string
	RabbitMQURL  string
	ExchangeName string
	QueueName    string
	DLXName      string
}

// LoadConfig reads configuration from environment variables with defaults.
func LoadConfig() *Config {
	port := getEnv("PORT", "8080")
	user := getEnv("RABBITMQ_USER", "guest")
	pass := getEnv("RABBITMQ_PASS", "guest")
	host := getEnv("RABBITMQ_HOST", "rabbitmq")
	rabbitPort := getEnv("RABBITMQ_PORT", "5672")

	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, rabbitPort)

	return &Config{
		Port:         port,
		RabbitMQURL:  rabbitURL,
		ExchangeName: getEnv("EXCHANGE_NAME", "sistema.eventos.bus"),
		QueueName:    getEnv("QUEUE_NAME", "api-go.procesamiento_inventario"),
		DLXName:      getEnv("DLX_NAME", "sistema.dlx"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
