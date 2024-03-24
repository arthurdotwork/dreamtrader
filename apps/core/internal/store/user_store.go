package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/arthurdotwork/dreamtrader/core/internal/entity"
	"github.com/arthurdotwork/dreamtrader/core/pkg/psql"
)

type UserStore interface {
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
}

type userStore struct {
	db psql.DBGetter
}

func NewUserStore(db psql.DBGetter) UserStore {
	return userStore{db: db}
}

func (s userStore) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	query := `
		INSERT INTO users (email, password, created_at, updated_at)
		VALUES (:email, :password, :created_at, :updated_at)
		ON CONFLICT (email) DO UPDATE SET deleted_at = NULL, updated_at = :updated_at,  password = :password
		RETURNING id, internal_id, deleted_at
	`

	stmt, err := s.db(ctx).PrepareNamedContext(ctx, query)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	defer stmt.Close()

	if err := stmt.GetContext(ctx, &user, user); err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s userStore) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	query := `
		SELECT id, internal_id, email, password, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	var user entity.User
	if err := s.db(ctx).GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, psql.ErrNoRows) {
			return entity.User{}, fmt.Errorf("failed to get user by email: %w", ErrNotFound)
		}

		return entity.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}
