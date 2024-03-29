package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/arthurdotwork/dreamtrader/core/internal/entity"
	"github.com/arthurdotwork/dreamtrader/core/internal/request"
	"github.com/arthurdotwork/dreamtrader/core/internal/service"
	"github.com/arthurdotwork/dreamtrader/core/internal/store"
	"github.com/arthurdotwork/dreamtrader/core/pkg/psql"
	"github.com/arthurdotwork/dreamtrader/core/pkg/test"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := psql.Connect(ctx, "postgres", "postgres", "localhost", "5432", "postgres")
	require.NoError(t, err)

	txn, rollback := test.Txn(t, ctx, db)
	t.Cleanup(rollback)

	authStore := store.NewAuthStore(txn)
	userStore := store.NewUserStore(txn)

	authService := service.NewAuthService(authStore, userStore)

	createdUser := test.CreateUser(ctx, t, txn)

	t.Run("it should return an error if the user does not exist", func(t *testing.T) {
		_, err := txn(ctx).ExecContext(ctx, `UPDATE users SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), createdUser.ID)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := txn(ctx).ExecContext(ctx, `UPDATE users SET deleted_at = NULL WHERE id = $1`, createdUser.ID)
			require.NoError(t, err)
		})

		_, _, err = authService.Authenticate(ctx, request.AuthenticateRequest{Email: createdUser.Email, Password: "password"})
		require.Error(t, err)
	})

	t.Run("it should return an error if the password is incorrect", func(t *testing.T) {
		_, _, err := authService.Authenticate(ctx, request.AuthenticateRequest{Email: createdUser.Email, Password: "wrong-password"})
		require.Error(t, err)
	})

	t.Run("it should authenticate the user", func(t *testing.T) {
		accessToken, refreshToken, err := authService.Authenticate(ctx, request.AuthenticateRequest{Email: createdUser.Email, Password: "password"})
		require.NoError(t, err)

		require.Equal(t, accessToken.UserID, createdUser.ID)
		require.NotEmpty(t, accessToken.AccessToken)
		require.Greater(t, accessToken.ExpiresAt, time.Now().UTC())

		require.Equal(t, refreshToken.UserID, createdUser.ID)
		require.NotEmpty(t, refreshToken.RefreshToken)
		require.Greater(t, refreshToken.ExpiresAt, time.Now().UTC())
	})
}

func TestVerifyAuthentication(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := psql.Connect(ctx, "postgres", "postgres", "localhost", "5432", "postgres")
	require.NoError(t, err)

	txn, rollback := test.Txn(t, ctx, db)
	t.Cleanup(rollback)

	authStore := store.NewAuthStore(txn)
	userStore := store.NewUserStore(txn)

	authService := service.NewAuthService(authStore, userStore)

	createdUser := test.CreateUser(ctx, t, txn)

	accessToken := entity.AuthAccessToken{
		UserID:      createdUser.ID,
		AccessToken: "access_token",
		ExpiresAt:   time.Now().UTC().Add(time.Hour * 2),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	createdAccessToken, err := authStore.CreateAuthAccessToken(ctx, accessToken)
	require.NoError(t, err)

	t.Run("it should return an error if there is no valid token", func(t *testing.T) {
		_, err := txn(ctx).ExecContext(ctx, `UPDATE auth_access_tokens SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), createdAccessToken.ID)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := txn(ctx).ExecContext(ctx, `UPDATE auth_access_tokens SET deleted_at = NULL WHERE id = $1`, createdAccessToken.ID)
			require.NoError(t, err)
		})

		err = authService.Verify(ctx, createdAccessToken.AccessToken)
		require.Error(t, err)
	})

	t.Run("it should return an error if the token is expired", func(t *testing.T) {
		_, err := txn(ctx).ExecContext(ctx, `UPDATE auth_access_tokens SET expires_at = $1 WHERE id = $2`, time.Now().UTC().Add(-time.Hour), createdAccessToken.ID)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := txn(ctx).ExecContext(ctx, `UPDATE auth_access_tokens SET expires_at = $1 WHERE id = $2`, time.Now().UTC().Add(time.Hour*2), createdAccessToken.ID)
			require.NoError(t, err)
		})

		err = authService.Verify(ctx, createdAccessToken.AccessToken)
		require.Error(t, err)
	})

	t.Run("it should verify the authentication", func(t *testing.T) {
		err := authService.Verify(ctx, createdAccessToken.AccessToken)
		require.NoError(t, err)
	})
}

func TestRefreshAuthentication(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := psql.Connect(ctx, "postgres", "postgres", "localhost", "5432", "postgres")
	require.NoError(t, err)

	txn, rollback := test.Txn(t, ctx, db)
	t.Cleanup(rollback)

	authStore := store.NewAuthStore(txn)
	userStore := store.NewUserStore(txn)

	authService := service.NewAuthService(authStore, userStore)

	createdUser := test.CreateUser(ctx, t, txn)

	refreshToken := entity.AuthRefreshToken{
		UserID:       createdUser.ID,
		RefreshToken: "refreshToken",
		ExpiresAt:    time.Now().UTC().Add(time.Hour * 2),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	createdRefreshToken, err := authStore.CreateAuthRefreshToken(ctx, refreshToken)
	require.NoError(t, err)

	t.Run("it should return an error if it can not find the refresh token", func(t *testing.T) {
		_, err := txn(ctx).ExecContext(ctx, `UPDATE auth_refresh_tokens SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), createdRefreshToken.ID)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := txn(ctx).ExecContext(ctx, `UPDATE auth_refresh_tokens SET deleted_at = NULL WHERE id = $1`, createdRefreshToken.ID)
			require.NoError(t, err)
		})

		_, _, err = authService.RefreshAuthentication(ctx, createdRefreshToken.RefreshToken)
		require.Error(t, err)
	})

	t.Run("it should return an error if the refresh token is expired", func(t *testing.T) {
		_, err := txn(ctx).ExecContext(ctx, `UPDATE auth_refresh_tokens SET expires_at = $1 WHERE id = $2`, time.Now().UTC().Add(-time.Hour), createdRefreshToken.ID)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := txn(ctx).ExecContext(ctx, `UPDATE auth_refresh_tokens SET expires_at = $1 WHERE id = $2`, time.Now().UTC().Add(time.Hour*2), createdRefreshToken.ID)
			require.NoError(t, err)
		})

		_, _, err = authService.RefreshAuthentication(ctx, createdRefreshToken.RefreshToken)
		require.Error(t, err)
	})

	t.Run("it should refresh the access token", func(t *testing.T) {
		accessToken, refreshToken, err := authService.RefreshAuthentication(ctx, createdRefreshToken.RefreshToken)
		require.NoError(t, err)
		require.Equal(t, accessToken.UserID, createdUser.ID)
		require.Greater(t, accessToken.ExpiresAt, time.Now().UTC())
		require.Equal(t, refreshToken.ID, createdRefreshToken.ID)
	})
}
