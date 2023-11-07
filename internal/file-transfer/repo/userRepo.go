package repo

import (
	"context"
	"file-transfer/pkg/db/dbmongo"
	"file-transfer/pkg/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.UserInfo) (string, error)
	FindById(ctx context.Context, id string) (*model.UserInfo, error)
	FindByUsername(ctx context.Context, username string) (*model.UserInfo, error)
}

type userRepoImpl struct {
	db *mongo.Client
}

var _ UserRepo = (*userRepoImpl)(nil)

func newUserRepo(db *mongo.Client) *userRepoImpl {
	return &userRepoImpl{db}
}

func NewUserRepo(db *mongo.Client) UserRepo {
	return &userRepoImpl{db}
}

func (u *userRepoImpl) Create(ctx context.Context, user *model.UserInfo) (string, error) {
	collection := u.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_USER)
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return "", err
}

func (u *userRepoImpl) FindById(ctx context.Context, id string) (*model.UserInfo, error) {
	collection := u.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_USER)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": bson.M{"$eq": objID}}
	user := &model.UserInfo{}
	if err := collection.FindOne(ctx, filter).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepoImpl) FindByUsername(ctx context.Context, username string) (*model.UserInfo, error) {
	collection := u.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_USER)

	filter := bson.M{"username": bson.M{"$eq": username}}
	user := &model.UserInfo{}
	if err := collection.FindOne(ctx, filter).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}
