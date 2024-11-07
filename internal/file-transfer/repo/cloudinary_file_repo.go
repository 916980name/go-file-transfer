package repo

import (
	"context"
	"file-transfer/pkg/model"
	"file-transfer/pkg/third"

	"go.mongodb.org/mongo-driver/mongo"
)

func (f *fileRepoImpl) CloudinaryNewFile(ctx context.Context, m *model.CloudinaryFile) (*mongo.InsertOneResult, error) {
	c := f.db.Database(third.CLOUDINARY_MONGODB).Collection(third.CLOUDINARY_MONGOCOLLECTION)
	return c.InsertOne(ctx, m)
}
