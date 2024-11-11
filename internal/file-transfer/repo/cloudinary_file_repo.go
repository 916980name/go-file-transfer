package repo

import (
	"context"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/model"
	"file-transfer/pkg/third"

	"go.mongodb.org/mongo-driver/mongo"
)

func (f *fileRepoImpl) CloudinaryNewFile(ctx context.Context, m *model.CloudinaryFile) (*mongo.InsertOneResult, error) {
	c := f.db.Database(third.CLOUDINARY_MONGODB).Collection(third.CLOUDINARY_MONGOCOLLECTION)
	return c.InsertOne(ctx, m)
}

func (f *fileRepoImpl) CloudinaryQueryAllFile(ctx context.Context, condition *v1.CloudinaryFileReq) ([]model.CloudinaryFile, error) {

	return nil, nil
}

func (f *fileRepoImpl) CloudinaryQueryFileById(ctx context.Context, assetId string) (*model.CloudinaryFile, error) {

	return nil, nil
}
