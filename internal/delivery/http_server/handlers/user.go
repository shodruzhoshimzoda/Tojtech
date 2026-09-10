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

	op := "ProductHandler.RegisterUser"
	var req dto.RegisterDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to decode request body", "invalid request body")
		return
	}

	_, err := h.usc.RegisterUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain_user.ErrUserAlreadyExists) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "user already exists", "duplicate user")
			return
		}

		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to create user", err, "Internal Server error", op)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	op := "ProductHandler.LoginUser"
	var req dto.LoginDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to decode request body", "invalid request body")
		return
	}

	token, err := h.usc.LoginUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain_user.ErrInvalidEmailOrPassword) {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to create user", err, "Internal Server error", op)
		return
	}

	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]string{
		"token": token,
	})
}
