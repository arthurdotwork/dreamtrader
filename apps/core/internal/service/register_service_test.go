package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/arthurdotwork/dreamtrader/core/internal/request"
	"github.com/arthurdotwork/dreamtrader/core/internal/service"
	"github.com/arthurdotwork/dreamtrader/core/internal/store"
	"github.com/arthurdotwork/dreamtrader/core/pkg/psql"
	"github.com/arthurdotwork/dreamtrader/core/pkg/test"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := psql.Connect(ctx, "postgres", "postgres", "localhost", "5432", "postgres")
	require.NoError(t, err)

	txn, rollback := test.Txn(t, ctx, db)
	t.Cleanup(rollback)

	userStore := store.NewUserStore(txn)
	registerService := service.NewRegisterService(userStore)

	createdUser := test.CreateUser(ctx, t, txn)

	t.Run("it should return no error if the user already exists", func(t *testing.T) {
		user, err := registerService.Register(ctx, request.CreateUserRequest{Email: "mail@domain.tld", Password: "password"})
		require.NoError(t, err)
		require.Equal(t, createdUser.ID, user.ID)
	})

	t.Run("it should create a new user", func(t *testing.T) {
		user, err := registerService.Register(ctx, request.CreateUserRequest{Email: "user@domain.tld", Password: "password"})
		require.NoError(t, err)
		require.NotEqual(t, createdUser.ID, user.ID)
		require.NotEqual(t, "password", user.Password)
	})

	t.Run("a deactivated user should be able to register with the same email", func(t *testing.T) {
		_, err := txn(ctx).ExecContext(ctx, `UPDATE users SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), createdUser.ID)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := txn(ctx).ExecContext(ctx, `UPDATE users SET deleted_at = NULL WHERE id = $1`, createdUser.ID)
			require.NoError(t, err)
		})

		user, err := registerService.Register(ctx, request.CreateUserRequest{Email: "mail@domain.tld", Password: "password"})
		require.NoError(t, err)
		require.Equal(t, createdUser.ID, user.ID)
		require.Nil(t, user.DeletedAt)
		require.WithinDuration(t, time.Now().UTC(), user.UpdatedAt, time.Second)
	})
}
