package repo

import (
	"context"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/db/dbmongo"
	"file-transfer/pkg/log"
	"file-transfer/pkg/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepo interface {
	Query(ctx context.Context, filter *v1.MessageQuery) ([]model.Message, error)
	QueryById(ctx context.Context, mId string) (*model.Message, error)
	Insert(ctx context.Context, m *model.Message) (*mongo.InsertOneResult, error)
	Delete(ctx context.Context, mId string, uId string) (*mongo.DeleteResult, error)
}

type messageRepoImpl struct {
	db *mongo.Client
}

var _ MessageRepo = (*messageRepoImpl)(nil)

func newMessageRepo(db *mongo.Client) *messageRepoImpl {
	return &messageRepoImpl{db}
}

func NewMessageRepo(db *mongo.Client) MessageRepo {
	return &messageRepoImpl{db}
}

func (t *messageRepoImpl) Query(ctx context.Context, condition *v1.MessageQuery) ([]model.Message, error) {
	collection := t.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_MESSAGE)
	// Define the query options for the Find method
	var skip int64 = 0
	if condition.PageNum-1 > 0 {
		skip = (condition.PageNum - 1) * condition.PageSize
	}
	var filter bson.M
	if condition.UserId != "" {
		// Define the query filter for the Find method
		filter = bson.M{
			"userId": bson.M{"$regex": condition.UserId, "$options": "i"},
		}
	}
	options := options.Find().
		SetSkip(skip).
		SetLimit(condition.PageSize).
		SetSort(bson.M{"createdAt": -1})

	// Call the Find method to retrieve the documents that match the query conditions
	cur, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		log.Fatalw(err.Error())
		return nil, err
	}
	defer cur.Close(ctx)
	arr := make([]model.Message, 0)
	for cur.Next(ctx) {
		var result model.Message
		err := cur.Decode(&result)
		if err != nil {
			log.Fatalw(err.Error())
			return nil, err
		}
		arr = append(arr, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatalw(err.Error())
		return nil, err
	}
	return arr, nil
}

func (t *messageRepoImpl) Insert(ctx context.Context, m *model.Message) (*mongo.InsertOneResult, error) {
	c := t.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_MESSAGE)
	return c.InsertOne(ctx, m)
}

func (t *messageRepoImpl) Delete(ctx context.Context, mId string, uId string) (*mongo.DeleteResult, error) {
	c := t.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_MESSAGE)
	objID, err := primitive.ObjectIDFromHex(mId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": bson.M{"$eq": objID}, "userId": bson.M{"$eq": uId}}
	return c.DeleteOne(ctx, filter)
}

func (t *messageRepoImpl) QueryById(ctx context.Context, mId string) (*model.Message, error) {
	c := t.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_MESSAGE)
	objID, err := primitive.ObjectIDFromHex(mId)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	var result model.Message
	err = c.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
