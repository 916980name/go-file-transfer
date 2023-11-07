package service

import (
	"context"
	"file-transfer/internal/file-transfer/repo"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/model"
)

type MessageService interface {
	QueryMessage(ctx context.Context, query *v1.MessageQuery) (tags []model.Message, err error)
}

type messageService struct {
	messageRepo repo.MessageRepo
}

var _ MessageService = (*messageService)(nil)

func NewMessageService(repo repo.MessageRepo) MessageService {
	return &messageService{messageRepo: repo}
}

func (s *messageService) QueryMessage(ctx context.Context, query *v1.MessageQuery) (tags []model.Message, err error) {
	return s.messageRepo.Query(ctx, query)
}