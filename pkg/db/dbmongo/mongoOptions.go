package dbmongo

import "github.com/spf13/viper"

var (
	// should set env: GO_FILE_TRANSFER_MONGO_HOST, because we have setting viper autoenv prefix
	ENV_MONGO_HOST = "MONGO_HOST"
)

type MongoOptions struct {
	ConnectionString string
	MaxPoolSize      uint64
	MinPoolSize      uint64
}

func NewMongoOptions() *MongoOptions {
	return &MongoOptions{
		ConnectionString: "mongodb://localhost:27017",
		MaxPoolSize:      10,
		MinPoolSize:      0,
	}
}

func ReadMongoOptions() *MongoOptions {
	connStr := viper.GetString(ENV_MONGO_HOST)
	if connStr == "" {
		connStr = viper.GetString("db.mongo.connectionString")
	}
	maxPoolSize := viper.GetUint64("db.mongo.maxPoolSize")
	minPoolSize := viper.GetUint64("db.mongo.minPoolSize")

	options := NewMongoOptions()
	if connStr != "" {
		options.ConnectionString = connStr
	}
	if maxPoolSize != 0 {
		options.MaxPoolSize = maxPoolSize
	}
	if minPoolSize != 0 {
		options.MinPoolSize = minPoolSize
	}
	return options
}
