package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arthurdotwork/dreamtrader/core/internal/entity"
	"github.com/arthurdotwork/dreamtrader/core/internal/request"
	"github.com/arthurdotwork/dreamtrader/core/internal/store"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type RegisterService struct {
	userStore store.UserStore
}

func NewRegisterService(userStore store.UserStore) RegisterService {
	return RegisterService{
		userStore: userStore,
	}
}

func (s RegisterService) Register(ctx context.Context, req request.CreateUserRequest) (entity.User, error) {
	user, err := s.userStore.GetUserByEmail(ctx, req.Email)
	switch {
	case err == nil:
		// user already exists, do not display an error to prevent user enumeration
		return user, nil
	case !errors.Is(err, store.ErrNotFound):
		log.Ctx(ctx).Error().Err(err).Msg("failed to get user by email")
		return entity.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to hash password")
		return entity.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user = entity.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	createdUser, err := s.userStore.CreateUser(ctx, user)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to create user")
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}
