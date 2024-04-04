package service

import (
	"context"
	"file-transfer/internal/file-transfer/repo"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/common"
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

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TEMP_FILE_DIR     string = "/tmp"
	TEMP_FILE_PATTERN string = "file-transfer-upload-*.tmp"
	// limit reader end with EOF, but don't know is it real end or reach the limit
	MAX_SINGLE_FILE_SIZE int64 = 50*1024*1024 + 1
)

var SAVE_FILE_PATH string

type FileService interface {
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, userId string) error
	QueryUserFile(ctx context.Context, q *v1.UserFileQuery) ([]v1.FileResponse, error)
	DownloadFile(ctx context.Context, userFileId string, userId string) (*v1.FileDownloadData, error)
	Share(ctx context.Context, mId string, userId string, expireParam *v1.MessageShareParam) (string, error)
	ReadShare(ctx context.Context, key string) (*v1.FileDownloadData, error)
}

type fileService struct {
	fileRepo  repo.FileRepo
	shareServ ShareService
}

var _ FileService = (*fileService)(nil)

func NewFileService(fileRepo repo.FileRepo, shareServ ShareService) FileService {
	workingPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	log.Infow("Read Working dir: " + workingPath)
	SAVE_FILE_PATH = viper.GetString("upload.path")
	log.Infow("Read Saving dir: " + workingPath)
	err = util.CreateDirectoryIfNotExists(SAVE_FILE_PATH)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	log.Infow("Check Save Upload dir: " + SAVE_FILE_PATH)
	return &fileService{fileRepo: fileRepo, shareServ: shareServ}
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
	fileMeta.Location = finalFilename
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

func (f *fileService) DownloadFile(ctx context.Context, userFileId string, userId string) (*v1.FileDownloadData, error) {
	if len(userFileId) < 1 {
		return nil, &errno.Errno{HTTP: http.StatusBadRequest, Message: "request illeagal"}
	}
	userFile, err := f.fileRepo.QueryUserFileById(ctx, userFileId)
	if err != nil {
		return nil, errno.ErrPageNotFound
	}
	if userFile.UserId != userId {
		return nil, errno.ErrPageNotFound
	}

	results, err := f.fileRepo.FindByMetaId(ctx, []string{userFile.MetaId})
	if err != nil || len(results) != 1 {
		return nil, errno.ErrPageNotFound
	}
	finalFilepath := filepath.Join(SAVE_FILE_PATH, results[0].Location)
	return &v1.FileDownloadData{
		Location: finalFilepath,
		Size:     results[0].Size,
		Name:     userFile.Name,
	}, nil
}

func (f *fileService) Share(ctx context.Context, mId string, userId string, expireParam *v1.MessageShareParam) (string, error) {
	file, err := f.fileRepo.QueryUserFileById(ctx, mId)
	if err != nil {
		return "", &errno.Errno{HTTP: http.StatusBadRequest, Message: "invalid"}
	}
	if file.UserId != userId {
		return "", &errno.Errno{HTTP: http.StatusNotAcceptable, Message: "invalid File"}
	}
	switch expireParam.ExpireType {
	case common.SHARE_EXPIRE_TYPE_DURATION:
		return f.shareServ.CreateShareUrl(ctx, common.SHARE_TYPE_FILE, mId, time.Duration(expireParam.Expire*int64(time.Minute)))
	case common.SHARE_EXPIRE_TYPE_TIMES:
		return f.shareServ.CreateShareUrlWithTimes(ctx, common.SHARE_TYPE_FILE, mId, FILE_SHARE_LINK_EXPIRE, int8(expireParam.Expire))
	default:
		return "", &errno.Errno{HTTP: http.StatusMethodNotAllowed, Message: "invalid type"}
	}
}

func (f *fileService) ReadShare(ctx context.Context, key string) (*v1.FileDownloadData, error) {
	userFileId, err := f.shareServ.ConsumeShareUrl(ctx, common.SHARE_TYPE_FILE, key, FILE_SHARE_LINK_EXPIRE)
	if err != nil {
		return nil, &errno.Errno{HTTP: http.StatusBadRequest, Message: "invalid"}
	}
	userFile, err := f.fileRepo.QueryUserFileById(ctx, userFileId)
	if err != nil {
		return nil, errno.ErrPageNotFound
	}

	results, err := f.fileRepo.FindByMetaId(ctx, []string{userFile.MetaId})
	if err != nil || len(results) != 1 {
		return nil, errno.ErrPageNotFound
	}
	finalFilepath := filepath.Join(SAVE_FILE_PATH, results[0].Location)
	return &v1.FileDownloadData{
		Location: finalFilepath,
		Size:     results[0].Size,
		Name:     userFile.Name,
	}, nil
}
