package httphelpers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/shodruzhoshimzoda/tojtech/pkg/reqlog"
)

// RespondError логирует через reqlog.Error и отдает JSON ошибку
func RespondError(ctx context.Context, w http.ResponseWriter, r *http.Request, status int, logMsg string, err error, clientMsg string, op string) {
	reqlog.Error(ctx, logMsg, err, slog.String("op", op))
	render.Status(r, status)
	render.JSON(w, r, map[string]string{"error": clientMsg})
}

// RespondWarn логирует через reqlog.Warn и отдает JSON ошибку
func RespondWarn(ctx context.Context, w http.ResponseWriter, r *http.Request, status int, logMsg string, clientMsg string) {
	reqlog.Warn(ctx, logMsg)
	render.Status(r, status)
	render.JSON(w, r, map[string]string{"error": clientMsg})
}

// RespondWarnWithDesc как RespondWarn но с error_description
func RespondWarnWithDesc(ctx context.Context, w http.ResponseWriter, r *http.Request, status int, logMsg string, clientMsg string, desc string) {
	reqlog.Warn(ctx, logMsg)
	render.Status(r, status)
	render.JSON(w, r, map[string]string{"error": clientMsg, "error_description": desc})
}

// RespondJSON просто отдает JSON с нужным статусом
func RespondJSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	render.Status(r, status)
	render.JSON(w, r, data)
}
