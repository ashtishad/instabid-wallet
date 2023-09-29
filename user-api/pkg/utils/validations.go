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

// ValidateFirstName validates first name
func validateFirstName(firstName string) error {
	if matched := regexp.MustCompile(`^[a-zA-Z]{1,64}$`).MatchString(firstName); !matched {
		return errors.New("first name must be alphabetic and between 1 and 64 characters long")
	}

	return nil
}

// ValidateLastName validates last name
func validateLastName(lastName string) error {
	if matched := regexp.MustCompile(`^[a-zA-Z\s]{1,128}$`).MatchString(lastName); !matched {
		return errors.New("last name must be alphabetic, may contain spaces, and be between 1 and 128 characters long")
	}

	return nil
}

// ValidateGender validates gender
func validateGender(gender string) error {
	if matched := regexp.MustCompile(`^(male|female|other)$`).MatchString(gender); !matched {
		return errors.New("gender must be one of: male, female, other")
	}

	return nil
}

// ValidateAddress validates address (optional field)
func validateAddress(address string) error {
	if len(address) > 256 {
		return errors.New("address cannot exceed 256 characters")
	}

	return nil
}

// ValidateCreateProfileInput validates the input dto for creating a new profile with the following criteria:
//   - FirstName: Must be alphabetic and between 1 and 64 characters long.
//   - LastName: Must be alphabetic, may contain spaces, and be between 1 and 128 characters long.
//   - Gender: Must be one of 'male', 'female', or 'other'.
//   - Address: If provided, must not exceed 256 characters.
func ValidateCreateProfileInput(input domain.NewProfileReqDTO) lib.APIError {
	var errs error
	var err error

	if err = validateFirstName(input.FirstName); err != nil {
		errs = errors.Join(errs, err)
	}

	if err = validateLastName(input.LastName); err != nil {
		errs = errors.Join(errs, err)
	}

	if err = validateGender(input.Gender); err != nil {
		errs = errors.Join(errs, err)
	}

	if err = validateAddress(input.Address); err != nil {
		errs = errors.Join(errs, err)
	}

	if errs != nil {
		return lib.BadRequestError(errs.Error())
	}

	return nil
}
