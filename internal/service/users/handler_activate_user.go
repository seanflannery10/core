package users

import (
	"crypto/sha256"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/validator"
)

type activateUserPayload struct {
	data.UpdateUserParams
	data.DeleteAllTokensForUserParams
	TokenPlaintext string `json:"token"`
}

func (p *activateUserPayload) Bind(r *http.Request) error {
	v := validator.New()

	validateTokenPlaintext(v, p.TokenPlaintext)

	if v.HasErrors() {
		return v.Get()
	}

	queries := helpers.ContextGetQueries(r)

	tokenHash := sha256.Sum256([]byte(p.TokenPlaintext))

	user, err := queries.GetUserFromToken(r.Context(), data.GetUserFromTokenParams{
		Hash:   tokenHash[:],
		Scope:  data.ScopeActivation,
		Expiry: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			v.AddError("token", "invalid or expired activation token")
			return v.Get()
		default:
			return err
		}
	}

	p.UpdateUserParams = data.UpdateUserParams{
		UpdateActivated: true,
		Activated:       true,
		ID:              user.ID,
		Version:         user.Version,
	}

	p.DeleteAllTokensForUserParams = data.DeleteAllTokensForUserParams{
		Scope:  data.ScopeActivation,
		UserID: user.ID,
	}

	return nil
}

func activateUserHandler(w http.ResponseWriter, r *http.Request) {
	input := &activateUserPayload{}

	err := render.Bind(r, input)
	if err != nil {
		_ = render.Render(w, r, errs.ErrBadRequest(err))
		return
	}

	queries := helpers.ContextGetQueries(r)

	user, err := queries.UpdateUser(r.Context(), input.UpdateUserParams)
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	err = queries.DeleteAllTokensForUser(r.Context(), input.DeleteAllTokensForUserParams)
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
	}

	render.Status(r, http.StatusOK)
	err = render.Render(w, r, NewUserResponse(user))
	if err != nil {
		return
	}
}
