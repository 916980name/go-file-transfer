package dbmongo

import (
	"context"
	"errors"
	"file-transfer/pkg/log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MONGO_DATABASE = "filetransfer"
	COLL_MESSAGE   = "message"
	COLL_USER      = "user"
	COLL_FILE_META = "filemeta"
	COLL_USER_FILE = "userfile"

	client     *mongo.Client
	clientOnce sync.Once
)

func initClient(ctx context.Context) error {
	var err error
	connOpt := ReadMongoOptions()
	clientOpt := options.Client()
	clientOpt.ApplyURI(connOpt.ConnectionString)
	clientOpt.SetMaxPoolSize(connOpt.MaxPoolSize)
	clientOpt.SetMinPoolSize(connOpt.MinPoolSize)

	client, err = mongo.Connect(ctx, clientOpt)
	if err != nil {
		log.Fatalw(err.Error())
	}

	// Test the connection to ensure that it is working
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalw(err.Error())
	}
	return err
}

func GetClient(ctx context.Context) *mongo.Client {
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
		err := client.Disconnect(ctx)
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
		log.Errorw("Failed to connect to MongoDB (" + err.Error() + "). Retrying in " + retryDelay.String())
		time.Sleep(retryDelay)
	}
}

func IsConnectionError(err error) bool {
	var mongoErr mongo.CommandError
	if errors.As(err, &mongoErr) {
		// Check for network error
		return mongoErr.HasErrorLabel("TransientTransactionError")
	}
	return false
}
