package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/arthurdotwork/dreamtrader/core/internal/entity"
	"github.com/arthurdotwork/dreamtrader/core/internal/store"
	"github.com/arthurdotwork/dreamtrader/core/pkg/psql"
	"github.com/arthurdotwork/dreamtrader/core/pkg/test"
	"github.com/stretchr/testify/require"
)

func TestCreateAuthAccessToken(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := psql.Connect(ctx, "postgres", "postgres", "localhost", "5432", "postgres")
	require.NoError(t, err)

	txn, rollback := test.Txn(t, ctx, db)
	t.Cleanup(rollback)

	authStore := store.NewAuthStore(txn)

	user := test.CreateUser(ctx, t, txn)

	t.Run("it should create an auth access token", func(t *testing.T) {
		accessToken := entity.AuthAccessToken{
			UserID:      user.ID,
			AccessToken: "access_token",
			ExpiresAt:   time.Now().UTC().Add(time.Hour * 2),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}

		createdAccessToken, err := authStore.CreateAuthAccessToken(ctx, accessToken)
		require.NoError(t, err)
		require.NotEmpty(t, createdAccessToken.ID)
		require.NotEmpty(t, createdAccessToken.InternalID)
	})
}

func TestGetAuthAccessTokenByAccessToken(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := psql.Connect(ctx, "postgres", "postgres", "localhost", "5432", "postgres")
	require.NoError(t, err)

	txn, rollback := test.Txn(t, ctx, db)
	t.Cleanup(rollback)

	authStore := store.NewAuthStore(txn)

	user := test.CreateUser(ctx, t, txn)

	accessToken := entity.AuthAccessToken{
		UserID:      user.ID,
		AccessToken: "access_token",
		ExpiresAt:   time.Now().UTC().Add(time.Hour * 2),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	createdAccessToken, err := authStore.CreateAuthAccessToken(ctx, accessToken)
	require.NoError(t, err)

	t.Run("it should return an error if it can not find the auth access token", func(t *testing.T) {
		_, err := txn(ctx).ExecContext(ctx, `UPDATE auth_access_tokens SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), createdAccessToken.ID)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := txn(ctx).ExecContext(ctx, `UPDATE auth_access_tokens SET deleted_at = NULL WHERE id = $1`, createdAccessToken.ID)
			require.NoError(t, err)
		})

		_, err = authStore.GetAuthAccessTokenByAccessToken(ctx, createdAccessToken.AccessToken)
		require.Error(t, err)
	})

	t.Run("it should return the auth access token", func(t *testing.T) {
		authAccessToken, err := authStore.GetAuthAccessTokenByAccessToken(ctx, createdAccessToken.AccessToken)
		require.NoError(t, err)
		require.Equal(t, createdAccessToken.ID, authAccessToken.ID)
	})
}

func TestCreateAuthRefreshToken(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := psql.Connect(ctx, "postgres", "postgres", "localhost", "5432", "postgres")
	require.NoError(t, err)

	txn, rollback := test.Txn(t, ctx, db)
	t.Cleanup(rollback)

	authStore := store.NewAuthStore(txn)

	user := test.CreateUser(ctx, t, txn)

	t.Run("it should create an auth access token", func(t *testing.T) {
		refreshToken := entity.AuthRefreshToken{
			UserID:       user.ID,
			RefreshToken: "refresh_token",
			ExpiresAt:    time.Now().UTC().Add(time.Hour * 2),
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		}

		createdAccessToken, err := authStore.CreateAuthRefreshToken(ctx, refreshToken)
		require.NoError(t, err)
		require.NotEmpty(t, createdAccessToken.ID)
		require.NotEmpty(t, createdAccessToken.InternalID)
	})
}

func TestGetAuthRefreshTokenByRefreshToken(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := psql.Connect(ctx, "postgres", "postgres", "localhost", "5432", "postgres")
	require.NoError(t, err)

	txn, rollback := test.Txn(t, ctx, db)
	t.Cleanup(rollback)

	authStore := store.NewAuthStore(txn)

	user := test.CreateUser(ctx, t, txn)

	refreshToken := entity.AuthRefreshToken{
		UserID:       user.ID,
		RefreshToken: "refresh_token",
		ExpiresAt:    time.Now().UTC().Add(time.Hour * 2),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	createdRefreshToken, err := authStore.CreateAuthRefreshToken(ctx, refreshToken)
	require.NoError(t, err)

	t.Run("it should return an error if it can not find the auth refresh token", func(t *testing.T) {
		_, err := txn(ctx).ExecContext(ctx, `UPDATE auth_refresh_tokens SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), createdRefreshToken.ID)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := txn(ctx).ExecContext(ctx, `UPDATE auth_refresh_tokens SET deleted_at = NULL WHERE id = $1`, createdRefreshToken.ID)
			require.NoError(t, err)
		})

		_, err = authStore.GetAuthRefreshTokenByRefreshToken(ctx, createdRefreshToken.RefreshToken)
		require.Error(t, err)
	})

	t.Run("it should return the auth refresh token", func(t *testing.T) {
		authRefreshToken, err := authStore.GetAuthRefreshTokenByRefreshToken(ctx, createdRefreshToken.RefreshToken)
		require.NoError(t, err)
		require.Equal(t, createdRefreshToken.ID, authRefreshToken.ID)
	})
}
