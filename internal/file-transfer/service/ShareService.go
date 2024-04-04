package service

import (
	"context"
	"file-transfer/pkg/common"
	"file-transfer/pkg/db/dbredis"
	"file-transfer/pkg/encrypt/aesencrypt"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/log"
	"file-transfer/pkg/util"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var shareTypePathMap = make(map[common.ShareKey]func(encodeKStr string) string)
var shareTypePrefixMap = make(map[common.ShareKey]string)

func init() {
	shareTypePathMap[common.SHARE_TYPE_LOGIN] = getSharePathLogin()
	shareTypePathMap[common.SHARE_TYPE_MESSAGE] = getSharePathMsg()
	shareTypePathMap[common.SHARE_TYPE_FILE] = getSharePathFile()

	shareTypePrefixMap[common.SHARE_TYPE_LOGIN] = dbredis.REDIS_LOGIN_SHARE_KEY_PREFIX
	shareTypePrefixMap[common.SHARE_TYPE_MESSAGE] = dbredis.REDIS_MESSAGE_SHARE_KEY_PREFIX
	shareTypePrefixMap[common.SHARE_TYPE_FILE] = dbredis.REDIS_FILE_SHARE_KEY_PREFIX
}

func getSharePathLogin() func(encodeKStr string) string {
	return func(encodeKStr string) string {
		return "/" + common.LOGIN_SHARE_PATH + "/" + encodeKStr
	}
}

func getSharePathMsg() func(encodeKStr string) string {
	return func(encodeKStr string) string {
		return "/" + common.MESSAGE_SHARE_PATH + "/" + encodeKStr
	}
}

func getSharePathFile() func(encodeKStr string) string {
	return func(encodeKStr string) string {
		return "/" + common.FILE_SHARE_PATH + "/" + encodeKStr
	}
}

type ShareService interface {
	CreateShareUrl(ctx context.Context, shareType common.ShareKey, value string, expire time.Duration) (string, error)
	CreateShareUrlWithTimes(ctx context.Context, shareType common.ShareKey, value string, expire time.Duration, times int8) (string, error)
	CheckShareUrl(ctx context.Context, shareType common.ShareKey, key string, expire time.Duration) (string, error)
	ConsumeShareUrl(ctx context.Context, shareType common.ShareKey, key string, expire time.Duration) (string, error)
}

type shareService struct {
	redisClient *redis.Client
}

var _ ShareService = (*shareService)(nil)

func NewShareService(rClient *redis.Client) ShareService {
	return &shareService{redisClient: rClient}
}

func genEncodeString(ctx context.Context) (string, error) {
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
	return encodeKStr, nil
}

func (s *shareService) CreateShareUrl(ctx context.Context, shareType common.ShareKey, value string, expire time.Duration) (string, error) {
	if len(value) < 1 {
		return "", errno.ErrInvalidParameter
	}
	kStr, err := genEncodeString(ctx)
	if err != nil {
		return "", err
	}
	// TODO: there could create many link for one user
	// set origin kStr
	bc := s.redisClient.Set(ctx, shareTypePrefixMap[shareType]+kStr, value, expire)
	if bc.Err() != nil {
		log.C(ctx).Warnw(bc.Err().Error())
		return "", errno.InternalServerError
	}
	// return encodeKStr
	return shareTypePathMap[shareType](kStr), nil
}

func (s *shareService) CreateShareUrlWithTimes(ctx context.Context, shareType common.ShareKey, value string, expire time.Duration, times int8) (string, error) {
	if len(value) < 1 {
		return "", errno.ErrInvalidParameter
	}
	kStr, err := genEncodeString(ctx)
	if err != nil {
		return "", err
	}
	bc := s.redisClient.Set(ctx, shareTypePrefixMap[shareType]+kStr, value, expire)
	if bc.Err() != nil {
		log.C(ctx).Warnw(bc.Err().Error())
		return "", errno.InternalServerError
	}
	bc = s.redisClient.Set(ctx, "count-"+kStr, times, expire)
	if bc.Err() != nil {
		log.C(ctx).Warnw(bc.Err().Error())
		return "", errno.InternalServerError
	}
	// return encodeKStr
	return shareTypePathMap[shareType](kStr), nil
}

func checkKeyExpire(ctx context.Context, key string, expire time.Duration) error {
	if len(key) < 1 {
		return errno.ErrInvalidParameter
	}
	key = util.Base64urlDecode(key)
	// check the 'key' valid, time not expired
	timeStr, err := aesencrypt.GetAESDecrypted(key)
	if err != nil {
		log.Warnw("share link decrypt fail", "error", err)
		return errno.ErrInvalidParameter
	}
	log.C(ctx).Debugw("timeStr: " + timeStr)
	// Convert timestamp string to int64
	timestamp, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		log.Warnw("Invalid timestamp format", "error", err)
		return errno.ErrInvalidParameter
	}

	expireTime := time.Unix(0, timestamp*int64(time.Millisecond)).Add(expire)
	if expireTime.Before(time.Now()) {
		log.Infow("share link expired at: " + expireTime.Format(time.RFC3339))
		return errno.ErrInvalidParameter
	}
	return nil
}

func (s *shareService) CheckShareUrl(ctx context.Context, shareType common.ShareKey, key string, expire time.Duration) (string, error) {
	err := checkKeyExpire(ctx, key, expire)
	if err != nil {
		return "", err
	}

	// then check redis
	sc := s.redisClient.Get(ctx, shareTypePrefixMap[shareType]+key)
	if sc.Err() != nil {
		log.C(ctx).Warnw(sc.Err().Error())
		return "", errno.ErrInvalidParameter
	}
	value := sc.Val()
	if value == "" {
		log.C(ctx).Infow("[" + fmt.Sprint(shareType) + "] share link not match: " + key)
		return "", errno.ErrInvalidParameter
	}
	ic := s.redisClient.Del(ctx, shareTypePrefixMap[shareType]+key)
	if ic.Err() != nil {
		log.C(ctx).Warnw(ic.Err().Error())
		return "", errno.InternalServerError
	}

	return value, nil
}

func (s *shareService) ConsumeShareUrl(ctx context.Context, shareType common.ShareKey, key string, expire time.Duration) (string, error) {
	err := checkKeyExpire(ctx, key, expire)
	if err != nil {
		return "", err
	}
	sc := s.redisClient.Get(ctx, shareTypePrefixMap[shareType]+key)
	if sc.Err() != nil {
		log.C(ctx).Warnw(sc.Err().Error())
		return "", errno.ErrInvalidParameter
	}
	value := sc.Val()
	if value == "" {
		log.C(ctx).Infow("[" + fmt.Sprint(shareType) + "] share link not match: " + key)
		return "", errno.ErrInvalidParameter
	}
	sc = s.redisClient.Get(ctx, "count-"+key)
	// if there is no count here, it's only expired by duration
	if sc.Err() != nil {
		log.C(ctx).Debugw(sc.Err().Error())
		return value, nil
	}

	ic := s.redisClient.Decr(ctx, "count-"+key)
	// if decr count error, something wrong
	if ic.Err() != nil {
		log.C(ctx).Warnw(ic.Err().Error())
		return value, nil
	}
	log.C(ctx).Debugw("access key: "+key, "count", ic.Val())
	if ic.Val() <= 0 {
		// do clean
		log.C(ctx).Debugw("Clean cache key: " + key)
		s.redisClient.Del(ctx, shareTypePrefixMap[shareType]+key)
		s.redisClient.Del(ctx, "count-"+key)
	}
	return value, nil
}
