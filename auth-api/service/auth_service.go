package service

import (
	"context"
	"log/slog"

	"github.com/ashtishad/instabid-wallet/auth-api/domain"
	"github.com/ashtishad/instabid-wallet/lib"
)

type AuthService interface {
	Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, lib.APIError)
}

type DefaultAuthService struct {
	repo domain.AuthRepository
	l    *slog.Logger
}

func NewAuthService(repo domain.AuthRepository, l *slog.Logger) DefaultAuthService {
	return DefaultAuthService{repo: repo, l: l}
}

func (s DefaultAuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, lib.APIError) {
	var apiErr lib.APIError
	var login *domain.Login

	if apiErr = validateLoginRequest(req); apiErr != nil {
		return nil, apiErr
	}

	login, apiErr = s.repo.FindByCredential(ctx, req)
	if apiErr != nil {
		return nil, apiErr
	}

	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims, s.l)

	var accessToken string

	if accessToken, apiErr = authToken.NewAccessToken(); apiErr != nil {
		return nil, apiErr
	}

	response := domain.LoginResponse{
		AccessToken: accessToken,
		Login: domain.Login{
			UserID:   login.UserID,
			Username: login.Username,
			Email:    login.Email,
			Role:     login.Role,
			Status:   login.Status,
		},
	}

	return &response, nil
}
