package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shodruzhoshimzoda/tojtech/internal/domain/dto"
	domain_user "github.com/shodruzhoshimzoda/tojtech/internal/domain/user"
	user_usecase "github.com/shodruzhoshimzoda/tojtech/internal/usecase/user"
	"github.com/shodruzhoshimzoda/tojtech/pkg/httphelpers"
)

type AuthHandler struct {
	usc *user_usecase.AuthUsercase
}

func NewAuthHandler(uc *user_usecase.AuthUsercase) *AuthHandler {
	return &AuthHandler{
		usc: uc,
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.RegisterUser"
	var req dto.RegisterDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to decode request body", "invalid request body")
		return
	}

	resp, err := h.usc.RegisterUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain_user.ErrInvalidEmailOrPassword) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "invalid email format", "invalid email or password")
			return
		}
		if errors.Is(err, dto.ErrEmailRequired) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to register user", "email is required")
			return
		}

		if errors.Is(err, dto.ErrPasswordRequired) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to register user", "password is required")
			return
		}

		if errors.Is(err, domain_user.ErrUserAlreadyExists) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusConflict, "user already exists", "user with this email already exists")
			return
		}

		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to create user", err, "Internal Server error", op)
		return
	}

	httphelpers.RespondJSON(w, r, http.StatusCreated, resp)
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.LoginUser"
	var req dto.LoginDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to decode request body", "invalid request body")
		return
	}

	resp, err := h.usc.LoginUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, dto.ErrEmailRequired) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to login user", "email is required")
			return
		}
		if errors.Is(err, dto.ErrPasswordRequired) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to login user", "password is required")
			return
		}

		// Возвращаем 401 Unauthorized для любых неверных учетных данных
		if errors.Is(err, domain_user.ErrInvalidEmailOrPassword) || errors.Is(err, domain_user.ErrUserNotFound) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusUnauthorized, "invalid credentials", "invalid email or password")
			return
		}

		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to login user", err, "Internal Server error", op)
		return
	}

	httphelpers.RespondJSON(w, r, http.StatusOK, resp)
}
