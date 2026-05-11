package repositories

import (
	"fmt"
	"reflect"

	"github.com/bigdann09/notifications/internal/dtos"
	"github.com/bigdann09/notifications/internal/models"
	"github.com/bigdann09/notifications/pkgs/pagination"
	"gorm.io/gorm"
)

type INotificationRepository interface {
	FindAllPaginated(query *dtos.NotificationQuery) (*pagination.Pagination[models.Notification], error)
}

type NotificationRepository struct {
	db    *gorm.DB
	table string
}

func NewNotificationRepository(db *gorm.DB) INotificationRepository {
	return &NotificationRepository{db: db, table: "notifications"}
}

func (repository *NotificationRepository) FindAllPaginated(query *dtos.NotificationQuery) (*pagination.Pagination[models.Notification], error) {
	query.Default()
	queryable := repository.db.Table(repository.table)
	if !reflect.DeepEqual(query.Type, nil) {
		queryable.Where("type = ?", query.Type)
	}
	if !reflect.DeepEqual(query.IsRead, nil) {
		if query.IsRead {
			queryable.Where("read_at IS NOT NULL")
		} else {
			queryable.Where("read_at IS NULL")
		}
	}
	if !reflect.DeepEqual(query.SortBy, nil) {
		queryable.Order(fmt.Sprintf("%s %s", query.SortBy, query.Order))
	}

	result := pagination.NewPagination[models.Notification](queryable, query.Page, query.Limit)
	return result, nil
}
