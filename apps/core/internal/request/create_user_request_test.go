package request_test

import (
	"testing"

	"github.com/arthurdotwork/dreamtrader/core/internal/request"
	"github.com/stretchr/testify/require"
)

func TestValidateCreateUserRequest(t *testing.T) {
	t.Parallel()

	t.Run("it should return an error if the email is invalid", func(t *testing.T) {
		t.Parallel()

		req := request.CreateUserRequest{Email: "invalid"}
		err := req.Validate()
		require.Error(t, err)
	})

	t.Run("it should return an error if the password is too short", func(t *testing.T) {
		t.Parallel()

		req := request.CreateUserRequest{Email: "mail@domain.tld", Password: "short"}
		err := req.Validate()
		require.Error(t, err)
	})

	t.Run("it should return an error if the password is made of empty spaces", func(t *testing.T) {
		t.Parallel()

		req := request.CreateUserRequest{Email: "mail@domain.tld", Password: "        "}
		err := req.Validate()
		require.Error(t, err)
	})

	t.Run("it should validate the request", func(t *testing.T) {
		t.Parallel()

		req := request.CreateUserRequest{Email: "mail@domain.tld", Password: "password"}
		err := req.Validate()
		require.NoError(t, err)
	})
}
