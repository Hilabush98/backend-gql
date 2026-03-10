package cache

import (
	"backend-gql/internal/logs"
	"database/sql"
	"fmt"
)

func InitCache(db *sql.DB) error {

	GlobalMatrixProfilesTools = InitMatrixPermission()
	err := GlobalMatrixProfilesTools.LoadFromDB(db)
	if err != nil {
		logs.Error("internal/cache/InitCache", fmt.Sprintf("Error: %w", err))
	}

	logs.Debug("internal/cache/InitCache", "Matriz de ProfilesTools cargada exitosamente.")

	GlobalMatrixcTool = InitToolsMap()
	err = GlobalMatrixcTool.LoadFromDB(db)
	if err != nil {

		logs.Error("internal/cache/InitCache", fmt.Sprintf("Error: %w", err))
	}
	logs.Debug("internal/cache/InitCache", "Matriz de Tools cargada exitosamente.")

	GlobalMatrixcProfiles = InitProfilesMap()
	err = GlobalMatrixcProfiles.LoadFromDB(db)
	if err != nil {
		logs.Error("internal/cache/InitCache", fmt.Sprintf("Error: %w", err))
	}
	logs.Debug("internal/cache/InitCache", "Matriz de Profiles cargada exitosamente.")

	GlobalMatrixcGroups = InitGroupsMap()
	err = GlobalMatrixcGroups.LoadFromDB(db)
	if err != nil {
		logs.Error("internal/cache/InitCache", fmt.Sprintf("Error: %w", err))
	}
	logs.Debug("internal/cache/InitCache", "Matriz de Groups cargada exitosamente")

	GlobalMatrixrGroupsProfiles = InitGroupsProfilesMap()
	err = GlobalMatrixrGroupsProfiles.LoadFromDB(db)

	if err != nil {
		logs.Error("internal/cache/InitCache", fmt.Sprintf("Error: %w", err))

	}
	logs.Debug("internal/cache/InitCache", "Matriz de GroupsProflies cargada exitosamente.")

	return err
}
