package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	mwlogger "github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server/handlers/middlwares"
	"github.com/shodruzhoshimzoda/tojtech/internal/domain/dto"
	domain_user "github.com/shodruzhoshimzoda/tojtech/internal/domain/user"
	user_usecase "github.com/shodruzhoshimzoda/tojtech/internal/usecase/user"
	"github.com/shodruzhoshimzoda/tojtech/pkg/httphelpers"
)

type AdminUserHandler struct {
	usc *user_usecase.AdminUserUsecase
}

func NewAdminUserHandler(uc *user_usecase.AdminUserUsecase) *AdminUserHandler {
	return &AdminUserHandler{usc: uc}
}


func (h *AdminUserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	const op = "AdminUserHandler.ListUsers"

	users, err := h.usc.GetUsers(r.Context())
	if err != nil {
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to list users", err, "internal server error", op)
		return
	}

	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]any{"users": users})
}


func (h *AdminUserHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {
	const op = "AdminUserHandler.ChangeRole"

	// 1. who is changed: uuid from the URL
	targetUUID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "invalid user uuid", "invalid user id")
		return
	}

	//  who asks: uuid from the token (RequireAuth put it into the context)
	callerRaw, ok := mwlogger.UserUUIDFromContext(r.Context())
	if !ok {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusUnauthorized, "caller uuid is missing in context", "unauthorized")
		return
	}
	callerUUID, err := uuid.Parse(callerRaw)
	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusUnauthorized, "caller uuid is invalid", "unauthorized")
		return
	}

	// 3. what to set: JSON body
	var req dto.ChangeRoleDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to decode request body", "invalid request body")
		return
	}

	// 4. do it and translate domain errors into HTTP statuses
	if err := h.usc.ChangeRole(r.Context(), callerUUID, targetUUID, req); err != nil {
		switch {
		case errors.Is(err, domain_user.ErrInvalidRole):
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "invalid role", err.Error())
		case errors.Is(err, domain_user.ErrCannotChangeOwnRole):
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusConflict, "admin tried to change own role", err.Error())
		case errors.Is(err, domain_user.ErrUserNotFound):
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "user not found", "user not found")
		default:
			httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to change role", err, "internal server error", op)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}