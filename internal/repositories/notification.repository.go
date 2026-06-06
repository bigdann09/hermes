package repositories

import (
	"fmt"

	"github.com/bigdann09/notifications/internal/dtos"
	"github.com/bigdann09/notifications/internal/models"
	"github.com/bigdann09/notifications/pkgs/pagination"
	"gorm.io/gorm"
)

type INotificationRepository interface {
	Create(notification *models.Notification) error
	FindAllPaginated(query *dtos.NotificationQuery) (*pagination.Pagination[models.Notification], error)
}

type NotificationRepository struct {
	db    *gorm.DB
	table string
}

func NewNotificationRepository(db *gorm.DB) INotificationRepository {
	return &NotificationRepository{db: db, table: "notifications"}
}

func (repository *NotificationRepository) Create(notification *models.Notification) error {
	return repository.db.Table(repository.table).Create(notification).Error
}

func (repository *NotificationRepository) FindAllPaginated(query *dtos.NotificationQuery) (*pagination.Pagination[models.Notification], error) {
	query.Default()
	queryable := repository.db.Table(repository.table).Where("user_id = ?", query.UserID)
	if query.Type != "" {
		queryable = queryable.Where("type = ?", query.Type)
	}
	if query.IsRead {
		queryable = queryable.Where("read_at IS NOT NULL")
	}
	if query.SortBy != "" {
		queryable = queryable.Order(fmt.Sprintf("%s %s", query.SortBy, query.Order))
	}

	result := pagination.NewPagination[models.Notification](queryable, query.Page, query.Limit)
	return result, nil
}
