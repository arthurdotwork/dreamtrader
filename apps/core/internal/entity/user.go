package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID  `db:"id"`
	InternalID int64      `db:"internal_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`

	Email    string `db:"email"`
	Password string `db:"password"`
}
