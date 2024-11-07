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

type FileRepo interface {
	InsertFileMeta(ctx context.Context, m *model.FileMeta) (*mongo.InsertOneResult, error)
	InsertUserFile(ctx context.Context, m *model.UserFile) (*mongo.InsertOneResult, error)
	FindOneBySha(ctx context.Context, sha string) (*model.FileMeta, error)
	FindByMetaId(ctx context.Context, ids []string) ([]model.FileMeta, error)

	FindOneByNameAndUser(ctx context.Context, name string, userId string) (*model.UserFile, error)
	QueryUserFile(ctx context.Context, condition *v1.UserFileQuery) ([]model.UserFile, error)
	QueryUserFileById(ctx context.Context, userFileId string) (*model.UserFile, error)
	DeleteUserFile(ctx context.Context, userFileId string) (*model.UserFile, error)
	DeleteMetaFile(ctx context.Context, metaFileId string) (*model.FileMeta, error)

	CloudinaryNewFile(ctx context.Context, m *model.CloudinaryFile) (*mongo.InsertOneResult, error)
}

type fileRepoImpl struct {
	db *mongo.Client
}

var _ FileRepo = (*fileRepoImpl)(nil)

func NewFileRepo(db *mongo.Client) FileRepo {
	return &fileRepoImpl{db}
}

func (f *fileRepoImpl) InsertFileMeta(ctx context.Context, m *model.FileMeta) (*mongo.InsertOneResult, error) {
	c := f.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_FILE_META)
	return c.InsertOne(ctx, m)
}
func (f *fileRepoImpl) InsertUserFile(ctx context.Context, m *model.UserFile) (*mongo.InsertOneResult, error) {
	c := f.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_USER_FILE)
	return c.InsertOne(ctx, m)
}

func (f *fileRepoImpl) FindOneBySha(ctx context.Context, sha string) (*model.FileMeta, error) {
	c := f.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_FILE_META)
	log.C(ctx).Debugw("FindOneBySha", "sha", sha)
	filter := bson.M{"sha": sha}
	var result model.FileMeta
	err := c.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (f *fileRepoImpl) FindOneByNameAndUser(ctx context.Context, name string, userId string) (*model.UserFile, error) {
	c := f.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_USER_FILE)
	log.C(ctx).Debugw("FindOneByNameAndUser", "name", name, "userId", userId)
	filter := bson.M{"name": name, "userId": userId}
	var result model.UserFile
	err := c.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (f *fileRepoImpl) QueryUserFile(ctx context.Context, condition *v1.UserFileQuery) ([]model.UserFile, error) {
	collection := f.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_USER_FILE)
	var skip int64 = 0
	if condition.PageNum-1 > 0 {
		skip = (condition.PageNum - 1) * condition.PageSize
	}
	var filter bson.M
	if condition.UserId != "" {
		filter = bson.M{
			"userId": bson.M{"$regex": condition.UserId, "$options": "i"},
		}
	}
	options := options.Find().
		SetSkip(skip).
		SetLimit(condition.PageSize).
		SetSort(bson.M{"createdAt": -1})

	cur, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		log.Fatalw(err.Error())
		return nil, err
	}
	defer cur.Close(ctx)
	return iterateUserFileResult(ctx, cur)
}

func (f *fileRepoImpl) QueryUserFileById(ctx context.Context, userFileId string) (*model.UserFile, error) {
	c := f.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_USER_FILE)
	objID, err := primitive.ObjectIDFromHex(userFileId)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	var result model.UserFile
	err = c.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func iterateUserFileResult(ctx context.Context, cur *mongo.Cursor) ([]model.UserFile, error) {
	arr := make([]model.UserFile, 0)
	for cur.Next(ctx) {
		var result model.UserFile
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

func (f *fileRepoImpl) FindByMetaId(ctx context.Context, ids []string) ([]model.FileMeta, error) {
	objIds := make([]primitive.ObjectID, len(ids))
	for _, id := range ids {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objIds = append(objIds, objId)
	}
	c := f.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_FILE_META)
	filter := bson.M{"_id": bson.M{"$in": objIds}}
	// Find the documents matching the filter
	cursor, err := c.Find(ctx, filter)
	if err != nil {
		log.Fatalw(err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	return iterateFileMetaResult(ctx, cursor)
}

func iterateFileMetaResult(ctx context.Context, cur *mongo.Cursor) ([]model.FileMeta, error) {
	arr := make([]model.FileMeta, 0)
	for cur.Next(ctx) {
		var result model.FileMeta
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

func (f *fileRepoImpl) DeleteUserFile(ctx context.Context, userFileId string) (*model.UserFile, error) {
	userC := f.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_USER_FILE)
	data, err := f.QueryUserFileById(ctx, userFileId)
	if err != nil {
		return nil, err
	}
	objID, err := primitive.ObjectIDFromHex(userFileId)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": bson.M{"$eq": objID}}
	_, err = userC.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (f *fileRepoImpl) DeleteMetaFile(ctx context.Context, metaId string) (*model.FileMeta, error) {
	filter := bson.M{"metaId": metaId}
	// check if user hold the file
	userC := f.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_USER_FILE)
	cursor, err := userC.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var metaData *model.FileMeta
	// Check if any records exist
	if !cursor.Next(ctx) {
		// No records matching the criteria
		metaC := f.db.Database(dbmongo.MONGO_DATABASE).Collection(dbmongo.COLL_FILE_META)
		objID, err := primitive.ObjectIDFromHex(metaId)
		if err != nil {
			return nil, err
		}
		filter = bson.M{"_id": bson.M{"$eq": objID}}
		err = metaC.FindOne(ctx, filter).Decode(&metaData)
		if err != nil {
			return nil, err
		}
		_, err = metaC.DeleteOne(ctx, filter)
		if err != nil {
			return nil, err
		}
	} // Close the cursor
	return metaData, nil
}
