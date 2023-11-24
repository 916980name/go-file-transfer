package service

import (
	"context"
	"file-transfer/internal/file-transfer/repo"
	v1 "file-transfer/pkg/api/v1"
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
}

type messageService struct {
	messageRepo repo.MessageRepo
}

var _ MessageService = (*messageService)(nil)

func NewMessageService(repo repo.MessageRepo) MessageService {
	return &messageService{messageRepo: repo}
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
