package v1

import (
	"file-transfer/pkg/common"
	"time"
)

type UserFileQuery struct {
	UserId   string `json:"userId,omitempty"`
	PageNum  int64  `json:"pageNum,omitempty"`
	PageSize int64  `json:"pageSize,omitempty"`
}

type FileResponse struct {
	Id        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
}

type FileShareParam struct {
	ExpireType common.ShareExpireTypeKey `json:"expireType,omitempty"`
	Expire     int64                     `json:"expire,omitempty"`
}

type FileDownloadData struct {
	Location string `json:"location"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
}

type CloudinaryFileUpReq struct {
	UserId    string
	Title     string
	Desc      string
	OriginUrl string
}
