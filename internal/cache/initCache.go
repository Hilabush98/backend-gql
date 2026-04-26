package cache

import (
	"backend-gql/internal/logs"
	"database/sql"
	"errors"
	"fmt"
)

func InitCache(db *sql.DB) error {
	var loadErr error

	GlobalMatrixProfilesTools = InitMatrixPermission()
	err := GlobalMatrixProfilesTools.LoadFromDB(db)
	if err != nil {
		logs.Error("internal/cache/InitCache", fmt.Sprintf("Error cargando ProfilesTools: %v", err))
		loadErr = errors.Join(loadErr, err)
	}

	logs.Debug("internal/cache/InitCache", "Matriz de ProfilesTools cargada exitosamente.")

	GlobalMatrixcTool = InitToolsMap()
	err = GlobalMatrixcTool.LoadFromDB(db)
	if err != nil {

		logs.Error("internal/cache/InitCache", fmt.Sprintf("Error cargando Tools: %v", err))
		loadErr = errors.Join(loadErr, err)
	}
	logs.Debug("internal/cache/InitCache", "Matriz de Tools cargada exitosamente.")

	GlobalMatrixcProfiles = InitProfilesMap()
	err = GlobalMatrixcProfiles.LoadFromDB(db)
	if err != nil {
		logs.Error("internal/cache/InitCache", fmt.Sprintf("Error cargando Profiles: %v", err))
		loadErr = errors.Join(loadErr, err)
	}
	logs.Debug("internal/cache/InitCache", "Matriz de Profiles cargada exitosamente.")

	GlobalMatrixcGroups = InitGroupsMap()
	err = GlobalMatrixcGroups.LoadFromDB(db)
	if err != nil {
		logs.Error("internal/cache/InitCache", fmt.Sprintf("Error cargando Groups: %v", err))
		loadErr = errors.Join(loadErr, err)
	}
	logs.Debug("internal/cache/InitCache", "Matriz de Groups cargada exitosamente")

	GlobalMatrixrGroupsProfiles = InitGroupsProfilesMap()
	err = GlobalMatrixrGroupsProfiles.LoadFromDB(db)

	if err != nil {
		logs.Error("internal/cache/InitCache", fmt.Sprintf("Error cargando GroupsProfiles: %v", err))
		loadErr = errors.Join(loadErr, err)

	}
	logs.Debug("internal/cache/InitCache", "Matriz de GroupsProflies cargada exitosamente.")

	return loadErr
}
