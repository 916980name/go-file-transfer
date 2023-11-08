package service

import (
	"context"
	"crypto/subtle"
	"file-transfer/internal/file-transfer/repo"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/common"
	"file-transfer/pkg/db/dbredis"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/log"
	"file-transfer/pkg/model"
	"file-transfer/pkg/util"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
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
}

var _ UserService = (*userService)(nil)

const DEFAULT_PASSWORD_LENGTH = 32

func NewUserService(repo repo.UserRepo, rClient *redis.Client) UserService {
	return &userService{userRepo: repo, redisClient: rClient}
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
	if len(userId) < 1 {
		return "", errno.ErrInvalidParameter
	}
	kStr, _ := util.GenerateRandomString(16)
	bc := s.redisClient.Set(ctx, dbredis.REDIS_LOGIN_SHARE_KEY_PREFIX+kStr, userId, 5*time.Minute)
	if bc.Err() != nil {
		log.C(ctx).Warnw(bc.Err().Error())
		return "", errno.InternalServerError
	}
	return viper.GetString(common.VIPER_HOST_URL) + "/" + common.LOGIN_SHARE_PATH + "/" + kStr, nil
}

func (s *userService) LoginByLoginUrl(ctx context.Context, key string) (*model.UserInfo, error) {
	if len(key) < 1 {
		return nil, errno.ErrInvalidParameter
	}
	sc := s.redisClient.Get(ctx, dbredis.REDIS_LOGIN_SHARE_KEY_PREFIX+key)
	if sc.Err() != nil {
		log.C(ctx).Warnw(sc.Err().Error())
		return nil, errno.ErrInvalidParameter
	}
	userId := sc.Val()
	if userId == "" {
		log.C(ctx).Infow("share link not match: " + key)
		return nil, errno.ErrInvalidParameter
	}
	ic := s.redisClient.Del(ctx, dbredis.REDIS_LOGIN_SHARE_KEY_PREFIX+key)
	if ic.Err() != nil {
		log.C(ctx).Warnw(ic.Err().Error())
		return nil, errno.InternalServerError
	}
	user, err := s.userRepo.FindById(ctx, userId)
	if err != nil {
		log.C(ctx).Warnw("user find failed "+userId, err)
		return nil, errno.ErrInvalidParameter
	}
	return user, nil
}
