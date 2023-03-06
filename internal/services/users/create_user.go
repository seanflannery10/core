package users

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/services"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type createUserPayload struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash []byte `json:"-"`
	env          *services.Env
}

func (p *createUserPayload) Bind(r *http.Request) error {
	v := validator.New()

	data.ValidateName(v, p.Name)
	data.ValidateEmail(v, p.Email)
	data.ValidatePasswordPlaintext(v, p.Password)

	if v.HasErrors() {
		return validator.NewValidationError(v.Errors)
	}

	queries := p.env.Queries

	ok, err := queries.CheckUser(r.Context(), p.Email)
	if err != nil {
		return err
	}

	if ok {
		v.AddError("email", "a user with this email address already exists")
		return validator.NewValidationError(v.Errors)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), 14)
	if err != nil {
		return err
	}

	p.PasswordHash = hash

	return nil
}

func CreateUserHandler(env *services.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &createUserPayload{env: env}

		if helpers.CheckAndBind(w, r, p) {
			return
		}

		user, err := env.Queries.CreateUser(r.Context(), data.CreateUserParams{
			Name:         p.Name,
			Email:        p.Email,
			PasswordHash: p.PasswordHash,
			Activated:    false,
		})
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		token, err := env.Queries.CreateTokenHelper(r.Context(), user.ID, 3*24*time.Hour, data.ScopeActivation)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		err = env.Mailer.Send(user.Email, "token_activation.tmpl", map[string]any{
			"activationToken": token.Plaintext,
		})
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		render.Status(r, http.StatusCreated)

		helpers.RenderAndCheck(w, r, &user)
	}
}
