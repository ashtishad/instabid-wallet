package lib

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	usernameRegex = `^[a-zA-Z0-9]{7,64}$`
	emailRegex    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

// ValidateEmail checks if input Must match the specified regex pattern EmailRegex.
func ValidateEmail(email string) error {
	if matched := regexp.MustCompile(emailRegex).MatchString(email); !matched {
		return fmt.Errorf("invalid email, you entered %s", email)
	}

	return nil
}

// ValidatePassword checks password must be at least 8 characters long and no more than 32 characters
func ValidatePassword(password string) error {
	if len(password) < 8 || len(password) > 32 {
		return errors.New("password must be at least 8 characters long and no more than 32 characters")
	}

	return nil
}

// ValidateUserName validates a username based on the following criteria:
// 1. The username must be between 7 and 64 characters long.
// 2. The username can only contain alphanumeric characters.
// 3. The username must not contain any spaces.
//
// It returns an error if the username does not meet these criteria.
func ValidateUserName(userName string) error {
	if ok := regexp.MustCompile(usernameRegex).MatchString(userName); !ok {
		return errors.New("invalid username: must be 7-64 alphanumeric characters with no spaces")
	}

	return nil
}
