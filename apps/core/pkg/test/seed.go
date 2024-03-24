package test

import (
	"context"
	"testing"
	"time"

	"github.com/arthurdotwork/dreamtrader/core/internal/entity"
	"github.com/arthurdotwork/dreamtrader/core/internal/store"
	"github.com/arthurdotwork/dreamtrader/core/pkg/psql"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserWithOptions func(*entity.User)

func CreateUserWithEmailOption(email string) CreateUserWithOptions {
	return func(u *entity.User) {
		u.Email = email
	}
}

func CreateUser(ctx context.Context, t *testing.T, txn psql.DBGetter, opts ...CreateUserWithOptions) entity.User {
	t.Helper()

	userStore := store.NewUserStore(txn)

	user := entity.User{
		Email:     "mail@domain.tld",
		Password:  "password",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	for _, opt := range opts {
		opt(&user)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	require.NoError(t, err)
	user.Password = string(hashedPassword)

	createdUser, err := userStore.CreateUser(ctx, user)
	require.NoError(t, err)

	return createdUser
}
