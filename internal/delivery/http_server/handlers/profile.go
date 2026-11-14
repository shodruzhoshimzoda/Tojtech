package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	mwlogger "github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server/handlers/middlwares"
	"github.com/shodruzhoshimzoda/tojtech/internal/domain/dto"
	domain_user "github.com/shodruzhoshimzoda/tojtech/internal/domain/user"
	"github.com/shodruzhoshimzoda/tojtech/pkg/httphelpers"
)

// currentUserUUID reads the caller's uuid that RequireAuth put into the context.
func currentUserUUID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	raw, ok := mwlogger.UserUUIDFromContext(r.Context())
	if !ok {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusUnauthorized, "user uuid is missing in context", "unauthorized")
		return uuid.Nil, false
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusUnauthorized, "user uuid is invalid", "unauthorized")
		return uuid.Nil, false
	}
	return id, true
}

// GetMe - GET /me
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandler.GetMe"

	id, ok := currentUserUUID(w, r)
	if !ok {
		return
	}

	resp, err := h.usc.GetProfile(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain_user.ErrUserNotFound) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusUnauthorized, "user from token not found", "unauthorized")
			return
		}
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to get profile", err, "Internal Server error", op)
		return
	}

	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]any{"user": resp})
}

// UpdateMe - PATCH /me (changes the username)
func (h *AuthHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandler.UpdateMe"

	id, ok := currentUserUUID(w, r)
	if !ok {
		return
	}

	var req dto.UpdateProfileDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to decode request body", "invalid request body")
		return
	}

	resp, err := h.usc.UpdateProfile(r.Context(), id, req)
	if err != nil {
		if domain_user.IsUsernameError(err) {
			httphelpers.RespondWarnWithDesc(r.Context(), w, r, http.StatusBadRequest, "invalid username", "invalid username", err.Error())
			return
		}
		if errors.Is(err, domain_user.ErrUsernameTaken) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusConflict, "username already taken", "this username is already taken")
			return
		}
		if errors.Is(err, domain_user.ErrUserNotFound) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusUnauthorized, "user from token not found", "unauthorized")
			return
		}
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to update profile", err, "Internal Server error", op)
		return
	}

	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]any{"user": resp})
}

// ChangePassword - POST /me/password
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandler.ChangePassword"

	id, ok := currentUserUUID(w, r)
	if !ok {
		return
	}

	var req dto.ChangePasswordDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to decode request body", "invalid request body")
		return
	}

	resp, err := h.usc.ChangePassword(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, dto.ErrPasswordRequired) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to change password", "current and new password are required")
			return
		}

		var policyErr *domain_user.PasswordPolicyError
		if errors.As(err, &policyErr) {
			httphelpers.RespondWarnWithDesc(r.Context(), w, r, http.StatusBadRequest, "password does not meet the policy", "weak password", policyErr.Error())
			return
		}

		// 400, not 401: the frontend treats 401 as "the session is dead" and logs the user out
		if errors.Is(err, domain_user.ErrWrongCurrentPassword) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "wrong current password", "current password is incorrect")
			return
		}
		if errors.Is(err, domain_user.ErrSamePassword) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "new password equals the old one", "new password must be different from the current one")
			return
		}
		if errors.Is(err, domain_user.ErrUserNotFound) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusUnauthorized, "user from token not found", "unauthorized")
			return
		}
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to change password", err, "Internal Server error", op)
		return
	}

	// fresh tokens: all old refresh tokens were revoked, this device gets a new session
	httphelpers.RespondJSON(w, r, http.StatusOK, resp)
}