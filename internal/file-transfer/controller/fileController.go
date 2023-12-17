package controller

import (
	"context"
	"file-transfer/internal/file-transfer/service"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/common"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/util"
	"net/http"

	"github.com/gorilla/mux"
)

type FileController struct {
	fileService service.FileService
}

func NewFileController(fileService service.FileService) FileController {
	return FileController{fileService: fileService}
}

func (fc *FileController) UploadFile(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	file, head, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer file.Close()

	userId := ctx.Value(common.CTX_USER_KEY).(string)

	err = fc.fileService.UploadFile(ctx, file, head, userId)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (fc *FileController) QueryUserFile(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	request := &v1.UserFileQuery{
		PageNum:  1,
		PageSize: 10,
	}
	err := util.HttpReadBody(r, request)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, errno.ErrInvalidParameter)
		return
	}
	request.UserId = ctx.Value(common.CTX_USER_KEY).(string)

	result, err := fc.fileService.QueryUserFile(ctx, request)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, err)
		return
	}
	errno.WriteResponse(ctx, w, result)
}

func (fc *FileController) DownloadFile(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fId := vars["fId"]
	if len(fId) < 1 {
		errno.WriteErrorResponse(ctx, w, &errno.Errno{Message: "invalid"})
		return
	}
	userId := ctx.Value(common.CTX_USER_KEY).(string)
	data, err := fc.fileService.DownloadFile(ctx, fId, userId)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, err)
		return
	}
	util.DownloadFileHandler(ctx, w, data)

}
