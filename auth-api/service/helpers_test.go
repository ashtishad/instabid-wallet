package service

import (
	"testing"

	"github.com/ashtishad/instabid-wallet/auth-api/domain"
)

func TestValidateLoginRequest(t *testing.T) {
	testCases := []struct {
		req      domain.LoginRequest
		hasError bool
		errMsg   string
	}{
		{domain.LoginRequest{}, true, "either username or email, along with a password, must be provided"},
		{domain.LoginRequest{Username: "testUser", Password: "password"}, false, ""},
		{domain.LoginRequest{Username: "testUser", Email: "test@email.com", Password: "password"}, true, "user can't sign in with both username and email"},
		{domain.LoginRequest{Email: "test@email.com", Password: "password"}, false, ""},
		{domain.LoginRequest{Username: "invalid user", Password: "password"}, true, "invalid username: must be 7-64 alphanumeric characters with no spaces"},
		{domain.LoginRequest{Email: "invalid-email", Password: "password"}, true, "invalid email, you entered invalid-email"},
		{domain.LoginRequest{Username: "testUser", Password: "short"}, true, "password must be at least 8 characters long and no more than 32 characters"},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			err := validateLoginRequest(tc.req)
			if tc.hasError {
				if err == nil {
					t.Errorf("expected an error but got none")
				} else if err.Error() != tc.errMsg {
					t.Errorf("expected error message %q, got %q", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, but got %q", err.Error())
				}
			}
		})
	}
}
