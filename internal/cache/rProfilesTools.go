package cache

import (
	"backend-gql/graph/model"
	"database/sql"

	"fmt"
	"log"
	"strings"
	"sync"
)

var GlobalMatrixProfilesTools *RProfilesToolsMap

type RProfilesToolsMap struct {
	mu sync.RWMutex

	profilesTools  map[string]map[string]string
	cProfilesTools []*model.ProfilesTools
}

func InitMatrixPermission() *RProfilesToolsMap {
	return &RProfilesToolsMap{
		profilesTools:  make(map[string]map[string]string),
		cProfilesTools: make([]*model.ProfilesTools, 0),
	}
}

func (p *RProfilesToolsMap) LoadFromDB(db *sql.DB) error {

	query := `SELECT PROFILE_ID,TOOL_ID,CREATED_ON,CREATED_BY,MODIFIED_ON,MODIFIED_BY,IS_ACTIVE,OPERATIONS FROM R_PROFILES_TOOLS WHERE IS_ACTIVE =1 `

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	newMatrix := make(map[string]map[string]string)
	cMatrix := make([]*model.ProfilesTools, 0)

	for rows.Next() {
		var profileID, toolID, operations string
		var createdOn, createdBy, modifiedOn, modifiedBy *string
		var isActive bool
		if err := rows.Scan(&profileID, &toolID, &createdOn, &createdBy, &modifiedOn, &modifiedBy, &isActive, &operations); err != nil {
			return err
		}

		if _, ok := newMatrix[toolID]; !ok {
			newMatrix[toolID] = make(map[string]string)
		}

		cMatrix = append(cMatrix, &model.ProfilesTools{
			ProfileID:  profileID,
			ToolID:     toolID,
			CreatedOn:  createdOn,
			CreatedBy:  NilValidate(createdBy),
			ModifiedOn: modifiedOn,
			ModifiedBy: NilValidate(modifiedBy),
			IsActive:   &isActive,
			Operations: &operations,
		})
		newMatrix[toolID][profileID] = strings.ToUpper(strings.TrimSpace(operations))
	}

	p.mu.Lock()
	p.cProfilesTools = cMatrix
	p.profilesTools = newMatrix
	p.mu.Unlock()

	return nil
}
func (p *RProfilesToolsMap) UpdateMatrix(db *sql.DB) ([]*model.ProfilesTools, error) {
	log.Println("Cache actualizado")

	query := `SELECT PROFILE_ID,TOOL_ID,CREATED_ON,CREATED_BY,MODIFIED_ON,MODIFIED_BY,IS_ACTIVE,OPERATIONS FROM R_PROFILES_TOOLS WHERE IS_ACTIVE =1 `

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar DB para actualización:", err)
	}
	defer rows.Close()

	newMatrix := make(map[string]map[string]string)
	cMatrix := make([]*model.ProfilesTools, 0)

	for rows.Next() {
		var p model.ProfilesTools
		var profileID, toolID, operations string
		var createdOn, createdBy, modifiedOn, modifiedBy *string
		var isActive bool
		if err := rows.Scan(&profileID, &toolID, &createdOn, &createdBy, &modifiedOn, &modifiedBy, &isActive, &operations); err != nil {
			return cMatrix, err
		}

		if _, ok := newMatrix[toolID]; !ok {
			newMatrix[toolID] = make(map[string]string)
		}
		cMatrix = append(cMatrix, &model.ProfilesTools{
			ProfileID:  profileID,
			ToolID:     toolID,
			CreatedOn:  createdOn,
			CreatedBy:  NilValidate(createdBy),
			ModifiedOn: modifiedOn,
			ModifiedBy: NilValidate(modifiedBy),
			IsActive:   &isActive,
			Operations: &operations,
		})
		newMatrix[toolID][profileID] = strings.ToUpper(strings.TrimSpace(*p.Operations))
	}

	p.mu.Lock()
	p.profilesTools = newMatrix
	p.cProfilesTools = cMatrix
	p.mu.Unlock()

	return cMatrix, err
}
func (p *RProfilesToolsMap) PrintMatrix() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	log.Println("\n=== MATRIZ DE PERMISOS ACTUAL ===")
	if len(p.profilesTools) == 0 {
		log.Println("La matriz está vacía.")
		return
	}
	for tools, profiles := range p.profilesTools {
		log.Println("TOOL", tools)
		for profile_id, operation := range profiles {
			log.Printf("profile: %s operation: %s \n", profile_id, operation)

		}
	}
}
func (m *RProfilesToolsMap) GetProfileByKey(profileID string, toolID string, actions []string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	profiles, ok := m.profilesTools[toolID]
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
func (p *RProfilesToolsMap) GetAllProfilesTools() []*model.ProfilesTools {

	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.cProfilesTools
}
