package kafka

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/IBM/sarama"
	"github.com/bigdann09/notifications/internal/config"
)

type IKafkaProducer interface {
	Close() error
}

type KafkaProducer struct {
	producer sarama.SyncProducer
	consumer any
}

func NewKafkaProducer(cfg *config.KafkaConfig) (*KafkaProducer, error) {
	if cfg.Brokers == "" {
		return nil, errors.New("kafka brokers not configured")
	} else if (cfg.Retries) == "" {
		return nil, errors.New("kafka retries not configured")
	}

	config := sarama.NewConfig()
	retries, _ := strconv.Atoi(cfg.Retries)
	config.Producer.Retry.Max = retries
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Net.MaxOpenRequests = 1
	config.Producer.Idempotent = true

	brokers := strings.Split(cfg.Brokers, ",")
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{producer: producer}, nil
}

func (k *KafkaProducer) Close() error {
	return k.producer.Close()
}

func (k *KafkaProducer) Publish(topic, key string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	}

	_, _, err = k.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
