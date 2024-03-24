package service

import (
	"context"
	"fmt"
	"time"

	"github.com/arthurdotwork/dreamtrader/core/internal/entity"
	"github.com/arthurdotwork/dreamtrader/core/internal/request"
	"github.com/arthurdotwork/dreamtrader/core/internal/store"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authStore store.AuthStore
	userStore store.UserStore
}

func NewAuthService(authStore store.AuthStore, userStore store.UserStore) AuthService {
	return AuthService{
		authStore: authStore,
		userStore: userStore,
	}
}

func (s AuthService) Authenticate(ctx context.Context, req request.AuthenticateRequest) (entity.AuthAccessToken, entity.AuthRefreshToken, error) {
	fmt.Printf("%+v\n", req)

	user, err := s.userStore.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get user by email")
		return entity.AuthAccessToken{}, entity.AuthRefreshToken{}, fmt.Errorf("failed to get user by email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to compare password")
		return entity.AuthAccessToken{}, entity.AuthRefreshToken{}, fmt.Errorf("failed to compare password")
	}

	accessToken := entity.AuthAccessToken{
		UserID:      user.ID,
		AccessToken: uuid.New().String(),
		ExpiresAt:   time.Now().UTC().Add(time.Hour * 2),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	refreshToken := entity.AuthRefreshToken{
		UserID:       user.ID,
		RefreshToken: uuid.New().String(),
		ExpiresAt:    time.Now().UTC().Add(time.Hour * 24 * 7),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if _, err := s.authStore.CreateAuthAccessToken(ctx, accessToken); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to create auth access token")
		return entity.AuthAccessToken{}, entity.AuthRefreshToken{}, fmt.Errorf("failed to create auth access token")
	}

	if _, err := s.authStore.CreateAuthRefreshToken(ctx, refreshToken); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to create auth refresh token")
		return entity.AuthAccessToken{}, entity.AuthRefreshToken{}, fmt.Errorf("failed to create auth refresh token")
	}

	return accessToken, refreshToken, nil
}
