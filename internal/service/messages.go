package service

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/oas"
	"github.com/seanflannery10/core/internal/shared/pagination"
	"github.com/seanflannery10/core/internal/shared/utils"
)

func (s *Handler) GetUserMessages(ctx context.Context, params oas.GetUserMessagesParams) (r oas.GetUserMessagesRes, _ error) {
	user := utils.ContextGetUser(ctx)

	messageResponse, err := getUserMessages(ctx, s.Queries, params.Page, params.PageSize, user.ID)
	if err != nil {
		return &oas.GetUserMessagesInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func getUserMessages(ctx context.Context, q data.Queries, page, pageSize int32, userID int64) (oas.MessagesResponse, error) {
	p := pagination.New(page, pageSize)

	messagesFromDB, err := q.GetUserMessages(ctx, data.GetUserMessagesParams{
		UserID: userID,
		Offset: p.Offset(),
		Limit:  p.Limit(),
	})
	if err != nil {
		return oas.MessagesResponse{}, fmt.Errorf("failed get user messages: %w", err)
	}

	count, err := q.GetUserMessageCount(ctx, userID)
	if err != nil {
		return oas.MessagesResponse{}, fmt.Errorf("failed get user message count: %w", err)
	}

	metadata, err := p.CalculateMetadata(count)
	if err != nil {
		return oas.MessagesResponse{}, pagination.ErrPageValueToHigh
	}

	messages := make([]oas.MessageResponse, len(messagesFromDB))
	for i, v := range messagesFromDB {
		messages[i] = oas.MessageResponse{
			ID:      v.UserID,
			Message: v.Message,
			Version: v.Version,
		}
	}

	messagesResponse := oas.MessagesResponse{Messages: messages, Metadata: metadata}

	return messagesResponse, nil
}

func (s *Handler) NewMessage(ctx context.Context, req *oas.MessageRequest) (oas.NewMessageRes, error) {
	user := utils.ContextGetUser(ctx)

	messageResponse, err := newMessage(ctx, s.Queries, req.Message, user.ID)
	if err != nil {
		return &oas.NewMessageInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func newMessage(ctx context.Context, q data.Queries, m string, userID int64) (oas.MessageResponse, error) {
	message, err := q.CreateMessage(ctx, data.CreateMessageParams{
		Message: m,
		UserID:  userID,
	})
	if err != nil {
		return oas.MessageResponse{}, fmt.Errorf("failed create message: %w", err)
	}

	messageResponse := oas.MessageResponse{
		ID:      message.ID,
		Message: message.Message,
		Version: message.Version,
	}

	return messageResponse, nil
}

func (s *Handler) GetMessage(ctx context.Context, params oas.GetMessageParams) (oas.GetMessageRes, error) {
	messageResponse, err := getMessage(ctx, s.Queries, params.ID)
	if err != nil {
		return &oas.GetMessageInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func getMessage(ctx context.Context, q data.Queries, messageID int64) (oas.MessageResponse, error) {
	message, err := q.GetMessage(ctx, messageID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return oas.MessageResponse{}, errNotFound
		default:
			return oas.MessageResponse{}, fmt.Errorf("failed get message: %w", err)
		}
	}

	messageResponse := oas.MessageResponse{
		ID:      message.ID,
		Message: message.Message,
		Version: message.Version,
	}

	return messageResponse, nil
}

func (s *Handler) UpdateMessage(ctx context.Context, req *oas.MessageRequest, params oas.UpdateMessageParams) (oas.UpdateMessageRes, error) {
	messageResponse, err := updateMessage(ctx, s.Queries, req.Message, params.ID)
	if err != nil {
		return &oas.UpdateMessageInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func updateMessage(ctx context.Context, q data.Queries, m string, messageID int64) (oas.MessageResponse, error) {
	message, err := q.UpdateMessage(ctx, data.UpdateMessageParams{
		Message: m,
		ID:      messageID,
	})
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return oas.MessageResponse{}, errNotFound
		default:
			return oas.MessageResponse{}, fmt.Errorf("failed update message: %w", err)
		}
	}

	// if r.Header.Get("X-Expected-Version") != "" {
	//	if strconv.FormatInt(int64(message.Version), 32) != r.Header.Get("X-Expected-Version") {
	//		_ = render.Render(w, r, errs.ErrEditConflict())
	//		return
	//	}
	//}

	messageResponse := oas.MessageResponse{
		ID:      message.ID,
		Message: message.Message,
		Version: message.Version,
	}

	return messageResponse, nil
}

func (s *Handler) DeleteMessage(ctx context.Context, params oas.DeleteMessageParams) (oas.DeleteMessageRes, error) {
	acceptanceResponse, err := deleteMessage(ctx, s.Queries, params.ID)
	if err != nil {
		return &oas.DeleteMessageInternalServerError{}, nil
	}

	return &acceptanceResponse, nil
}

func deleteMessage(ctx context.Context, q data.Queries, id int64) (oas.AcceptanceResponse, error) {
	err := q.DeleteMessage(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return oas.AcceptanceResponse{}, errNotFound
		default:
			return oas.AcceptanceResponse{}, fmt.Errorf("failed delete message: %w", err)
		}
	}

	return oas.AcceptanceResponse{Message: "message deleted"}, nil
}
