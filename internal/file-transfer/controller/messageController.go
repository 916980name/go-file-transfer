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
		errno.WriteErrorResponse(ctx, w, errno.ErrInvalidParameter)
		return
	}
	messageRequest.UserId = ctx.Value(common.CTX_USER_KEY).(string)
	util.DebugPrintObj(ctx, messageRequest)

	result, err := mc.service.QueryMessage(ctx, messageRequest)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, err)
		return
	}
	errno.WriteResponse(ctx, w, result)
}

func (mc *MessageController) SendMessage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	messageRequest := &v1.MessageSendRequest{}
	err := util.HttpReadBody(r, messageRequest)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, err)
		return
	}
	userId := ctx.Value(common.CTX_USER_KEY).(string)
	util.DebugPrintObj(ctx, messageRequest)

	err = mc.service.SendMessage(ctx, messageRequest, userId)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (mc *MessageController) DeleteMessage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mId := vars["mId"]
	if len(mId) < 1 {
		errno.WriteErrorResponse(ctx, w, &errno.Errno{Message: "invalid"})
		return
	}

	userId := ctx.Value(common.CTX_USER_KEY).(string)

	err := mc.service.DeleteMessage(ctx, mId, userId)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
