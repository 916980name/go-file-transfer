package service

import (
	"context"
	"file-transfer/internal/file-transfer/repo"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/common"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/log"
	"file-transfer/pkg/model"
	"file-transfer/pkg/util"
	"net/http"
	"time"
)

type MessageService interface {
	QueryMessage(ctx context.Context, query *v1.MessageQuery) ([]v1.MessageResponse, error)
	SendMessage(ctx context.Context, r *v1.MessageSendRequest, userId string) error
	DeleteMessage(ctx context.Context, mId string, userId string) error
	ShareMessage(ctx context.Context, mId string, userId string, expireParam *v1.MessageShareParam) (string, error)
	ReadShareMessage(ctx context.Context, key string) (string, error)
}

type messageService struct {
	messageRepo repo.MessageRepo
	shareServ   ShareService
}

var (
	MESSAGE_SHARE_LINK_EXPIRE time.Duration = 24 * time.Hour
)

var _ MessageService = (*messageService)(nil)

func NewMessageService(repo repo.MessageRepo, shareServ ShareService) MessageService {
	return &messageService{messageRepo: repo, shareServ: shareServ}
}

func (s *messageService) QueryMessage(ctx context.Context, query *v1.MessageQuery) ([]v1.MessageResponse, error) {
	if len(query.UserId) < 1 {
		return nil, &errno.Errno{Message: "request illeagal"}
	}
	if query.PageNum < 1 {
		query.PageNum = 1
	}
	log.C(ctx).Debugw("read msg", query)
	list, err := s.messageRepo.Query(ctx, query)
	if err != nil {
		return nil, errno.ErrBind
	}
	transformed := make([]v1.MessageResponse, len(list))
	for i, msg := range list {
		transformed[i] = v1.MessageResponse{
			Id:        msg.Id,
			Info:      msg.Info,
			CreatedAt: msg.CreatedAt,
		}
	}
	return transformed, err
}

func (s *messageService) SendMessage(ctx context.Context, r *v1.MessageSendRequest, userId string) error {
	m := &model.Message{
		UserId:    userId,
		Info:      r.Info,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	util.DebugPrintObj(ctx, m)
	result, err := s.messageRepo.Insert(ctx, m)
	util.DebugPrintObj(ctx, result)
	return err
}

func (s *messageService) DeleteMessage(ctx context.Context, mId string, userId string) error {
	if len(mId) < 1 || len(userId) < 1 {
		return &errno.Errno{Message: "invalid"}
	}
	result, err := s.messageRepo.Delete(ctx, mId, userId)
	if result != nil && result.DeletedCount != 1 {
		log.C(ctx).Debugw("DeleteMessage", result.DeletedCount)
		return &errno.Errno{HTTP: http.StatusForbidden, Message: "invalid"}
	}
	return err
}

func (s *messageService) ShareMessage(ctx context.Context, mId string, userId string, expireParam *v1.MessageShareParam) (string, error) {
	msg, err := s.messageRepo.QueryById(ctx, mId)
	if err != nil {
		return "", &errno.Errno{HTTP: http.StatusBadRequest, Message: "invalid"}
	}
	if msg.UserId != userId {
		return "", &errno.Errno{HTTP: http.StatusNotAcceptable, Message: "invalid Message"}
	}
	switch expireParam.ExpireType {
	case common.SHARE_EXPIRE_TYPE_DURATION:
		return s.shareServ.CreateShareUrl(ctx, common.SHARE_TYPE_MESSAGE, mId, time.Duration(expireParam.Expire*int64(time.Minute)))
	case common.SHARE_EXPIRE_TYPE_TIMES:
		return s.shareServ.CreateShareUrlWithTimes(ctx, common.SHARE_TYPE_MESSAGE, mId, MESSAGE_SHARE_LINK_EXPIRE, int8(expireParam.Expire))
	default:
		return "", &errno.Errno{HTTP: http.StatusMethodNotAllowed, Message: "invalid type"}
	}
}

func (s *messageService) ReadShareMessage(ctx context.Context, key string) (string, error) {
	mId, err := s.shareServ.ConsumeShareUrl(ctx, common.SHARE_TYPE_MESSAGE, key, MESSAGE_SHARE_LINK_EXPIRE)
	if err != nil {
		return "", &errno.Errno{HTTP: http.StatusBadRequest, Message: "invalid"}
	}
	msg, err := s.messageRepo.QueryById(ctx, mId)
	if err != nil {
		return "", &errno.Errno{HTTP: http.StatusBadRequest, Message: "invalid"}
	}
	return msg.Info, nil
}
