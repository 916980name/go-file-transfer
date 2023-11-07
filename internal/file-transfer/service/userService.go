package service

import (
	"context"
	"file-transfer/internal/file-transfer/repo"
	"file-transfer/pkg/model"
	"file-transfer/pkg/util"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, username string) (*model.UserInfo, error)
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
