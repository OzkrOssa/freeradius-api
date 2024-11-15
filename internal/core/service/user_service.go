package service

import (
	"context"
	"errors"
	"github.com/OzkrOssa/freeradius-api/internal/core/domain"
	"github.com/OzkrOssa/freeradius-api/internal/core/port"
	"github.com/OzkrOssa/freeradius-api/internal/core/utils"
)

type UserService struct {
	repo  port.UserRepository
	cache port.CacheRepository
}

func NewUserService(repo port.UserRepository, cache port.CacheRepository) *UserService {
	return &UserService{repo, cache}
}

func (u UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, domain.InternalError
	}

	user.Password = hashedPassword

	user, err = u.repo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, domain.ConflictDataError) {
			return nil, err
		}
		return nil, domain.ConflictDataError
	}

	key := utils.GenerateCacheKey("user", user.ID)

	serializedUser, err := utils.Serialize(user)
	if err != nil {
		return nil, domain.InternalError
	}

	err = u.cache.Set(ctx, key, serializedUser, 0)
	if err != nil {
		return nil, domain.InternalError
	}

	err = u.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, domain.InternalError
	}

	return user, nil
}

func (u UserService) GetUser(ctx context.Context, id uint64) (*domain.User, error) {
	cacheKey := utils.GenerateCacheKey("user", id)
	cachedUser, err := u.cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, domain.InternalError
	}

	var user *domain.User

	if cachedUser != nil {
		var user *domain.User
		err := utils.Deserialize(cachedUser, &user)
		if err != nil {
			return nil, domain.InternalError
		}
		return user, nil
	}

	user, err = u.repo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, domain.DataNotFoundError) {
			return nil, err
		}
		return nil, domain.InternalError
	}

	userSerialized, err := utils.Serialize(user)
	if err != nil {
		return nil, domain.InternalError
	}

	err = u.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, domain.InternalError
	}

	return user, nil

}

func (u UserService) ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error) {

	var users []domain.User

	params := utils.GenerateCacheKeyParams(skip, limit)
	cacheKey := utils.GenerateCacheKey("users", params)

	cachedUsers, err := u.cache.Get(ctx, cacheKey)
	if err == nil {
		err := utils.Deserialize(cachedUsers, &users)
		if err != nil {
			return nil, domain.InternalError
		}
		return users, nil
	}

	users, err = u.repo.ListUsers(ctx, skip, limit)
	if err != nil {
		return nil, domain.InternalError
	}

	usersSerialized, err := utils.Serialize(users)
	if err != nil {
		return nil, domain.InternalError
	}

	err = u.cache.Set(ctx, cacheKey, usersSerialized, 0)
	if err != nil {
		return nil, domain.InternalError
	}

	return users, nil
}

func (u UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := u.repo.GetUserById(ctx, user.ID)

	if err != nil {
		if errors.Is(err, domain.DataNotFoundError) {
			return nil, err
		}
		return nil, domain.InternalError
	}

	emptyData := user.Name == "" &&
		user.Email == "" &&
		user.Password == ""

	sameData := existingUser.Name == user.Name &&
		existingUser.Email == user.Email

	if emptyData || sameData {
		return nil, domain.ConflictDataError
	}

	var hashedPassword string

	if user.Password != "" {
		hashedPassword, err = utils.HashPassword(user.Password)
		if err != nil {
			return nil, domain.ConflictDataError
		}
	}

	user.Password = hashedPassword

	_, err = u.repo.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, domain.DataNotFoundError) {
			return nil, err
		}
		return nil, domain.InternalError
	}

	cacheKey := utils.GenerateCacheKey("user", user.ID)

	err = u.cache.Delete(ctx, cacheKey)
	if err != nil {
		return nil, domain.InternalError
	}

	userSerialized, err := utils.Serialize(user)
	if err != nil {
		return nil, err
	}

	err = u.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, domain.InternalError
	}

	err = u.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, domain.InternalError
	}

	return user, nil
}

func (u UserService) DeleteUser(ctx context.Context, id uint64) error {
	_, err := u.repo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, domain.DataNotFoundError) {
			return err
		}
		return domain.ConflictDataError
	}

	cacheKey := utils.GenerateCacheKey("user", id)

	err = u.cache.Delete(ctx, cacheKey)
	if err != nil {
		return domain.InternalError
	}

	err = u.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return domain.InternalError
	}

	return u.repo.DeleteUser(ctx, id)
}
