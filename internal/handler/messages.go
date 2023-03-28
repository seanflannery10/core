package handler

import (
	"context"

	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/logic"
	"github.com/seanflannery10/core/internal/shared/utils"
)

func (s *Handler) GetUserMessages(ctx context.Context, params api.GetUserMessagesParams) (r api.GetUserMessagesRes, _ error) {
	user := utils.ContextGetUser(ctx)

	messageResponse, err := logic.GetUserMessages(ctx, s.Queries, params.Page, params.PageSize, user.ID)
	if err != nil {
		return &api.GetUserMessagesInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func (s *Handler) NewMessage(ctx context.Context, req *api.MessageRequest) (api.NewMessageRes, error) {
	user := utils.ContextGetUser(ctx)

	messageResponse, err := logic.NewMessage(ctx, s.Queries, req.Message, user.ID)
	if err != nil {
		return &api.NewMessageInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func (s *Handler) GetMessage(ctx context.Context, params api.GetMessageParams) (api.GetMessageRes, error) {
	messageResponse, err := logic.GetMessage(ctx, s.Queries, params.ID)
	if err != nil {
		return &api.GetMessageInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func (s *Handler) UpdateMessage(ctx context.Context, req *api.MessageRequest, params api.UpdateMessageParams) (api.UpdateMessageRes, error) {
	messageResponse, err := logic.UpdateMessage(ctx, s.Queries, req.Message, params.ID)
	if err != nil {
		return &api.UpdateMessageInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func (s *Handler) DeleteMessage(ctx context.Context, params api.DeleteMessageParams) (api.DeleteMessageRes, error) {
	acceptanceResponse, err := logic.DeleteMessage(ctx, s.Queries, params.ID)
	if err != nil {
		return &api.DeleteMessageInternalServerError{}, nil
	}

	return &acceptanceResponse, nil
}
