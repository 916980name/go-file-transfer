package controller

import (
	"context"
	"file-transfer/internal/file-transfer/service"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/common"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/util"
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

func (mc *MessageController) ReadMessage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	messageRequest := &v1.MessageQuery{
		PageNum:  1,
		PageSize: 10,
	}
	err := util.HttpReadBody(r, messageRequest)
	if err != nil {
		errno.WriteResponse(ctx, w, err, nil)
		return
	}
	messageRequest.UserId = ctx.Value(common.CTX_USER_KEY).(string)
	util.DebugPrintObj(ctx, messageRequest)

	result, err := mc.service.QueryMessage(ctx, messageRequest)
	if err != nil {
		errno.WriteResponse(ctx, w, err, nil)
		return
	}
	errno.WriteResponse(ctx, w, nil, result)
}

func (mc *MessageController) SendMessage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	messageRequest := &v1.MessageSendRequest{}
	err := util.HttpReadBody(r, messageRequest)
	if err != nil {
		errno.WriteResponse(ctx, w, err, nil)
		return
	}
	userId := ctx.Value(common.CTX_USER_KEY).(string)
	util.DebugPrintObj(ctx, messageRequest)

	err = mc.service.SendMessage(ctx, messageRequest, userId)
	if err != nil {
		errno.WriteResponse(ctx, w, err, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
}
