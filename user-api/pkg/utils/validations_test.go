package utils

import (
	"strings"
	"testing"

	"github.com/ashtishad/instabid-wallet/user-api/internal/domain"
)

// TestValidateCreateUserInput tests ValidateCreateUserInput function.
func TestValidateCreateUserInput(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.NewUserReqDTO
		wantErr bool
		errText string
	}{
		{
			name: "Valid input",
			input: domain.NewUserReqDTO{
				UserName: "testUser",
				Password: "password123",
				Email:    "email@test.com",
				Status:   "active",
				Role:     "user",
			},
			wantErr: false,
		},
		{
			name: "Invalid email",
			input: domain.NewUserReqDTO{
				UserName: "testUser",
				Password: "password123",
				Email:    "invalid-email",
				Status:   "active",
				Role:     "user",
			},
			wantErr: true,
			errText: "invalid email, you entered invalid-email",
		},
		{
			name: "Empty username",
			input: domain.NewUserReqDTO{
				UserName: "",
				Password: "password123",
				Email:    "email@test.com",
				Status:   "active",
				Role:     "user",
			},
			wantErr: true,
			errText: "username must be between 7 and 64 characters long",
		},
		{
			name: "Short password",
			input: domain.NewUserReqDTO{
				UserName: "testUser",
				Password: "short",
				Email:    "email@test.com",
				Status:   "active",
				Role:     "user",
			},
			wantErr: true,
			errText: "password must be at least 8 characters long and no more than 32 characters",
		},
		{
			name: "Invalid status",
			input: domain.NewUserReqDTO{
				UserName: "testUser",
				Password: "password123",
				Email:    "email@test.com",
				Status:   "unknown",
				Role:     "user",
			},
			wantErr: true,
			errText: "status must be one of: active, inactive, deleted",
		},
		{
			name: "Invalid role",
			input: domain.NewUserReqDTO{
				UserName: "testUser",
				Password: "password123",
				Email:    "email@test.com",
				Status:   "active",
				Role:     "alien",
			},
			wantErr: true,
			errText: "role must be one of: user, admin, moderator, merchant",
		},
		{
			name: "Multiple errors",
			input: domain.NewUserReqDTO{
				UserName: "",
				Password: "short",
				Email:    "invalid-email",
				Status:   "unknown",
				Role:     "alien",
			},
			wantErr: true,
			errText: "invalid email, you entered invalid-email\npassword must be at least 8 characters long and no more than 32 characters\nusername must be between 7 and 64 characters long\nstatus must be one of: active, inactive, deleted\nrole must be one of: user, admin, moderator, merchant",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ValidateCreateUserInput(tt.input)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("ValidateCreateUserInput() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}

			if gotErr != nil && gotErr.Error() != tt.errText {
				t.Errorf("ValidateCreateUserInput() got error text = %v, want %v", gotErr.Error(), tt.errText)
			}
		})
	}
}

// TestValidateCreateProfileInput tests ValidateCreateProfileInput function.
func TestValidateCreateProfileInput(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.NewProfileReqDTO
		wantErr bool
		errText string
	}{
		{
			name: "Valid input",
			input: domain.NewProfileReqDTO{
				FirstName: "John",
				LastName:  "Doe",
				Gender:    "male",
				Address:   "1234 Elm St",
			},
			wantErr: false,
		},
		{
			name: "Invalid first name",
			input: domain.NewProfileReqDTO{
				FirstName: "J@hn",
				LastName:  "Doe Susan",
				Gender:    "male",
			},
			wantErr: true,
			errText: "first name must be alphabetic and between 1 and 64 characters long",
		},
		{
			name: "Invalid last name",
			input: domain.NewProfileReqDTO{
				FirstName: "John",
				LastName:  "D@e Trixy",
				Gender:    "male",
			},
			wantErr: true,
			errText: "last name must be alphabetic, may contain spaces, and be between 1 and 128 characters long",
		},
		{
			name: "Invalid gender",
			input: domain.NewProfileReqDTO{
				FirstName: "John",
				LastName:  "Doe",
				Gender:    "alien",
			},
			wantErr: true,
			errText: "gender must be one of: male, female, other",
		},
		{
			name: "Valid Input and Empty address",
			input: domain.NewProfileReqDTO{
				FirstName: "John",
				LastName:  "Doe",
				Gender:    "male",
				Address:   "",
			},
			wantErr: false,
		},
		{
			name: "Address too long",
			input: domain.NewProfileReqDTO{
				FirstName: "John",
				LastName:  "Doe",
				Gender:    "male",
				Address:   strings.Repeat("a", 257),
			},
			wantErr: true,
			errText: "address cannot exceed 256 characters",
		},
		{
			name: "Multiple errors",
			input: domain.NewProfileReqDTO{
				FirstName: "John ",
				LastName:  "D@e",
				Gender:    "alien",
			},
			wantErr: true,
			errText: "first name must be alphabetic and between 1 and 64 characters long\nlast name must be alphabetic, may contain spaces, and be between 1 and 128 characters long\ngender must be one of: male, female, other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ValidateCreateProfileInput(tt.input)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("ValidateCreateProfileInput() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}

			if gotErr != nil && gotErr.Error() != tt.errText {
				t.Errorf("ValidateCreateProfileInput() got error text = %v, want %v", gotErr.Error(), tt.errText)
			}
		})
	}
}
