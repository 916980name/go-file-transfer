package controller

import (
	"context"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/common"
	"file-transfer/pkg/errno"
	"net/http"
)

func (fc *FileController) CloudinaryUploadFile(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	file, head, err := r.FormFile("file")
	title := r.FormValue("title")
	desc := r.FormValue("desc")
	source := r.FormValue("source")
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer file.Close()

	userId := ctx.Value(common.Trace_request_uid{}).(string)
	req := &v1.CloudinaryFileUpReq{
		UserId:    userId,
		Title:     title,
		Desc:      desc,
		OriginUrl: source,
	}

	m, err := fc.fileService.CloudinaryUploadFile(ctx, file, head, req)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, err)
		return
	}
	errno.WriteResponse(ctx, w, m)
}
