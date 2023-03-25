package service

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/oas"
)

// func (s *Handler) GetUserMessages(ctx context.Context, params oas.GetUserMessagesParams) (r oas.GetUserMessagesRes, _ error) {
//	return r, ht.ErrNotImplemented
//}

// func getUserMessages(ctx context.Context, q data.Queries, email string) http.HandlerFunc {
//	messages, err := q.GetUserMessages(ctx, data.GetUserMessagesParams{
//		UserID: env.User.ID,
//		Offset: p.Pagination.Offset(),
//		Limit:  p.Pagination.Limit(),
//	})
//	if err != nil {
//		_ = render.Render(w, r, errs.ErrServerError(err))
//	}
//
//	count, err := q.GetUserMessageCount(ctx, env.User.ID)
//	if err != nil {
//		_ = render.Render(w, r, errs.ErrServerError(err))
//	}
//
//	metadata := p.Pagination.CalculateMetadata(count)
//
//	if p.Pagination.Validator.HasErrors() {
//		_ = render.Render(w, r, errs.ErrFailedValidation(p.Pagination.Validator.Errors))
//	}
//}

func (s *Handler) NewMessage(ctx context.Context, req *oas.MessageRequest) (r oas.NewMessageRes, _ error) {
	const uid = 123

	messageResponse, err := newMessage(ctx, s.Queries, req.Message, uid)
	if err != nil {
		return &oas.NewMessageInternalServerError{}, nil
	}

	return &messageResponse, nil
}

func newMessage(ctx context.Context, q data.Queries, text string, uid int64) (oas.MessageResponse, error) {
	message, err := q.CreateMessage(ctx, data.CreateMessageParams{
		Message: text,
		UserID:  uid,
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

func (s *Handler) GetMessage(ctx context.Context, params oas.GetMessageParams) (r oas.GetMessageRes, _ error) {
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

func updateMessage(ctx context.Context, q data.Queries, m string, id int64) (oas.MessageResponse, error) {
	message, err := q.UpdateMessage(ctx, data.UpdateMessageParams{
		Message: m,
		ID:      id,
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

func (s *Handler) DeleteMessage(ctx context.Context, params oas.DeleteMessageParams) (r oas.DeleteMessageRes, _ error) {
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
