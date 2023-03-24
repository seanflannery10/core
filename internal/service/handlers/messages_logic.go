package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/seanflannery10/core/internal/api"

	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/data"
)

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

func newMessage(ctx context.Context, q data.Queries, text string, uid int64) (api.MessageResponse, error) {
	message, err := q.CreateMessage(ctx, data.CreateMessageParams{
		Message: text,
		UserID:  uid,
	})
	if err != nil {
		return api.MessageResponse{}, fmt.Errorf("failed create message: %w", err)
	}

	messageResponse := api.MessageResponse{
		ID:      message.ID,
		Message: message.Message,
		Version: message.Version,
	}

	return messageResponse, nil
}

func getMessage(ctx context.Context, q data.Queries, messageID int64) (api.MessageResponse, error) {
	message, err := q.GetMessage(ctx, messageID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return api.MessageResponse{}, errNotFound
		default:
			return api.MessageResponse{}, fmt.Errorf("failed get message: %w", err)
		}
	}

	messageResponse := api.MessageResponse{
		ID:      message.ID,
		Message: message.Message,
		Version: message.Version,
	}

	return messageResponse, nil
}

func updateMessage(ctx context.Context, q data.Queries, m string, id int64) (api.MessageResponse, error) {
	message, err := q.UpdateMessage(ctx, data.UpdateMessageParams{
		Message: m,
		ID:      id,
	})
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return api.MessageResponse{}, errNotFound
		default:
			return api.MessageResponse{}, fmt.Errorf("failed update message: %w", err)
		}
	}

	// if r.Header.Get("X-Expected-Version") != "" {
	//	if strconv.FormatInt(int64(message.Version), 32) != r.Header.Get("X-Expected-Version") {
	//		_ = render.Render(w, r, errs.ErrEditConflict())
	//		return
	//	}
	//}

	messageResponse := api.MessageResponse{
		ID:      message.ID,
		Message: message.Message,
		Version: message.Version,
	}

	return messageResponse, nil
}

func deleteMessage(ctx context.Context, q data.Queries, id int64) (api.AcceptanceResponse, error) {
	err := q.DeleteMessage(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return api.AcceptanceResponse{}, errNotFound
		default:
			return api.AcceptanceResponse{}, fmt.Errorf("failed delete message: %w", err)
		}
	}

	return api.AcceptanceResponse{Message: "message deleted"}, nil
}
