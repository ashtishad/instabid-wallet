package utils

import (
	"errors"
	"regexp"

	"github.com/ashtishad/instabid-wallet/lib"
	"github.com/ashtishad/instabid-wallet/user-api/internal/domain"
)

// ValidateCreateUserInput validates the input dto for creating a new user with the following criteria:
//   - Email: Must match the specified regex pattern (EmailRegex).
//   - Password: Must be at least 8 characters long and no more than 32 characters.
//   - Username: Must be 7-64 alphanumeric characters with no spaces.
//   - Status: If provided, must be one of 'active', 'inactive', or 'deleted'.
//   - Role: If provided, must be one of 'user', 'admin', 'moderator', or 'merchant'.
func ValidateCreateUserInput(input domain.NewUserReqDTO) lib.APIError {
	var errs error
	var err error

	if err = lib.ValidateEmail(input.Email); err != nil {
		errs = errors.Join(errs, err)
	}

	if err = lib.ValidatePassword(input.Password); err != nil {
		errs = errors.Join(errs, err)
	}

	if err = lib.ValidateUserName(input.UserName); err != nil {
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

// validateStatus checks status must be one of: active, inactive, deleted
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

// validateFirstName validates first name must be alphabetic and between 1 and 64 characters long
func validateFirstName(firstName string) error {
	if matched := regexp.MustCompile(`^[a-zA-Z]{1,64}$`).MatchString(firstName); !matched {
		return errors.New("first name must be alphabetic and between 1 and 64 characters long")
	}

	return nil
}

// validateLastName validates last name must be alphabetic, may contain spaces, and be between 1 and 128 characters long
func validateLastName(lastName string) error {
	if matched := regexp.MustCompile(`^[a-zA-Z\s]{1,128}$`).MatchString(lastName); !matched {
		return errors.New("last name must be alphabetic, may contain spaces, and be between 1 and 128 characters long")
	}

	return nil
}

// validateGender validates gender must be one of: male, female, other
func validateGender(gender string) error {
	if matched := regexp.MustCompile(`^(male|female|other)$`).MatchString(gender); !matched {
		return errors.New("gender must be one of: male, female, other")
	}

	return nil
}

// validateAddress validates address (optional field) and address cannot exceed 256 characters
func validateAddress(address string) error {
	if len(address) > 256 {
		return errors.New("address cannot exceed 256 characters")
	}

	return nil
}
