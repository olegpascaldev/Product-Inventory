package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// ProductCreatedEvent представляет событие, публикуемое при создании продукта.
type ProductCreatedEvent struct {
	ProductID    string `json:"product_id"`
	InitialStock int32  `json:"initial_stock"`
}

// NewProducer создает новый синхронный producer Kafka.
func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

// Функция PublishProductCreated отправляет событие ProductCreated в Kafka.
func (p *Producer) PublishProductCreated(ctx context.Context, productID uuid.UUID, initialStock int32) error {
	event := ProductCreatedEvent{
		ProductID:    productID.String(),
		InitialStock: initialStock,
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(productID.String()),
		Value: sarama.ByteEncoder(data),
	}
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}
	log.Printf("Message sent to partition %d at offset %d", partition, offset)
	return nil
}

// Функция Close завершает работу producer.

func (p *Producer) Close() error {
	return p.producer.Close()
}
