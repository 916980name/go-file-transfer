package v1

type MessageQuery struct {
	UserId   string `json:"userId,omitempty"`
	PageNum  int64  `json:"pageNum,omitempty"`
	PageSize int64  `json:"pageSize,omitempty"`
}

type MessageSendRequest struct {
	Info string `json:"info"`
}
