package service

import (
	"context"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/log"
	"file-transfer/pkg/model"
	"file-transfer/pkg/third"
	"mime/multipart"
	"net/http"

	clapi "github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func (f *fileService) CloudinaryUploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, req *v1.CloudinaryFileUpReq) (*model.CloudinaryFile, error) {
	/*
		tempFile, err := os.CreateTemp(TEMP_FILE_DIR, TEMP_FILE_PATTERN)
		log.C(ctx).Debugw("create temp: " + tempFile.Name())
		if err != nil {
			return err
		}
		defer func() {
			log.C(ctx).Debugw("close/remove temp: " + tempFile.Name())
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()

		mReader := io.Reader(file)
		_, err = io.Copy(tempFile, mReader)
		if err != nil {
			return err
		}
	*/

	cld, err := third.Cloudinary_credentials()
	if err != nil {
		log.C(ctx).Errorw("cloudinary init failed", err)
		return nil, &errno.Errno{HTTP: http.StatusInternalServerError, Message: "cloud error"}
	}
	result, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         third.CLOUDINARY_FOLDER,
		PublicID:       header.Filename,
		UniqueFilename: clapi.Bool(false),
		Overwrite:      clapi.Bool(true)})
	if err != nil || result.Error.Message != "" {
		log.C(ctx).Errorw("upload cloudinary failed", err)
		log.C(ctx).Errorw("upload cloudinary failed", result.Error.Message)
		return nil, &errno.Errno{HTTP: http.StatusInternalServerError, Message: "cloud error"}
	}

	m := model.ClCopy(result)
	m.Title = req.Title
	m.Desc = req.Desc
	m.OriginUrl = req.OriginUrl
	_, err = f.fileRepo.CloudinaryNewFile(ctx, m)
	if err != nil {
		log.C(ctx).Errorw("InsertFile Record failed", "err", err, "clResult", result)
		return nil, &errno.Errno{HTTP: http.StatusInternalServerError, Message: "repo error"}
	}

	log.C(ctx).Infow("Upload suc", "clResult", result)
	return m, nil
}
