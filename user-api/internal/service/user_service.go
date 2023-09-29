package service

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"

	"github.com/ashtishad/instabid-wallet/lib"
	"github.com/ashtishad/instabid-wallet/user-api/internal/domain"
	"github.com/ashtishad/instabid-wallet/user-api/pkg/hashpass"
	"github.com/ashtishad/instabid-wallet/user-api/pkg/utils"
)

type UserService interface {
	NewUser(ctx context.Context, req domain.NewUserReqDTO) (*domain.UserRespDTO, lib.APIError)
	NewProfile(ctx context.Context, uuid string, req domain.NewProfileReqDTO) (*domain.ProfileRespDTO, lib.APIError)
}

type DefaultUserService struct {
	repo domain.UserRepository
	l    *slog.Logger
}

func NewUserService(repo domain.UserRepository, l *slog.Logger) *DefaultUserService {
	return &DefaultUserService{repo: repo, l: l}
}

func (s *DefaultUserService) NewUser(ctx context.Context,
	req domain.NewUserReqDTO) (*domain.UserRespDTO, lib.APIError) {
	if apiErr := utils.ValidateCreateUserInput(req); apiErr != nil {
		return nil, apiErr
	}

	if req.Status == "" {
		req.Status = utils.UserStatusActive
	}

	if req.Role == "" {
		req.Role = utils.UserRoleUser
	}

	hashedPass, err := hashpass.Generate(ctx, req.Password, s.l)
	if err != nil {
		return nil, err
	}

	u := domain.User{
		UserName:   strings.ToLower(req.UserName),
		Email:      strings.ToLower(req.Email),
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

func (s *DefaultUserService) NewProfile(ctx context.Context, uuid string,
	req domain.NewProfileReqDTO) (*domain.ProfileRespDTO, lib.APIError) {
	if apiErr := utils.ValidateCreateProfileInput(req); apiErr != nil {
		return nil, apiErr
	}

	up := domain.Profile{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Gender:    req.Gender,
	}

	if req.Address != "" {
		up.Address = sql.NullString{
			String: req.Address,
			Valid:  true,
		}
	}

	res, apiErr := s.repo.InsertProfile(ctx, uuid, up)
	if apiErr != nil {
		return nil, apiErr
	}

	resDto := domain.ProfileRespDTO{
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Gender:    res.Gender,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}

	if res.Address.Valid {
		resDto.Address = res.Address.String
	}

	return &resDto, nil
}
