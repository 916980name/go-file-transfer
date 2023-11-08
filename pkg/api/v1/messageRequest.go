package v1

import "time"

type MessageQuery struct {
	UserId   string `json:"userId,omitempty"`
	PageNum  int64  `json:"pageNum,omitempty"`
	PageSize int64  `json:"pageSize,omitempty"`
}

type MessageSendRequest struct {
	Info string `json:"info"`
}

type MessageResponse struct {
	Id        string    `json:"id,omitempty"`
	Info      string    `json:"info"`
	CreatedAt time.Time `json:"createdAt"`
}
