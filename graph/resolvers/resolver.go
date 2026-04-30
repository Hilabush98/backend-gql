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
	RProfilesTools  *cache.RProfilesToolsMap
	ToolsCache      *cache.CToolsMap
	ProfilesCache   *cache.CProfilesMap
	GroupsCache     *cache.CGroupsMap
	RGroupsProfiles *cache.RGroupsProfilesMap
	DB     *sql.DB
	DBAuth *sql.DB
}
