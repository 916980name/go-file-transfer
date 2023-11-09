package service

import (
	"context"
	"crypto/subtle"
	"file-transfer/internal/file-transfer/repo"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/common"
	"file-transfer/pkg/db/dbredis"
	"file-transfer/pkg/encrypt/aesencrypt"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/log"
	"file-transfer/pkg/model"
	"file-transfer/pkg/util"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
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
	nowMilli := fmt.Sprint(time.Now().UnixMilli())
	log.C(ctx).Debugw("now: " + nowMilli)
	kStr, err := aesencrypt.GetAESEncrypted(nowMilli)
	if err != nil {
		log.C(ctx).Warnw(err.Error())
		return "", errno.InternalServerError
	}
	// base64 string could contains "+""/"
	log.C(ctx).Debugw("kStr: " + kStr)
	encodeKStr := util.Base64urlEncode(kStr)
	log.C(ctx).Debugw("encode kStr: " + encodeKStr)
	// TODO: there could create many link for one user
	// set origin kStr
	bc := s.redisClient.Set(ctx, dbredis.REDIS_LOGIN_SHARE_KEY_PREFIX+kStr, userId, SHARE_LINK_EXPIRE)
	if bc.Err() != nil {
		log.C(ctx).Warnw(bc.Err().Error())
		return "", errno.InternalServerError
	}
	// return encodeKStr
	return viper.GetString(common.VIPER_HOST_URL) + "/" + common.LOGIN_SHARE_PATH + "/" + encodeKStr, nil
}

func (s *userService) LoginByLoginUrl(ctx context.Context, key string) (*model.UserInfo, error) {
	if len(key) < 1 {
		return nil, errno.ErrInvalidParameter
	}
	key = util.Base64urlDecode(key)
	// check the 'key' valid, time not expired
	timeStr, err := aesencrypt.GetAESDecrypted(key)
	if err != nil {
		log.Warnw("share link decrypt fail", "error", err)
		return nil, errno.ErrInvalidParameter
	}
	shareTime, err := time.ParseDuration(timeStr + "ms")
	if err != nil {
		log.Warnw("share link time parse fail", "error", err)
		return nil, errno.ErrInvalidParameter
	}
	expireTime := time.Unix(0, shareTime.Nanoseconds()*int64(time.Millisecond)).Add(SHARE_LINK_EXPIRE)
	if expireTime.Before(time.Now()) {
		log.Infow("share link expired at: " + expireTime.Format(time.RFC3339))
		return nil, errno.ErrInvalidParameter
	}

	// then check redis
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
