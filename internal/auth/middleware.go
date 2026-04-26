package auth

import (
	"backend-gql/internal/logs"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type contextKey string

const permissionsKey contextKey = "userPermissions"

type Permission struct {
	ToolID *string `json:"tool_id"`

	ProfileID []string `json:"profile_id"`
}
type PermissionResponse struct {
	Permissions Permission `json:"Permissions"`
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			permsHeader := r.Header.Get("auth")
			var permissionsData PermissionResponse
			if len(permsHeader) == 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "missing auth header",
				})
				return
			}

			err := json.Unmarshal([]byte(permsHeader), &permissionsData)
			if err != nil {
				logs.Error("internal/auth/middleware.go", fmt.Sprintf("invalid auth header json: %v", err))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "invalid auth header json",
				})
				return
			}

			ctx := context.WithValue(r.Context(), permissionsKey, &permissionsData)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
func PermissionsFromContext(ctx context.Context) *PermissionResponse {
	if perms, ok := ctx.Value(permissionsKey).(*PermissionResponse); ok {
		return perms
	}
	return nil
}
