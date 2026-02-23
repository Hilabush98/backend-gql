package graph

import (
	"backend-gql/internal/cache"
	"database/sql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	PermissionCache *cache.PermissionHashMap
	ToolsCache      *cache.CToolsMap
	ProfilesCache   *cache.CProfileMap
	DB              *sql.DB
}
