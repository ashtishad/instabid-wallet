package service

import (
	"context"
	"log/slog"

	"github.com/ashtishad/instabid-wallet/lib"
	"github.com/ashtishad/instabid-wallet/user-api/internal/domain"
	"github.com/ashtishad/instabid-wallet/user-api/pkg/hashpass"
)

type UserService interface {
	NewUser(ctx context.Context, req domain.NewUserReqDTO) (*domain.UserRespDTO, lib.APIError)
}

type DefaultUserService struct {
	repo domain.UserRepository
	l    *slog.Logger
}

func NewUserService(repo domain.UserRepository) *DefaultUserService {
	return &DefaultUserService{repo: repo}
}

func (s *DefaultUserService) NewUser(ctx context.Context, req domain.NewUserReqDTO) (*domain.UserRespDTO, lib.APIError) {
	hashedPass, err := hashpass.Generate(ctx, req.Password, s.l)
	if err != nil {
		return nil, err
	}

	u := domain.User{
		UserName:   req.UserName,
		Email:      req.Email,
		Status:     req.Status,
		Role:       req.Role,
		HashedPass: hashedPass,
	}

	user, err := s.repo.Insert(ctx, u)
	if err != nil {
		return nil, err
	}

	userDTO := domain.UserRespDTO{
		UserID:    user.UserID,
		UserName:  user.UserName,
		Email:     user.Email,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &userDTO, err
}
