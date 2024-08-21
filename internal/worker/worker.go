package worker

import (
	"context"
	"github.com/IBM/sarama"
	"log"
	"message/internal/kafka"
	"message/storage"
	"os"
	"strconv"
	"time"
)

func New(ctx context.Context, st storage.Storage, producer *sarama.SyncProducer) {
	timeout, err := strconv.Atoi(os.Getenv("SERVER_WORKER_TIMEOUT"))
	if err != nil {
		log.Fatalf("Failed to convert SERVER_WORKER_TIMEOUT to int: %v", err)
	}
	topic := os.Getenv("SERVER_TOPIC")
	log.Printf("Worker is running with timeout %d seconds", timeout)

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker is done")
			return
		default:
			log.Println("Worker is waiting")
			time.Sleep(time.Duration(timeout) * time.Second)
			log.Println("Worker started")
			messages := st.NeedSent()
			for _, message := range messages {
				select {
				case <-ctx.Done():
					log.Println("Worker is done")
					return
				default:
					err := kafka.SendMessage(message, producer, topic)
					if err != nil {
						log.Printf("Failed to send message: %v", err)
						continue
					}
					log.Println("Message sent to Kafka")
					err = st.ChangeStatusSent(ctx, message.ID)
					if err != nil {
						log.Printf("Failed to change sent status: %v", err)
						continue
					}
					log.Println("Status changed for message with ID: ", message.ID)
				}
			}
		}
	}
}
