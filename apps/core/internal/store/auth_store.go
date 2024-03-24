package store

import (
	"context"
	"fmt"

	"github.com/arthurdotwork/dreamtrader/core/internal/entity"
	"github.com/arthurdotwork/dreamtrader/core/pkg/psql"
)

type AuthStore interface {
	CreateAuthAccessToken(ctx context.Context, accessToken entity.AuthAccessToken) (entity.AuthAccessToken, error)
	GetAuthAccessTokenByAccessToken(ctx context.Context, accessToken string) (entity.AuthAccessToken, error)

	CreateAuthRefreshToken(ctx context.Context, refreshToken entity.AuthRefreshToken) (entity.AuthRefreshToken, error)
	GetAuthRefreshTokenByRefreshToken(ctx context.Context, refreshToken string) (entity.AuthRefreshToken, error)
}

type authStore struct {
	db psql.DBGetter
}

func NewAuthStore(db psql.DBGetter) AuthStore {
	return authStore{
		db: db,
	}
}

func (s authStore) CreateAuthAccessToken(ctx context.Context, accessToken entity.AuthAccessToken) (entity.AuthAccessToken, error) {
	query := `
		INSERT INTO auth_access_tokens (user_id, access_token, expires_at, created_at, updated_at)
		VALUES (:user_id, :access_token, :expires_at, :created_at, :updated_at)
		RETURNING id, internal_id
	`

	stmt, err := s.db(ctx).PrepareNamed(query)
	if err != nil {
		return entity.AuthAccessToken{}, fmt.Errorf("failed to create auth access token: %w", err)
	}
	defer stmt.Close()

	if err := stmt.GetContext(ctx, &accessToken, accessToken); err != nil {
		return entity.AuthAccessToken{}, fmt.Errorf("failed to create auth access token: %w", err)
	}

	return accessToken, nil
}

func (s authStore) GetAuthAccessTokenByAccessToken(ctx context.Context, accessToken string) (entity.AuthAccessToken, error) {
	query := `
		SELECT id, internal_id, user_id, access_token, expires_at, created_at, updated_at
		FROM auth_access_tokens
		WHERE access_token = $1 AND deleted_at IS NULL
	`

	var authAccessToken entity.AuthAccessToken
	if err := s.db(ctx).GetContext(ctx, &authAccessToken, query, accessToken); err != nil {
		return entity.AuthAccessToken{}, fmt.Errorf("failed to get auth access token by access token: %w", err)
	}

	return authAccessToken, nil
}

func (s authStore) CreateAuthRefreshToken(ctx context.Context, refreshToken entity.AuthRefreshToken) (entity.AuthRefreshToken, error) {
	query := `
		INSERT INTO auth_refresh_tokens (user_id, refresh_token, expires_at, created_at, updated_at)
		VALUES (:user_id, :refresh_token, :expires_at, :created_at, :updated_at)
		RETURNING id, internal_id
	`

	stmt, err := s.db(ctx).PrepareNamed(query)
	if err != nil {
		return entity.AuthRefreshToken{}, fmt.Errorf("failed to create auth refresh token: %w", err)
	}
	defer stmt.Close()

	if err := stmt.GetContext(ctx, &refreshToken, refreshToken); err != nil {
		return entity.AuthRefreshToken{}, fmt.Errorf("failed to create auth refresh token: %w", err)
	}

	return refreshToken, nil
}

func (s authStore) GetAuthRefreshTokenByRefreshToken(ctx context.Context, refreshToken string) (entity.AuthRefreshToken, error) {
	query := `
		SELECT id, internal_id, user_id, refresh_token, expires_at, created_at, updated_at
		FROM auth_refresh_tokens
		WHERE refresh_token = $1 AND deleted_at IS NULL
	`

	var authRefreshToken entity.AuthRefreshToken
	if err := s.db(ctx).GetContext(ctx, &authRefreshToken, query, refreshToken); err != nil {
		return entity.AuthRefreshToken{}, fmt.Errorf("failed to get auth refresh token by refresh token: %w", err)
	}

	return authRefreshToken, nil
}
