package handler

import (
	"context"

	"github.com/seanflannery10/core/internal/logic"
	"github.com/seanflannery10/core/internal/oas"
	"github.com/seanflannery10/core/internal/shared/utils"
)

func (s *Handler) GetUserMessages(ctx context.Context, params oas.GetUserMessagesParams) (r oas.GetUserMessagesRes, _ error) {
	user := utils.ContextGetUser(ctx)

	messageResponse, err := logic.GetUserMessages(ctx, s.Queries, params.Page, params.PageSize, user.ID)
	if err != nil {
		return &oas.GetUserMessagesInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func (s *Handler) NewMessage(ctx context.Context, req *oas.MessageRequest) (oas.NewMessageRes, error) {
	user := utils.ContextGetUser(ctx)

	messageResponse, err := logic.NewMessage(ctx, s.Queries, req.Message, user.ID)
	if err != nil {
		return &oas.NewMessageInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func (s *Handler) GetMessage(ctx context.Context, params oas.GetMessageParams) (oas.GetMessageRes, error) {
	messageResponse, err := logic.GetMessage(ctx, s.Queries, params.ID)
	if err != nil {
		return &oas.GetMessageInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func (s *Handler) UpdateMessage(ctx context.Context, req *oas.MessageRequest, params oas.UpdateMessageParams) (oas.UpdateMessageRes, error) {
	messageResponse, err := logic.UpdateMessage(ctx, s.Queries, req.Message, params.ID)
	if err != nil {
		return &oas.UpdateMessageInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func (s *Handler) DeleteMessage(ctx context.Context, params oas.DeleteMessageParams) (oas.DeleteMessageRes, error) {
	acceptanceResponse, err := logic.DeleteMessage(ctx, s.Queries, params.ID)
	if err != nil {
		return &oas.DeleteMessageInternalServerError{}, nil
	}

	return &acceptanceResponse, nil
}
