package service

import (
	"context"
	"file-transfer/internal/file-transfer/repo"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/model"
	"file-transfer/pkg/util"
	"time"
)

type MessageService interface {
	QueryMessage(ctx context.Context, query *v1.MessageQuery) (tags []model.Message, err error)
	SendMessage(ctx context.Context, r *v1.MessageSendRequest, userId string) error
}

type messageService struct {
	messageRepo repo.MessageRepo
}

var _ MessageService = (*messageService)(nil)

func NewMessageService(repo repo.MessageRepo) MessageService {
	return &messageService{messageRepo: repo}
}

func (s *messageService) QueryMessage(ctx context.Context, query *v1.MessageQuery) (tags []model.Message, err error) {
	if len(query.UserId) < 1 {
		return nil, &errno.Errno{Message: "request illeagal"}
	}
	return s.messageRepo.Query(ctx, query)
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
