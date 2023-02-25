package users

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/validator"
	"golang.org/x/exp/slog"
)

var errUserExists = errors.New("user exists")

type createUserPayload struct {
	data.CreateUserParams
	Password string `json:"password"`
}

func (p *createUserPayload) Bind(r *http.Request) error {
	v := validator.New()

	validateName(v, p.CreateUserParams.Name)
	validateEmail(v, p.CreateUserParams.Email)
	validatePasswordPlaintext(v, p.Password)

	if v.HasErrors() {
		return v.Get()
	}

	queries := helpers.ContextGetQueries(r)

	ok, err := queries.CheckUser(r.Context(), p.Email)
	if err != nil {
		return err
	}

	if ok {
		return errUserExists
	}

	hash, err := data.GetPasswordHash(p.Password)
	if err != nil {
		return err
	}

	p.CreateUserParams.Name = p.Name
	p.CreateUserParams.Email = p.Email
	p.CreateUserParams.PasswordHash = hash

	return nil
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	input := &createUserPayload{}

	err := render.Bind(r, input)
	if err != nil {
		switch {
		case errors.Is(err, errUserExists):
			_ = render.Render(w, r, errs.ErrUserExists)
		default:
			_ = render.Render(w, r, errs.ErrBadRequest(err))
		}

		return
	}

	queries := helpers.ContextGetQueries(r)

	user, err := queries.CreateUser(r.Context(), input.CreateUserParams)
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	token, err := queries.NewToken(r.Context(), user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	mailer := helpers.ContextGetMailer(r)

	err = mailer.Send(user.Email, "token_activation.tmpl", map[string]any{
		"activationToken": token.Plaintext,
	})
	if err != nil {
		slog.Error("email error", err)
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, NewUserResponse(user))
	if err != nil {
		return
	}
}
