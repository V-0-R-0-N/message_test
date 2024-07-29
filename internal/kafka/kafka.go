package kafka

import (
	"github.com/IBM/sarama"
	"message/models"
	"net"
	"os"
	"strconv"
	"time"

	"log"
)

const (
	TIMEOUT = 10 // seconds
)

// InitProducer инициализирует продюсера Kafka
func InitProducer() *sarama.SyncProducer {
	// Конфигурация клиента Kafka
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	// Получение переменных окружения
	kafkaHost := os.Getenv("SERVER_KAFKA_HOST")
	kafkaPort := os.Getenv("SERVER_KAFKA_PORT")
	timeout, err := strconv.Atoi(os.Getenv("SERVER_KAFKA_TIMEOUT"))
	if err != nil {
		log.Fatalf("Failed to convert WORKER_TIMEOUT to int: %v", err)
	}
	// Создание хоста
	hostPort := net.JoinHostPort(kafkaHost, kafkaPort)
	time.Sleep(time.Duration(timeout) * time.Second)
	// Создание синхронного продюсера
	producer, err := sarama.NewSyncProducer([]string{hostPort}, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
	}
	return &producer
	// TODO defer после вызова функции
}

// SendMessage отправляет сообщение в Kafka
func SendMessage(mes *models.Message, producer *sarama.SyncProducer, topic string) error {
	// Преобразование сообщения в JSON
	body, err := models.MessageToJSONForKafka(mes)
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
		return err
	}
	// Создание сообщения для отправки в Kafka
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(body),
	}
	// Отправка сообщения
	partition, offset, err := (*producer).SendMessage(message)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
		return err
	}
	log.Printf("Message is stored in partition %d, offset %d\n", partition, offset)
	return nil
}
