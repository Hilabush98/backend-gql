package cache

import (
	"database/sql"
	"log"
)

func InitCache(db *sql.DB) error {

	GlobalMatrixProfilesTools = InitMatrixPermission()
	err := GlobalMatrixProfilesTools.LoadFromDB(db)
	//GlobalMatrixPermission.PrintMatrix()
	if err != nil {
		log.Fatalf("Error crítico cargando permisos: %v", err)
	}
	log.Println("Matriz de permisos cargada exitosamente.")
	GlobalMatrixcTool = InitToolsMap()
	err = GlobalMatrixcTool.LoadFromDB(db)
	//GlobalMatrixcTool.PrintMapTool()
	if err != nil {
		log.Fatalf("Error crítico cargando Tools: %v", err)
	}
	log.Println("Matriz de Tools cargada exitosamente.")

	GlobalMatrixcProfiles = InitProfilesMap()
	err = GlobalMatrixcProfiles.LoadFromDB(db)
	//GlobalMatrixcProfile.PrintMapTool()
	if err != nil {
		log.Fatalf("Error crítico cargando Profiles: %v", err)
	}
	log.Println("Matriz de Profiles cargada exitosamente.")

	GlobalMatrixcGroups = InitGroupsMap()
	err = GlobalMatrixcGroups.LoadFromDB(db)
	//GlobalMatrixcGroups.PrintMapGroups()
	if err != nil {
		log.Fatalf("Error crítico cargando Groups: %v", err)
	}
	log.Println("Matriz de Groups cargada exitosamente.")

	GlobalMatrixrGroupsProfiles = InitGroupsProfilesMap()
	err = GlobalMatrixrGroupsProfiles.LoadFromDB(db)
	//GlobalMatrixrGroupsProfiles.PrintMatrix()
	if err != nil {
		log.Fatalf("Error crítico cargando Groups: %v", err)
	}
	log.Println("Matriz de Groups cargada exitosamente.")
	return err
}
