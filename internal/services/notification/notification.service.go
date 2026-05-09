package notification

type INotificationService interface {
	Send() error
}

type NotificationService struct {
}

func NewNotificationService() INotificationService {
	return &NotificationService{}
}

func (s *NotificationService) Send() error {
	return nil
}
