package kafka

import (
	"context"
	"errors"
	"strings"

	"github.com/IBM/sarama"
	"github.com/bigdann09/notifications/internal/config"
)

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
}

func NewKafkaConsumer(cfg *config.KafkaConfig, groupID string) (*KafkaConsumer, error) {
	if cfg.Brokers == "" {
		return nil, errors.New("kafka brokers not configured")
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	brokers := strings.Split(cfg.Brokers, ",")
	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{consumer: consumer}, nil
}

func (k *KafkaConsumer) Close() error {
	return k.consumer.Close()
}

func (k *KafkaConsumer) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	for {
		if err := k.consumer.Consume(ctx, topics, handler); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
