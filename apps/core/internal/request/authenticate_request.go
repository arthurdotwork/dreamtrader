package request

import (
	"errors"
	"fmt"
	"strings"
)

type AuthenticateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r AuthenticateRequest) Validate() error {
	var errs []error

	if !strings.Contains(r.Email, "@") {
		errs = append(errs, fmt.Errorf("email is invalid"))
	}

	if len(strings.TrimSpace(r.Password)) < 8 {
		errs = append(errs, fmt.Errorf("password is too short"))
	}

	return errors.Join(errs...)
}
