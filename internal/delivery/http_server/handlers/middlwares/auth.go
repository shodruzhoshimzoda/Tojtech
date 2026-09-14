package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shodruzhoshimzoda/tojtech/pkg/httphelpers"
)

type ctxKey string

const (
	ctxUserUUIDKey ctxKey = "user_uuid"
	ctxUerRoleKey  ctxKey = "user_role"
)

// RequireAuth - check that request has a valid JWT, put role and uuid in context
func RequireAuth(jwtSecret []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
			if !ok || tokenString == "" {
				httphelpers.RespondWarn(r.Context(), w, r, http.StatusUnauthorized, "missing bearer token", "unauthorized")
				return
			}

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return jwt.ErrSignatureInvalid, nil
				}

				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				httphelpers.RespondWarn(r.Context(), w, r, http.StatusUnauthorized, "invalid or expired token", "unauthorized")
				return
			}

			sub := claims["sub"].(string)
			role := claims["role"].(string)

			ctx := context.WithValue(r.Context(), ctxUserUUIDKey, sub)
			ctx = context.WithValue(ctx, ctxUerRoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})

	}
}

// RequireRole - comes after RequireAth and check that role in the list of access roles
func RequireRole(allowedRoles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role := r.Context().Value(ctxUerRoleKey).(string)
			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			httphelpers.RespondWarn(r.Context(), w, r, http.StatusForbidden, "insufficient permissions", "you dont have permissions to perform this action")

		})

	}
}

// UserUUIDFromContext - helper for handlers which needs UUID
func UserUUIDFromContext(ctx context.Context) (string, bool) {
	u, ok := ctx.Value(ctxUserUUIDKey).(string)

	return u, ok

}
