package cache

import (
	"backend-gql/graph/model"
	"database/sql"

	"fmt"
	"log"
	"strings"
	"sync"
)

var GlobalMatrixPermission *PermissionHashMap

type PermissionHashMap struct {
	mu sync.RWMutex

	matrix map[string]map[string]string
}

func InitMatrixPermission() *PermissionHashMap {
	return &PermissionHashMap{
		matrix: make(map[string]map[string]string),
	}
}

func (p *PermissionHashMap) LoadFromDB(db *sql.DB) error {

	query := `SELECT PROFILE_ID,TOOL_ID,OPERATIONS FROM R_PROFILES_TOOLS WHERE IS_ACTIVE =1 `

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	newMatrix := make(map[string]map[string]string)

	for rows.Next() {
		var profileID, toolID, operation string
		if err := rows.Scan(&profileID, &toolID, &operation); err != nil {
			return err
		}

		if _, ok := newMatrix[toolID]; !ok {
			newMatrix[toolID] = make(map[string]string)
		}

		newMatrix[toolID][profileID] = strings.ToUpper(strings.TrimSpace(operation))
	}

	p.mu.Lock()
	p.matrix = newMatrix
	p.mu.Unlock()

	return nil
}
func (p *PermissionHashMap) UpdateMatrix(db *sql.DB) ([]*model.ProfilesTools, error) {
	log.Println("Cache actualizado")

	query := `SELECT PROFILE_ID, TOOL_ID, CREATED_ON, CREATED_BY, MODIFIED_ON, MODIFIED_BY,IS_ACTIVE,OPERATIONS FROM R_PROFILES_TOOLS WHERE IS_ACTIVE =1 `
	var data []*model.ProfilesTools

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar DB para actualización:", err)
	}
	defer rows.Close()

	newMatrix := make(map[string]map[string]string)

	for rows.Next() {
		var p model.ProfilesTools

		if err := rows.Scan(&p.ProfileID, &p.ToolID, &p.CreatedOn, &p.CreatedBy, &p.ModifiedOn, &p.ModifiedBy, &p.IsActive, &p.Operations); err != nil {
			log.Println("Error al leer: ", err.Error())
			return data, err
		}

		if _, ok := newMatrix[p.ToolID]; !ok {
			newMatrix[p.ToolID] = make(map[string]string)
		}
		data = append(data, &p)

		newMatrix[p.ToolID][p.ProfileID] = strings.ToUpper(strings.TrimSpace(*p.Operations))
	}

	p.mu.Lock()
	p.matrix = newMatrix
	p.mu.Unlock()

	return data, err
}
func (p *PermissionHashMap) PrintMatrix() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	log.Println("\n=== MATRIZ DE PERMISOS ACTUAL ===")
	if len(p.matrix) == 0 {
		log.Println("La matriz está vacía.")
		return
	}
	for tools, profiles := range p.matrix {
		log.Println("TOOL", tools)
		for profile_id, operation := range profiles {
			log.Printf("profile: %s operation: %s \n", profile_id, operation)

		}
	}
}
func (m *PermissionHashMap) LookInMatrix(profileID string, toolID string, actions []string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	profiles, ok := m.matrix[toolID]
	if !ok {
		return ""
	}

	permission, exist := profiles[profileID]
	if !exist {
		return ""
	}

	if permission == "*" {
		return "*"
	}

	for _, ope := range actions {
		if strings.Contains(permission, ope) {
			log.Printf("Match encontrado: Perfil %s, Herramienta %s, Acción %s", profileID, toolID, ope)
			return permission
		}
	}

	return ""
}
