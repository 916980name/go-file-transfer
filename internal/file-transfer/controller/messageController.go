package controller

import (
	"context"
	"file-transfer/internal/file-transfer/service"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/errno"
	"net/http"
)

type MessageController struct {
	service service.MessageService
}

func NewMessageController(service service.MessageService) MessageController {
	return MessageController{
		service: service,
	}
}

func (mc *MessageController) Onemessage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	result, err := mc.service.QueryMessage(ctx, &v1.MessageQuery{UserId: "A"})
	if err != nil {
		errno.WriteResponse(ctx, w, err, nil)
		return
	}
	errno.WriteResponse(ctx, w, nil, result)
}
