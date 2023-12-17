package service

import (
	"context"
	"file-transfer/internal/file-transfer/repo"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/log"
	"file-transfer/pkg/model"
	"file-transfer/pkg/util"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TEMP_FILE_DIR     string = "/tmp"
	TEMP_FILE_PATTERN string = "file-transfer-upload-*.tmp"
	// limit reader end with EOF, but don't know is it real end or reach the limit
	MAX_SINGLE_FILE_SIZE int64 = 50*1024*1024 + 1
)

var SAVE_FILE_PATH string

func init() {
	workingPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	log.Infow("Read Working dir: " + workingPath)
	SAVE_FILE_PATH = filepath.Join(workingPath, "uploads")
	err = util.CreateDirectoryIfNotExists(SAVE_FILE_PATH)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	log.Infow("Check Save Upload dir: " + SAVE_FILE_PATH)
}

type FileService interface {
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, userId string) error
	QueryUserFile(ctx context.Context, q *v1.UserFileQuery) ([]v1.FileResponse, error)
}

type fileService struct {
	fileRepo repo.FileRepo
}

var _ FileService = (*fileService)(nil)

func NewFileService(fileRepo repo.FileRepo) FileService {
	return &fileService{fileRepo: fileRepo}
}

func (f *fileService) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, userId string) error {
	if header.Size >= MAX_SINGLE_FILE_SIZE {
		msg := fmt.Sprintf("file size exceed %d", MAX_SINGLE_FILE_SIZE)
		log.C(ctx).Warnw("upload failed, " + msg)
		return &errno.Errno{HTTP: http.StatusBadRequest, Message: msg}
	}

	// check exist file name
	exist, _ := f.fileRepo.FindOneByNameAndUser(ctx, header.Filename, userId)
	if exist != nil {
		msg := fmt.Sprintf("upload exist: %s", header.Filename)
		log.C(ctx).Infow(msg)
		return &errno.Errno{HTTP: http.StatusBadRequest, Message: msg}
	}

	createTime := time.Now()
	fileMeta := &model.FileMeta{
		CreatedAt: createTime,
		Size:      header.Size,
	}

	userFile := &model.UserFile{
		CreatedAt: createTime,
		Name:      header.Filename,
		UserId:    userId,
	}

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

	// Copy file contents to a temporary file while checking the size
	limitedReader := io.LimitReader(file, MAX_SINGLE_FILE_SIZE)
	_, err = io.Copy(tempFile, limitedReader)
	if err != nil {
		// If there was an error while copying, remove the partially written file
		// works in defer
		return err
	}

	// Check the actual file size
	fileInfo, _ := tempFile.Stat()
	fileSize := fileInfo.Size()
	if fileSize >= MAX_SINGLE_FILE_SIZE {
		// If the file size exceeds the limit, remove the file and return an error
		// works in defer
		msg := fmt.Sprintf("file size exceed %d", MAX_SINGLE_FILE_SIZE)
		log.C(ctx).Warnw("upload failed, " + msg)
		return &errno.Errno{HTTP: http.StatusBadRequest, Message: msg}
	}

	sha, err := util.CalculateFileSHA1(tempFile.Name())
	if err != nil {
		log.C(ctx).Errorw("sha1 failed", err)
		return &errno.Errno{HTTP: http.StatusInternalServerError, Message: "save error"}
	}

	result, _ := f.fileRepo.FindOneBySha(ctx, sha)
	if result != nil {
		msg := fmt.Sprintf("upload file exist: sha %s, path: %s", sha, result.Location)
		log.C(ctx).Infow(msg)
		// rm tempfile // works in defer
		// write userfile
		userFile.MetaId = result.Id
		_, err = f.fileRepo.InsertUserFile(ctx, userFile)
		if err != nil {
			log.C(ctx).Errorw("InsertUserFile failed", "userFile", userFile)
			return &errno.Errno{HTTP: http.StatusInternalServerError, Message: "save error"}
		}
		return nil
	}

	// Rename the temporary file to the desired location
	finalFilename, _ := util.GenerateRandomString(16) // Replace with your desired file path
	finalFilename = fmt.Sprintf("%d%d%d%d-%s", createTime.Year(), createTime.Month(), createTime.Day(), createTime.Hour(), finalFilename)
	finalFilepath := filepath.Join(SAVE_FILE_PATH, finalFilename)
	err = os.Rename(tempFile.Name(), finalFilepath)
	if err != nil {
		// If there was an error while renaming, remove the temporary file
		// works in defer
		log.C(ctx).Errorw("upload failed", err)
		return &errno.Errno{HTTP: http.StatusInternalServerError, Message: "save error"}
	}
	fileMeta.Location = finalFilepath
	fileMeta.Sha = sha

	res, err := f.fileRepo.InsertFileMeta(ctx, fileMeta)
	if err != nil {
		log.C(ctx).Errorw("InsertFileMeta failed", "fileMeta", fileMeta)
		return &errno.Errno{HTTP: http.StatusInternalServerError, Message: "save error"}
	}

	fileId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.C(ctx).Errorw("FileMeta ID error", "fileMeta", fileMeta)
		return &errno.Errno{HTTP: http.StatusInternalServerError, Message: "save error"}
	}
	userFile.MetaId = fileId.Hex()

	_, err = f.fileRepo.InsertUserFile(ctx, userFile)
	if err != nil {
		log.C(ctx).Errorw("InsertUserFile failed", "userFile", userFile)
		return &errno.Errno{HTTP: http.StatusInternalServerError, Message: "save error"}
	}

	log.C(ctx).Infow("Upload suc", "userFile", userFile)
	return nil
}

func (f *fileService) QueryUserFile(ctx context.Context, q *v1.UserFileQuery) ([]v1.FileResponse, error) {
	if len(q.UserId) < 1 {
		return nil, &errno.Errno{HTTP: http.StatusBadRequest, Message: "request illeagal"}
	}
	list, err := f.fileRepo.QueryUserFile(ctx, q)
	if err != nil {
		log.Errorw("QueryUserFile", err)
		return nil, errno.InternalServerError
	}
	ids := make([]string, len(list))
	for i, item := range list {
		ids[i] = item.MetaId
	}

	fileList, err := f.fileRepo.FindByMetaId(ctx, ids)
	if err != nil {
		log.Errorw("QueryUserFile", err)
		return nil, errno.InternalServerError
	}
	fileMap := make(map[string]model.FileMeta)
	for _, obj := range fileList {
		fileMap[obj.Id] = obj
	}
	result := make([]v1.FileResponse, len(list))
	for i, item := range list {
		r := v1.FileResponse{
			Id:        item.Id,
			Name:      item.Name,
			Size:      fileMap[item.MetaId].Size,
			CreatedAt: item.CreatedAt,
		}
		result[i] = r
	}
	return result, nil
}
