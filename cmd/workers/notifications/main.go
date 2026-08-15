package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/bigdann09/notifications/internal/config"
	"github.com/bigdann09/notifications/internal/infrastructure/database"
	"github.com/bigdann09/notifications/internal/infrastructure/kafka"
	"github.com/bigdann09/notifications/internal/repositories"
	"github.com/bigdann09/notifications/internal/services/notification"
	"github.com/bigdann09/notifications/internal/services/notification/channels"
	"github.com/bigdann09/notifications/pkgs/logger"
	"go.uber.org/zap"
)

const NOTIFICATION_GROUP = "notifications"
const NOTIFICATION_TOPIC = "notifications"

type NotificationConsumer struct {
	logger  *zap.Logger
	service *notification.DispatcherService
}

func (*NotificationConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*NotificationConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (n *NotificationConsumer) ConsumeClaim(
	sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		if NOTIFICATION_TOPIC != msg.Topic {
			n.logger.Warn(
				"received message with wrong topic",
				zap.String("topic", msg.Topic),
				zap.String("key", string(msg.Key)),
			)
			sess.MarkMessage(msg, "")
			continue
		}

		n.service.Dispatch(msg)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(&cfg.App)
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		panic(err)
	}

	notification_repository := repositories.NewNotificationRepository(db)
	channel_list := []channels.IChannel{
		channels.NewDatabaseChannel(notification_repository),
		channels.NewMailChannel(&cfg.Mail),
		channels.NewPushChannel(log),
	}

	srv := notification.NewDispatcherService(log, channel_list)
	notification_consumer := &NotificationConsumer{
		logger:  log,
		service: srv,
	}

	consumer, err := kafka.NewKafkaConsumer(&cfg.Kafka, NOTIFICATION_GROUP)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	err_ch := make(chan error, 1)
	go func() {
		if err := consumer.Consume(ctx, []string{NOTIFICATION_TOPIC}, notification_consumer); err != nil {
			err_ch <- err
		}
	}()

	select {
	case sig := <-quit:
		log.Info("shutting down", zap.String("signal", sig.String()))
		defer cancel()
		err := consumer.Close()
		if err != nil {
			log.Error("failed to shut down kafka notification consumer", zap.Error(err))
		}
	case err := <-err_ch:
		if err != nil {
			log.Fatal("failed to consume messages", zap.Error(err))
		}
	}
}
