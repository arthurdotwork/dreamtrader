package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuthAccessToken struct {
	ID         uuid.UUID  `db:"id"`
	InternalID int64      `db:"internal_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`

	UserID      uuid.UUID `db:"user_id"`
	AccessToken string    `db:"access_token"`
	ExpiresAt   time.Time `db:"expires_at"`
}

type AuthRefreshToken struct {
	ID         uuid.UUID  `db:"id"`
	InternalID int64      `db:"internal_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`

	UserID       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
}
