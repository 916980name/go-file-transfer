package service

import (
	"context"
	"crypto/subtle"
	"file-transfer/internal/file-transfer/repo"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/common"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/log"
	"file-transfer/pkg/model"
	"file-transfer/pkg/util"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	SHARE_LINK_EXPIRE time.Duration = 5 * time.Minute
)

type UserService interface {
	CreateUser(ctx context.Context, username string) (*model.UserInfo, error)
	Login(ctx context.Context, request v1.UserLoginRequest) (*model.UserInfo, error)
	CreateLoginUrl(ctx context.Context, userId string) (string, error)
	LoginByLoginUrl(ctx context.Context, key string) (*model.UserInfo, error)
}

type userService struct {
	userRepo    repo.UserRepo
	redisClient *redis.Client
	shareServ   ShareService
}

var _ UserService = (*userService)(nil)

const DEFAULT_PASSWORD_LENGTH = 32

func NewUserService(repo repo.UserRepo, rClient *redis.Client, shareServ ShareService) UserService {
	return &userService{userRepo: repo, redisClient: rClient, shareServ: shareServ}
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
		log.C(ctx).Warnw("authentication failed", err)
		return nil, &errno.Errno{
			Message: "authentication failed",
		}
	}
	if subtle.ConstantTimeCompare([]byte(user.Password), []byte(request.Password)) != 1 {
		log.C(ctx).Warnw("authentication failed, bad password")
		return nil, &errno.Errno{
			Message: "authentication failed",
		}
	}

	return user, nil
}

func (s *userService) CreateLoginUrl(ctx context.Context, userId string) (string, error) {
	return s.shareServ.CreateShareUrl(ctx, common.SHARE_TYPE_LOGIN, userId, SHARE_LINK_EXPIRE)
}

func (s *userService) LoginByLoginUrl(ctx context.Context, key string) (*model.UserInfo, error) {
	userId, err := s.shareServ.CheckShareUrl(ctx, common.SHARE_TYPE_LOGIN, key)
	if err != nil {
		return nil, errno.ErrInvalidParameter
	}
	user, err := s.userRepo.FindById(ctx, userId)
	if err != nil {
		log.C(ctx).Warnw("user find failed "+userId, err)
		return nil, errno.ErrInvalidParameter
	}
	return user, nil
}
