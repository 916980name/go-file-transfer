package util

import (
	"context"
	"encoding/json"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/log"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func HttpReadBody(r *http.Request, customType interface{}) error {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return err
	}

	// Create an instance of the custom type

	// Unmarshal the request body into the custom type
	err = json.Unmarshal(body, &customType)
	if err != nil {
		// http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return err
	}
	return nil
}

func DownloadFileHandler(ctx context.Context, w http.ResponseWriter, data *v1.FileDownloadData) {
	// Open the file
	file, err := os.Open(data.Location)
	if err != nil {
		log.C(ctx).Errorw("downloadFileHandler open file failed", "file", data.Location)
		errno.WriteErrorResponse(ctx, w, errno.InternalServerError)
		return
	}
	defer file.Close()

	// Set the headers
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.PathEscape(data.Name))) // Replace with the desired filename

	// Stream the file to the response
	_, err = io.Copy(w, file)
	if err != nil {
		log.C(ctx).Errorw("downloadFileHandler Failed to stream file", "file", data.Location, "error", err)
		errno.WriteErrorResponse(ctx, w, errno.InternalServerError)
		return
	}
}
