package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ashtishad/instabid-wallet/lib"
	"github.com/ashtishad/instabid-wallet/user-api/internal/domain"
)

// ValidateCreateUserInput validates the input dto for creating a new user with the following criteria:
//   - Email: Must match the specified regex pattern (EmailRegex).
//   - Password: Must be at least 8 characters long and no more than 32 characters.
//   - UserName: Must be between 7 and 64 characters long.
//   - Status: If provided, must be one of 'active', 'inactive', or 'deleted'.
//   - Role: If provided, must be one of 'user', 'admin', 'moderator', or 'merchant'.
func ValidateCreateUserInput(input domain.NewUserReqDTO) lib.APIError {
	var errorMessages []string

	if matched := regexp.MustCompile(EmailRegex).MatchString(input.Email); !matched {
		errorMessages = append(errorMessages, fmt.Sprintf("invalid email, you entered %s", input.Email))
	}

	if len(input.Password) < 8 || len(input.Password) > 32 {
		errorMessages = append(errorMessages, "password must be at least 8 characters long and no more than 32 characters")
	}

	if len(input.UserName) < 7 || len(input.UserName) > 64 {
		errorMessages = append(errorMessages, "username must be between 7 and 64 characters long")
	}

	if input.Status != "" {
		if matched := regexp.MustCompile(StatusRegex).MatchString(input.Status); !matched {
			errorMessages = append(errorMessages, "status must be one of: active, inactive, deleted")
		}
	}

	if input.Role != "" {
		if matched := regexp.MustCompile(RoleRegex).MatchString(input.Role); !matched {
			errorMessages = append(errorMessages, "role must be one of: user, admin, moderator, merchant")
		}
	}

	if len(errorMessages) > 0 {
		return lib.BadRequestError(strings.Join(errorMessages, "; "))
	}

	return nil
}
