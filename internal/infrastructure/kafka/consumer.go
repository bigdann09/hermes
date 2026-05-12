package kafka

import (
	"context"
	"strings"

	"github.com/IBM/sarama"
	"github.com/bigdann09/notifications/internal/config"
)

type KafkaConsumer struct {
	group sarama.ConsumerGroup
}

type ConsumerHandler func(ctx context.Context, topic string, key, payload []byte) error

func NewKafkaConsumer(cfg *config.KafkaConfig, groupID string) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	brokers := strings.Split(cfg.Brokers, ",")
	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{group: group}, nil
}

func (c *KafkaConsumer) Close() error {
	return c.group.Close()
}

func (c *KafkaConsumer) Consume(ctx context.Context, topics []string, handler ConsumerHandler) error {
	consumer := &groupHandler{
		handler: handler,
	}

	for {
		if err := c.group.Consume(ctx, topics, consumer); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// groupHandler implements sarama.ConsumerGroupHandler
type groupHandler struct {
	handler ConsumerHandler
}

func (h *groupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *groupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *groupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		err := h.handler(sess.Context(), msg.Topic, msg.Key, msg.Value)
		if err == nil {
			sess.MarkMessage(msg, "")
		}
	}
	return nil
}
