package service

import (
	"errors"

	"github.com/ashtishad/instabid-wallet/auth-api/domain"
	"github.com/ashtishad/instabid-wallet/lib"
)

// validateLoginRequest validates the fields of a LoginRequest.
// Either Username or Email must be provided, along with a Password.
// It returns errors if the validation fails.
func validateLoginRequest(req domain.LoginRequest) lib.APIError {
	var errs error

	if req.Username != "" && req.Email != "" {
		return lib.BadRequestError("user can't sign in with both username and email")
	}

	if (req.Username == "" && req.Email == "") || req.Password == "" {
		return lib.BadRequestError("either username or email, along with a password, must be provided")
	}

	if req.Username != "" {
		if err := lib.ValidateUserName(req.Username); err != nil {
			errs = err
		}
	}

	if req.Email != "" {
		if err := lib.ValidateEmail(req.Email); err != nil {
			errs = err
		}
	}

	if err := lib.ValidatePassword(req.Password); err != nil {
		errs = errors.Join(errs, err)
	}

	if errs != nil {
		return lib.BadRequestError(errs.Error())
	}

	return nil
}
