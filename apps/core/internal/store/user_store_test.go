package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/arthurdotwork/dreamtrader/core/internal/entity"
	"github.com/arthurdotwork/dreamtrader/core/internal/store"
	"github.com/arthurdotwork/dreamtrader/core/pkg/psql"
	"github.com/arthurdotwork/dreamtrader/core/pkg/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := psql.Connect(ctx, "postgres", "postgres", "localhost", "5432", "postgres")
	require.NoError(t, err)

	txn, rollback := test.Txn(t, ctx, db)
	t.Cleanup(rollback)

	userStore := store.NewUserStore(txn)

	t.Run("it should create a user", func(t *testing.T) {
		user := entity.User{
			Email:     "mail@domain.tld",
			Password:  "password",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		createdUser, err := userStore.CreateUser(ctx, user)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, createdUser.ID)
		require.NotEmpty(t, createdUser.InternalID)
	})

	t.Run("it should reactivate a deleted user", func(t *testing.T) {
		user := entity.User{
			Email:     "deleted@domain.tld",
			Password:  "password",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		createdUser, err := userStore.CreateUser(ctx, user)
		require.NoError(t, err)

		_, err = txn(ctx).ExecContext(ctx, `UPDATE users SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), createdUser.ID)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := txn(ctx).ExecContext(ctx, `UPDATE users SET deleted_at = NULL WHERE id = $1`, createdUser.ID)
			require.NoError(t, err)
		})

		reactivatedUser, err := userStore.CreateUser(ctx, user)
		require.NoError(t, err)
		require.Equal(t, createdUser.ID, reactivatedUser.ID)
		require.Nil(t, reactivatedUser.DeletedAt)
		require.WithinDuration(t, time.Now().UTC(), reactivatedUser.UpdatedAt, time.Second)
	})
}

func TestGetUserByEmail(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := psql.Connect(ctx, "postgres", "postgres", "localhost", "5432", "postgres")
	require.NoError(t, err)

	txn, rollback := test.Txn(t, ctx, db)
	t.Cleanup(rollback)

	userStore := store.NewUserStore(txn)

	user := entity.User{
		Email:     "mail@domain.tld",
		Password:  "password",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	createdUser, err := userStore.CreateUser(ctx, user)
	require.NoError(t, err)

	t.Run("it should return an error if it can not find the user", func(t *testing.T) {
		_, err := txn(ctx).ExecContext(ctx, `UPDATE users SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), createdUser.ID)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := txn(ctx).ExecContext(ctx, `UPDATE users SET deleted_at = null WHERE id = $1`, createdUser.ID)
			require.NoError(t, err)
		})

		_, err = userStore.GetUserByEmail(ctx, createdUser.Email)
		require.Error(t, err)
		require.ErrorIs(t, err, store.ErrNotFound)
	})

	t.Run("it should return the user", func(t *testing.T) {
		foundUser, err := userStore.GetUserByEmail(ctx, createdUser.Email)
		require.NoError(t, err)
		require.Equal(t, createdUser.ID, foundUser.ID)
		require.Equal(t, createdUser.Email, foundUser.Email)
	})
}
