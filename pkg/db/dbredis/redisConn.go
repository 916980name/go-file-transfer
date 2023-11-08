package dbredis

import (
	"context"
	"file-transfer/pkg/log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	REDIS_LOGIN_SHARE_KEY_PREFIX = "ls-"

	client     *redis.Client
	clientOnce sync.Once
)

func initClient(ctx context.Context) error {
	var err error
	connOpt := ReadRedisOptions()
	opts, err := redis.ParseURL(connOpt.ConnectionString)
	if err != nil {
		return err
	}
	client = redis.NewClient(opts)
	return nil
}

func GetClient(ctx context.Context) *redis.Client {
	clientOnce.Do(func() {
		if client == nil {
			cancel := func() {}
			if ctx == nil {
				ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
				defer cancel()
			}
			initClient(ctx)
		}
	})
	return client
}

func CloseClient(ctx context.Context) error {
	if client != nil {
		err := client.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func RetryConnect(ctx context.Context, retryDelay time.Duration) error {
	var err error
	for {
		err = initClient(ctx)
		if err == nil {
			return nil
		}
		log.Errorw("Failed to connect to redis (" + err.Error() + "). Retrying in " + retryDelay.String())
		time.Sleep(retryDelay)
	}
}
