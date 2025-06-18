package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/domain"
)

// Consumer implements the domain.EventConsumer interface using Kafka
type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	logger        *zap.Logger
	config        *Config
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(config *Config, consumerGroupID string, logger *zap.Logger) (*Consumer, error) {
	// Configure Sarama
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaConfig.Consumer.Group.Session.Timeout = 10000
	saramaConfig.Consumer.Group.Heartbeat.Interval = 3000
	saramaConfig.Consumer.Return.Errors = true

	// Create consumer group
	consumerGroup, err := sarama.NewConsumerGroup(config.Brokers, consumerGroupID, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		consumerGroup: consumerGroup,
		logger:        logger.Named("kafka_consumer"),
		config:        config,
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

// ConsumeOrderEvents starts consuming order events
func (c *Consumer) ConsumeOrderEvents(handler domain.OrderEventHandler) error {
	return c.consumeEvents([]string{c.config.Topics.OrderEvents}, handler)
}

// ConsumeInventoryEvents starts consuming inventory events
func (c *Consumer) ConsumeInventoryEvents(handler domain.OrderEventHandler) error {
	return c.consumeEvents([]string{c.config.Topics.InventoryEvents}, handler)
}

// ConsumePaymentEvents starts consuming payment events
func (c *Consumer) ConsumePaymentEvents(handler domain.OrderEventHandler) error {
	return c.consumeEvents([]string{c.config.Topics.PaymentEvents}, handler)
}

// consumeEvents starts consuming events from the specified topics
func (c *Consumer) consumeEvents(topics []string, handler domain.OrderEventHandler) error {
	consumerHandler := &eventConsumerHandler{
		handler: handler,
		logger:  c.logger,
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			// Check if context was cancelled
			if c.ctx.Err() != nil {
				return
			}

			err := c.consumerGroup.Consume(c.ctx, topics, consumerHandler)
			if err != nil {
				c.logger.Error("Error consuming events", 
					zap.Error(err),
					zap.Strings("topics", topics),
				)
			}
		}
	}()

	// Handle consumer errors
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for err := range c.consumerGroup.Errors() {
			c.logger.Error("Consumer group error", zap.Error(err))
		}
	}()

	return nil
}

// Close closes the Kafka consumer
func (c *Consumer) Close() error {
	c.logger.Info("Closing Kafka consumer")
	c.cancel()
	c.wg.Wait()
	
	if c.consumerGroup != nil {
		return c.consumerGroup.Close()
	}
	return nil
}

// eventConsumerHandler implements sarama.ConsumerGroupHandler
type eventConsumerHandler struct {
	handler domain.OrderEventHandler
	logger  *zap.Logger
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *eventConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	h.logger.Info("Consumer group session started")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *eventConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	h.logger.Info("Consumer group session ended")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (h *eventConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Handle messages in a loop
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			h.logger.Debug("Received message", 
				zap.String("topic", message.Topic),
				zap.Int32("partition", message.Partition),
				zap.Int64("offset", message.Offset),
			)

			// Process the message
			err := h.processMessage(message)
			if err != nil {
				h.logger.Error("Failed to process message", 
					zap.Error(err),
					zap.String("topic", message.Topic),
					zap.Int32("partition", message.Partition),
					zap.Int64("offset", message.Offset),
				)
				// Continue processing other messages even if one fails
			}

			// Mark message as processed
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

// processMessage processes a Kafka message and delegates to the event handler
func (h *eventConsumerHandler) processMessage(message *sarama.ConsumerMessage) error {
	// Deserialize the event
	event, err := domain.FromJSON(message.Value)
	if err != nil {
		return fmt.Errorf("failed to deserialize event: %w", err)
	}

	// Add message metadata
	event.AddMetadata("kafka_topic", message.Topic)
	event.AddMetadata("kafka_partition", fmt.Sprintf("%d", message.Partition))
	event.AddMetadata("kafka_offset", fmt.Sprintf("%d", message.Offset))

	// Handle the event
	err = h.handler.HandleEvent(event)
	if err != nil {
		return fmt.Errorf("failed to handle event: %w", err)
	}

	h.logger.Debug("Successfully processed event", 
		zap.String("event_id", event.ID),
		zap.String("event_type", string(event.Type)),
		zap.String("order_id", event.OrderID),
	)

	return nil
}
