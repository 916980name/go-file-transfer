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

func (mc *MessageController) ReadMessageDefault(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	messageRequest := &v1.MessageQuery{
		PageNum:  1,
		PageSize: 10,
	}
	mc.readMessage(ctx, w, messageRequest)
}

func (mc *MessageController) ReadMessageByPage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	messageRequest := &v1.MessageQuery{
		PageNum:  1,
		PageSize: 10,
	}
	err := util.HttpReadBody(r, messageRequest)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, errno.ErrInvalidParameter)
		return
	}
	mc.readMessage(ctx, w, messageRequest)
}

func (mc *MessageController) readMessage(ctx context.Context, w http.ResponseWriter, messageRequest *v1.MessageQuery) {
	messageRequest.UserId = ctx.Value(common.CTX_USER_KEY).(string)

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

func (mc *MessageController) ShareMessage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mId := vars["mId"]
	if len(mId) < 1 {
		errno.WriteErrorResponse(ctx, w, &errno.Errno{Message: "invalid"})
		return
	}
	shareRequest := &v1.MessageShareParam{
		ExpireType: common.SHARE_EXPIRE_TYPE_DURATION,
		Expire:     1}
	err := util.HttpReadBody(r, shareRequest)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, err)
		return
	}
	userId := ctx.Value(common.CTX_USER_KEY).(string)

	url, err := mc.service.ShareMessage(ctx, mId, userId, shareRequest)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, err)
		return
	}
	errno.WriteResponse(ctx, w, url)
}

func (mc *MessageController) ReadShareMessage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	if len(key) < 1 {
		errno.WriteErrorResponse(ctx, w, &errno.Errno{Message: "invalid"})
		return
	}

	msg, err := mc.service.ReadShareMessage(ctx, key)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, err)
		return
	}
	errno.WriteResponse(ctx, w, msg)
}
