package dbmongo

import (
	"context"
	"encoding/json"
	"file-transfer/pkg/config"
	"file-transfer/pkg/log"
	"file-transfer/pkg/model"
	"path/filepath"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestConnQuery(t *testing.T) {
	configFile, err := filepath.Abs("../../../_output/file-transfer.yaml")
	if err != nil {
		t.Logf("file not found: %s", configFile)
	}
	ctx := context.Background()
	config.ReadConfig(configFile)
	client := GetClient(ctx)
	defer CloseClient(ctx)

	collection := client.Database(MONGO_DATABASE).Collection(COLL_MESSAGE)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatalw(err.Error())
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result model.Message
		err := cur.Decode(&result)
		if err != nil {
			log.Fatalw(err.Error())
		}
		s, _ := json.Marshal(result)
		log.Infow(string(s))
		// do something with result....
	}
	if err := cur.Err(); err != nil {
		log.Fatalw(err.Error())
	}
}
