package handler

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/logic"
	"github.com/seanflannery10/core/internal/shared/utils"
)

func (s *Handler) GetUserMessages(ctx context.Context, params api.GetUserMessagesParams) (*api.MessagesResponse, error) {
	user := utils.ContextGetUser(ctx)

	messageResponse, err := logic.GetUserMessages(ctx, s.Queries, params.Page.Value, params.PageSize.Value, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed get user messages")
	}

	return messageResponse, nil
}

func (s *Handler) NewMessage(ctx context.Context, req *api.MessageRequest) (*api.MessageResponse, error) {
	user := utils.ContextGetUser(ctx)

	messageResponse, err := logic.NewMessage(ctx, s.Queries, req.Message, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed mew message")
	}

	return messageResponse, nil
}

func (s *Handler) GetMessage(ctx context.Context, params api.GetMessageParams) (*api.MessageResponse, error) {
	user := utils.ContextGetUser(ctx)

	messageResponse, err := logic.GetMessage(ctx, s.Queries, params.ID, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed get message")
	}

	return messageResponse, nil
}

func (s *Handler) UpdateMessage(ctx context.Context, req *api.MessageRequest, params api.UpdateMessageParams) (*api.MessageResponse, error) {
	user := utils.ContextGetUser(ctx)

	messageResponse, err := logic.UpdateMessage(ctx, s.Queries, req.Message, params.ID, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed update message")
	}

	// if r.Header.Get("X-Expected-Version") != "" {
	//	if strconv.FormatInt(int64(message.Version), 32) != r.Header.Get("X-Expected-Version") {
	//		_ = render.Render(w, r, errs.ErrEditConflict())
	//		return
	//	}
	//}

	return messageResponse, nil
}

func (s *Handler) DeleteMessage(ctx context.Context, params api.DeleteMessageParams) (*api.AcceptanceResponse, error) {
	user := utils.ContextGetUser(ctx)

	acceptanceResponse, err := logic.DeleteMessage(ctx, s.Queries, params.ID, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed delete message")
	}

	return acceptanceResponse, nil
}
