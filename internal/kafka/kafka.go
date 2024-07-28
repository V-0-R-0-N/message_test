package kafka

import (
	"github.com/IBM/sarama"
	"message/models"
	"net"
	"time"

	"log"
)

const (
	KAFKAHOST = "kafka"
	KAFKAPORT = "9092"
	TIMEOUT   = 10 // seconds
)

// InitProducer инициализирует продюсера Kafka
func InitProducer() *sarama.SyncProducer {
	// Конфигурация клиента Kafka
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	// Создание хоста
	hostPort := net.JoinHostPort(KAFKAHOST, KAFKAPORT)
	time.Sleep(TIMEOUT * time.Second)
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
