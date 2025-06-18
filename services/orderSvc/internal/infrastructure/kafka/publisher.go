package kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/domain"
)

// Publisher implements the domain.EventPublisher interface using Kafka
type Publisher struct {
	producer sarama.SyncProducer
	logger   *zap.Logger
	config   *Config
}

// Config holds Kafka configuration
type Config struct {
	Brokers []string
	Topics  TopicConfig
}

// TopicConfig defines Kafka topic names
type TopicConfig struct {
	OrderEvents     string
	InventoryEvents string
	PaymentEvents   string
}

// NewPublisher creates a new Kafka publisher
func NewPublisher(config *Config, logger *zap.Logger) (*Publisher, error) {
	// Configure Sarama
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Producer.Retry.Max = 3
	saramaConfig.Producer.Retry.Backoff = 100 * time.Millisecond
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas
	saramaConfig.Producer.Compression = sarama.CompressionSnappy
	saramaConfig.Producer.Flush.Frequency = 500 * time.Millisecond

	// Create producer
	producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &Publisher{
		producer: producer,
		logger:   logger.Named("kafka_publisher"),
		config:   config,
	}, nil
}

// PublishOrderEvent publishes an order lifecycle event
func (p *Publisher) PublishOrderEvent(event *domain.OrderEvent) error {
	return p.publishEvent(event, p.config.Topics.OrderEvents)
}

// PublishInventoryEvent publishes an inventory-related event
func (p *Publisher) PublishInventoryEvent(event *domain.OrderEvent) error {
	return p.publishEvent(event, p.config.Topics.InventoryEvents)
}

// PublishPaymentEvent publishes a payment-related event
func (p *Publisher) PublishPaymentEvent(event *domain.OrderEvent) error {
	return p.publishEvent(event, p.config.Topics.PaymentEvents)
}

// publishEvent publishes an event to the specified topic
func (p *Publisher) publishEvent(event *domain.OrderEvent, topic string) error {
	// Add metadata
	event.AddMetadata("service", "orderSvc")
	event.AddMetadata("environment", "production") // This should come from config
	
	// Serialize event to JSON
	eventData, err := event.ToJSON()
	if err != nil {
		p.logger.Error("Failed to serialize event", 
			zap.Error(err),
			zap.String("event_id", event.ID),
			zap.String("event_type", string(event.Type)),
		)
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	// Create Kafka message
	message := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(event.OrderID), // Use order ID as partition key
		Value:     sarama.ByteEncoder(eventData),
		Timestamp: event.Timestamp,
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("event_type"),
				Value: []byte(event.Type),
			},
			{
				Key:   []byte("event_id"),
				Value: []byte(event.ID),
			},
			{
				Key:   []byte("order_id"),
				Value: []byte(event.OrderID),
			},
			{
				Key:   []byte("user_id"),
				Value: []byte(event.UserID),
			},
		},
	}

	// Publish message
	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		p.logger.Error("Failed to publish event", 
			zap.Error(err),
			zap.String("topic", topic),
			zap.String("event_id", event.ID),
			zap.String("event_type", string(event.Type)),
			zap.String("order_id", event.OrderID),
		)
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.logger.Info("Successfully published event", 
		zap.String("topic", topic),
		zap.String("event_id", event.ID),
		zap.String("event_type", string(event.Type)),
		zap.String("order_id", event.OrderID),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
	)

	return nil
}

// Close closes the Kafka producer
func (p *Publisher) Close() error {
	if p.producer != nil {
		return p.producer.Close()
	}
	return nil
}
