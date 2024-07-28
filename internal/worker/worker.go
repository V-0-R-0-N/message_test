package worker

import (
	"context"
	"github.com/IBM/sarama"
	"log"
	"message/internal/kafka"
	"message/storage"
	"time"
)

const (
	TOPIC   = "message"
	TIMEOUT = 100 // seconds
)

func New(ctx context.Context, st storage.Storage, producer *sarama.SyncProducer) {
	log.Printf("Worker is running with timeout %d seconds", TIMEOUT)

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker is done")
			return
		default:
			log.Println("Worker is waiting")
			time.Sleep(TIMEOUT * time.Second)
			log.Println("Worker started")
			messages := storage.NeedSent(st)
			for _, message := range messages {
				select {
				case <-ctx.Done():
					log.Println("Worker is done")
					return
				default:
					err := kafka.SendMessage(message, producer, TOPIC)
					if err != nil {
						log.Printf("Failed to send message: %v", err)
						continue
					}
					log.Println("Message sent to Kafka")
					err = storage.ChangeStatusSent(ctx, st, message.ID)
					if err != nil {
						log.Printf("Failed to change sent status: %v", err)
						continue
					}
					log.Println("Status changed")
				}
			}
		}
	}
}
