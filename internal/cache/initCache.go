package cache

import (
	"database/sql"
	"log"
)

func InitCache(db *sql.DB) error {

	GlobalMatrixPermission = InitMatrixPermission()
	err := GlobalMatrixPermission.LoadFromDB(db)
	GlobalMatrixPermission.PrintMatrix()
	if err != nil {
		log.Fatalf("Error crítico cargando permisos: %v", err)
	}
	log.Println("Matriz de permisos cargada exitosamente.")
	GlobalMatrixcTool = InitToolsMap()
	err = GlobalMatrixcTool.LoadFromDB(db)
	GlobalMatrixcTool.PrintMapTool()
	if err != nil {
		log.Fatalf("Error crítico cargando Tools: %v", err)
	}
	log.Println("Matriz de Tools cargada exitosamente.")

	log.Println("Matriz de permisos cargada exitosamente.")
	GlobalMatrixcProfile = InitProfilesMap()
	err = GlobalMatrixcProfile.LoadFromDB(db)
	GlobalMatrixcProfile.PrintMapTool()
	if err != nil {
		log.Fatalf("Error crítico cargando Profiles: %v", err)
	}
	log.Println("Matriz de Profiles cargada exitosamente.")
	return err
}
