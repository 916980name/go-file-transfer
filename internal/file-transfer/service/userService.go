package service

import (
	"context"
	"crypto/subtle"
	"file-transfer/internal/file-transfer/repo"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/model"
	"file-transfer/pkg/util"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, username string) (*model.UserInfo, error)
	Login(ctx context.Context, request v1.UserLoginRequest) (*model.UserInfo, error)
}

type userService struct {
	userRepo repo.UserRepo
}

var _ UserService = (*userService)(nil)

const DEFAULT_PASSWORD_LENGTH = 32

func NewUserService(repo repo.UserRepo) UserService {
	return &userService{userRepo: repo}
}

func (s *userService) CreateUser(ctx context.Context, username string) (*model.UserInfo, error) {
	password, ep := util.GenerateRandomString(DEFAULT_PASSWORD_LENGTH)
	if ep != nil {
		return nil, ep
	}
	user := &model.UserInfo{
		Username:  username,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	user, error := s.userRepo.FindById(ctx, id)
	return user, error
}

func (s *userService) Login(ctx context.Context, request v1.UserLoginRequest) (*model.UserInfo, error) {
	if len(request.Username) == 0 || len(request.Password) == 0 {
		return nil, &errno.Errno{
			Message: "request field can't empty",
		}
	}
	user, err := s.userRepo.FindByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}
	if subtle.ConstantTimeCompare([]byte(user.Password), []byte(request.Password)) != 1 {
		return nil, &errno.Errno{
			Message: "authentication failed",
		}
	}

	return user, nil
}
