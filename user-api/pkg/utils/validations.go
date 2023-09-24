package utils

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/ashtishad/instabid-wallet/lib"
	"github.com/ashtishad/instabid-wallet/user-api/internal/domain"
)

func validateEmail(email string) error {
	if matched := regexp.MustCompile(EmailRegex).MatchString(email); !matched {
		return fmt.Errorf("invalid email, you entered %s", email)
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 || len(password) > 32 {
		return errors.New("password must be at least 8 characters long and no more than 32 characters")
	}

	return nil
}

func validateUserName(userName string) error {
	if len(userName) < 7 || len(userName) > 64 {
		return errors.New("username must be between 7 and 64 characters long")
	}

	return nil
}

func validateStatus(status string) error {
	if matched := regexp.MustCompile(StatusRegex).MatchString(status); !matched && status != "" {
		return errors.New("status must be one of: active, inactive, deleted")
	}

	return nil
}

func validateRole(role string) error {
	if matched := regexp.MustCompile(RoleRegex).MatchString(role); !matched && role != "" {
		return errors.New("role must be one of: user, admin, moderator, merchant")
	}

	return nil
}

// ValidateCreateUserInput validates the input dto for creating a new user with the following criteria:
//   - Email: Must match the specified regex pattern (EmailRegex).
//   - Password: Must be at least 8 characters long and no more than 32 characters.
//   - UserName: Must be between 7 and 64 characters long.
//   - Status: If provided, must be one of 'active', 'inactive', or 'deleted'.
//   - Role: If provided, must be one of 'user', 'admin', 'moderator', or 'merchant'.
func ValidateCreateUserInput(input domain.NewUserReqDTO) lib.APIError {
	var errs error
	var err error

	if err = validateEmail(input.Email); err != nil {
		errs = errors.Join(errs, err)
	}

	if err = validatePassword(input.Password); err != nil {
		errs = errors.Join(errs, err)
	}

	if err = validateUserName(input.UserName); err != nil {
		errs = errors.Join(errs, err)
	}

	if err = validateStatus(input.Status); err != nil {
		errs = errors.Join(errs, err)
	}

	if err = validateRole(input.Role); err != nil {
		errs = errors.Join(errs, err)
	}

	if errs != nil {
		return lib.BadRequestError(errs.Error())
	}

	return nil
}
