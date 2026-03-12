package kafka

import (
	"context"
	"encoding/json"
	"log"

	"inventory-service/internal/usecase"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

// Consumer represents a Kafka consumer group handler.
type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	usecase       *usecase.StockUsecase
	topic         string
}

// ProductCreatedEvent matches the event published by Product Service.
type ProductCreatedEvent struct {
	ProductID    string `json:"product_id"`
	InitialStock int32  `json:"initial_stock"`
}

// NewConsumer creates a new Kafka consumer group.
func NewConsumer(brokers []string, groupID, topic string, uc *usecase.StockUsecase) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumerGroup: consumerGroup,
		usecase:       uc,
		topic:         topic,
	}, nil
}

// Start begins consuming messages. It blocks until an error or context cancel.
func (c *Consumer) Start(ctx context.Context) error {
	for {
		err := c.consumerGroup.Consume(ctx, []string{c.topic}, c)
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// Close shuts down the consumer.
func (c *Consumer) Close() error {
	return c.consumerGroup.Close()
}

// Setup is called once when a new session begins.
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is called at the end of a session.
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from a claim.
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: topic=%s, partition=%d, offset=%d", msg.Topic, msg.Partition, msg.Offset)

		var event ProductCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			session.MarkMessage(msg, "")
			continue
		}

		productID, err := uuid.Parse(event.ProductID)
		if err != nil {
			log.Printf("Invalid product ID in event: %v", err)
			session.MarkMessage(msg, "")
			continue
		}

		// Create stock entry
		if err := c.usecase.CreateStock(session.Context(), productID, event.InitialStock); err != nil {
			log.Printf("Failed to create stock for product %s: %v", productID, err)
			// In production, implement retry or dead letter queue
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
