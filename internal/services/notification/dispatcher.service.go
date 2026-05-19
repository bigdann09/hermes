package notification

import (
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type DispatcherService struct {
	logger *zap.Logger
}

func NewDispatcherService(logger *zap.Logger) *DispatcherService {
	return &DispatcherService{
		logger: logger,
	}
}

func (dispatcher *DispatcherService) Dispatch(msg *sarama.ConsumerMessage) {
	fmt.Println("in the dispatcher", string(msg.Value))
}
