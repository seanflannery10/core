package logic

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/shared/pagination"
)

func GetUserMessages(ctx context.Context, q *data.Queries, page, pageSize int32, userID int64) (*api.MessagesResponse, error) {
	p := pagination.New(page, pageSize)

	messagesFromDB, err := q.GetUserMessages(ctx, data.GetUserMessagesParams{UserID: userID, Offset: p.Offset(), Limit: p.Limit()})
	if err != nil {
		return nil, fmt.Errorf("failed get user messages: %w", err)
	}

	count, err := q.GetUserMessageCount(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed get user message count: %w", err)
	}

	metadata, err := p.CalculateMetadata(count)
	if err != nil {
		return nil, pagination.ErrPageValueToHigh
	}

	messages := make([]api.MessageResponse, len(messagesFromDB))
	for i, v := range messagesFromDB {
		messages[i] = api.MessageResponse{ID: v.ID, Message: v.Message, Version: v.Version}
	}

	messagesResponse := &api.MessagesResponse{Messages: messages, Metadata: metadata}

	return messagesResponse, nil
}

func NewMessage(ctx context.Context, q *data.Queries, m string, userID int64) (*api.MessageResponse, error) {
	message, err := q.CreateMessage(ctx, data.CreateMessageParams{Message: m, UserID: userID})
	if err != nil {
		return nil, fmt.Errorf("failed create message: %w", err)
	}

	messageResponse := &api.MessageResponse{ID: message.ID, Message: message.Message, Version: message.Version}

	return messageResponse, nil
}

func GetMessage(ctx context.Context, q *data.Queries, mid, uid int64) (*api.MessageResponse, error) {
	message, err := q.GetMessage(ctx, data.GetMessageParams{ID: mid, UserID: uid})
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrMessageNotFound
		default:
			return nil, fmt.Errorf("failed get message: %w", err)
		}
	}

	messageResponse := &api.MessageResponse{ID: message.ID, Message: message.Message, Version: message.Version}

	return messageResponse, nil
}

func UpdateMessage(ctx context.Context, q *data.Queries, m string, mid, uid int64) (*api.MessageResponse, error) {
	message, err := q.UpdateMessage(ctx, data.UpdateMessageParams{Message: m, ID: mid, UserID: uid})
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrMessageNotFound
		default:
			return nil, fmt.Errorf("failed update message: %w", err)
		}
	}

	messageResponse := &api.MessageResponse{ID: message.ID, Message: message.Message, Version: message.Version}

	return messageResponse, nil
}

func DeleteMessage(ctx context.Context, q *data.Queries, mid, uid int64) (*api.AcceptanceResponse, error) {
	_, err := q.DeleteMessage(ctx, data.DeleteMessageParams{ID: mid, UserID: uid})
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrMessageNotFound
		default:
			return nil, fmt.Errorf("failed delete message: %w", err)
		}
	}

	acceptanceResponse := &api.AcceptanceResponse{Message: "message deleted"}

	return acceptanceResponse, nil
}
