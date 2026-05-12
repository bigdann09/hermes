package user

import (
	"fmt"
	"time"

	"github.com/bigdann09/notifications/internal/infrastructure/cache"
	"github.com/bigdann09/notifications/internal/models"
	"github.com/bigdann09/notifications/internal/repositories"
	"github.com/bigdann09/notifications/pkgs/apiresponse"
	"go.uber.org/zap"
)

type IUserService interface {
	FindByID(id string) (*models.User, *apiresponse.APIResponse)
}

type UserService struct {
	repository repositories.IUserRepository
	logger     *zap.Logger
	cache      *cache.Cache
}

func NewUserService(repository repositories.IUserRepository, cache *cache.Cache, logger *zap.Logger) *UserService {
	return &UserService{repository: repository, cache: cache, logger: logger}
}

func (srv *UserService) FindByID(id string) (*models.User, *apiresponse.APIResponse) {
	var cached models.User
	cache_key := fmt.Sprintf("user:%s", id)
	if err := srv.cache.Get(cache_key, &cached); err == nil {
		srv.logger.Info("user retrieved successfully from cache")
		return &cached, nil
	}

	user, err := srv.repository.FindByID(id)
	if err != nil {
		srv.logger.Error("could not retrieve user", zap.Error(err))
		return nil, apiresponse.NotFound("User not found")
	}

	srv.cache.Set(cache_key, user, 10*time.Minute)
	return user, nil
}
