package cache

import (
	"backend-gql/graph/model"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
)

var GlobalMatrixrGroupsProfiles *RGroupsProfilesMap

type RGroupsProfilesMap struct {
	groupsProfiles  map[string]string
	cGroupsProfiles []*model.GroupsProfiles
	mu              sync.RWMutex
}

func InitGroupsProfilesMap() *RGroupsProfilesMap {
	return &RGroupsProfilesMap{
		groupsProfiles:  make(map[string]string),
		cGroupsProfiles: make([]*model.GroupsProfiles, 0),
	}
}

func (p *RGroupsProfilesMap) LoadFromDB(db *sql.DB) error {
	query := `SELECT GROUP_ID, PROFILE_ID, CREATED_ON,CREATED_BY,MODIFIED_ON, MODIFIED_BY, IS_ACTIVE FROM R_GROUPS_PROFILES WHERE IS_ACTIVE =1 `

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	newMatrix := make(map[string]string)
	cMatrix := make([]*model.GroupsProfiles, 0)

	for rows.Next() {
		var groupID, profileID string
		var createdOn, createdBy, modifiedOn, modifiedBy *string
		var isActive bool
		if err := rows.Scan(&groupID, &profileID, &createdOn, &createdBy, &modifiedOn, &modifiedBy, &isActive); err != nil {
			return err
		}

		cMatrix = append(cMatrix, &model.GroupsProfiles{
			GroupID:    groupID,
			ProfileID:  profileID,
			CreatedOn:  createdOn,
			CreatedBy:  NilValidate(createdBy),
			ModifiedOn: modifiedOn,
			ModifiedBy: NilValidate(modifiedBy),
			IsActive:   &isActive,
		})
		newMatrix[groupID] = strings.ToUpper(strings.TrimSpace(profileID))
	}

	p.mu.Lock()
	p.cGroupsProfiles = cMatrix
	p.groupsProfiles = newMatrix
	p.mu.Unlock()

	return nil
}

func (p *RGroupsProfilesMap) PrintMatrix() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	log.Println("\n=== MATRIZ DE PERMISOS ACTUAL ===")
	if len(p.groupsProfiles) == 0 {
		log.Println("La matriz esta vacia.")
		return
	}

	for groupID, profileID := range p.groupsProfiles {
		log.Printf("group: %s profile: %s\n", groupID, profileID)
	}

	for _, g := range p.cGroupsProfiles {
		log.Println("Grupo:", g)
	}
}

func (g *RGroupsProfilesMap) GetGroupProfileById(id string) string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.groupsProfiles[id]
}

func (g *RGroupsProfilesMap) GetAllGroupsProfiles() []*model.GroupsProfiles {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.cGroupsProfiles
}

func (g *RGroupsProfilesMap) UpdateMatrix(db *sql.DB) ([]*model.GroupsProfiles, error) {
	query := `SELECT GROUP_ID, PROFILE_ID FROM R_GROUPS_PROFILES WHERE IS_ACTIVE =1`
	var data []*model.GroupsProfiles

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar DB para actualizacion: %w", err)
	}
	defer rows.Close()

	newMatrix := make(map[string]string)

	for rows.Next() {
		var p model.GroupsProfiles
		if err := rows.Scan(&p.GroupID, &p.ProfileID); err != nil {
			log.Println("Error al leer:", err.Error())
			return data, err
		}
		data = append(data, &p)
		newMatrix[p.GroupID] = strings.ToUpper(strings.TrimSpace(p.ProfileID))
	}

	g.mu.Lock()
	g.groupsProfiles = newMatrix
	g.cGroupsProfiles = data
	g.mu.Unlock()

	return data, err
}
