package handlers

import (
	"context"

	"github.com/seanflannery10/core/internal/api"
)

// func (s *Service) GetUserMessages(ctx context.Context, params api.GetUserMessagesParams) (r api.GetUserMessagesRes, _ error) {
//	return r, ht.ErrNotImplemented
//}

func (s *Service) NewMessage(ctx context.Context, req *api.MessageRequest) (r api.NewMessageRes, _ error) {
	const uid = 123

	messageResponse, err := newMessage(ctx, s.queries, req.Message, uid)
	if err != nil {
		return &api.NewMessageInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func (s *Service) GetMessage(ctx context.Context, params api.GetMessageParams) (r api.GetMessageRes, _ error) {
	messageResponse, err := getMessage(ctx, s.queries, params.ID)
	if err != nil {
		return &api.GetMessageInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func (s *Service) UpdateMessage(ctx context.Context, req *api.MessageRequest, params api.UpdateMessageParams) (r api.UpdateMessageRes, _ error) {
	messageResponse, err := updateMessage(ctx, s.queries, req.Message, params.ID)
	if err != nil {
		return &api.UpdateMessageInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func (s *Service) DeleteMessage(ctx context.Context, params api.DeleteMessageParams) (r api.DeleteMessageRes, _ error) {
	acceptanceResponse, err := deleteMessage(ctx, s.queries, params.ID)
	if err != nil {
		return &api.DeleteMessageInternalServerError{}, nil
	}

	return &acceptanceResponse, nil
}
