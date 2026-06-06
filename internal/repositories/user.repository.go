package repositories

import (
	"github.com/bigdann09/notifications/internal/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(user *models.User) error
	FindAll() ([]models.User, error)
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
}

type UserRepository struct {
	db    *gorm.DB
	table string
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db, table: "users"}
}

func (repo *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := repo.db.Table(repo.table).Find(&users).Error
	return users, err
}

func (repo *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	err := repo.db.Table(repo.table).Where("id = ?", id).First(&user).Error
	return &user, err
}

func (repo *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := repo.db.Table(repo.table).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (repo *UserRepository) Create(user *models.User) error {
	return repo.db.Table(repo.table).Create(user).Error
}

func (repo *UserRepository) Update(user *models.User) error {
	return repo.db.Table(repo.table).Updates(user).Error
}

func (repo *UserRepository) Delete(id string) error {
	return repo.db.Table(repo.table).Delete(&models.User{}, id).Error
}
