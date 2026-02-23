package auth

import (
	"context"
	"encoding/json"
	"log"
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
			log.Printf("DEBUG: Contenido crudo del header auth: [%s]", permsHeader)
			var permissionsData PermissionResponse

			if len(permsHeader) == 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Header error",
				})
				return
			}

			if permsHeader != "" {
				err := json.Unmarshal([]byte(permsHeader), &permissionsData)
				if err != nil {
					log.Printf("Error decodificando JSON de permisos: %v", err)
				}
			}
			ctx := context.WithValue(r.Context(), permissionsKey, &permissionsData)
			r = r.WithContext(ctx)
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
