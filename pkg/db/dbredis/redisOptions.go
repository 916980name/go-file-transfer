package dbredis

import "github.com/spf13/viper"

var (
	// should set env: GO_FILE_TRANSFER_REDIS_HOST, because we have setting viper autoenv prefix
	ENV_REDIS_HOST = "REDIS_HOST"
)

type RedisOptions struct {
	ConnectionString string
}

func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		ConnectionString: "redis://localhost:6379",
	}
}

func ReadRedisOptions() *RedisOptions {
	connStr := viper.GetString(ENV_REDIS_HOST)
	if connStr == "" {
		connStr = viper.GetString("db.redis.connectionString")
	}

	options := NewRedisOptions()
	options.ConnectionString = connStr
	return options
}
