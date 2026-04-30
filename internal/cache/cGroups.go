package cache

import (
	"backend-gql/graph/model"
	"database/sql"
	"log"
	"sync"
)

var GlobalMatrixcGroups *CGroupsMap

type CGroupsMap struct {
	Groups map[string]model.Group
	mu     sync.RWMutex
}

func InitGroupsMap() *CGroupsMap {
	return &CGroupsMap{
		Groups: make(map[string]model.Group),
	}
}
func (p *CGroupsMap) LoadFromDB(db *sql.DB) error {

	query := `SELECT GROUP_ID, NAME, DESCRIPTION, CREATED_ON, CREATED_BY, MODIFIED_ON, MODIFIED_BY, ORDER_BY, IS_ACTIVE FROM APP.C_GROUPS WHERE IS_ACTIVE =true `

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	newGroups := make(map[string]model.Group)

	for rows.Next() {
		var groupID string
		var createdOn, name, orderBy, createdBy, modifiedOn, modifiedBy, description *string
		var isActive bool
		err := rows.Scan(
			&groupID,
			&name,
			&description,
			&createdOn,
			&createdBy,
			&modifiedOn,
			&modifiedBy,
			&orderBy,
			&isActive,
		)
		if err != nil {
			return err
		}

		t := model.Group{
			GroupID:     groupID,
			Name:        NilValidate(name),
			OrderBy:     NilValidate(orderBy),
			CreatedOn:   createdOn,
			CreatedBy:   NilValidate(createdBy),
			ModifiedOn:  modifiedOn,
			ModifiedBy:  NilValidate(modifiedBy),
			Description: NilValidate(description),
			IsActive:    &isActive,
		}

		newGroups[t.GroupID] = t
	}

	p.mu.Lock()
	p.Groups = newGroups
	p.mu.Unlock()

	return nil

}

func (p *CGroupsMap) PrintMapGroups() {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for id, t := range p.Groups {

		log.Println("ID: ", id)
		log.Println("Name: ", *t.Name)
		log.Println("Description: ", Str(t.Description))
		log.Println("OrderBy: ", *t.OrderBy)
		log.Println("CreatedOn: ", Str(t.CreatedOn))
		log.Println("CreatedBy:", Str(t.CreatedBy))
		log.Println("ModifiedBy: ", Str(t.ModifiedBy))
		log.Println("Is_Active: ", *t.IsActive)

		log.Println("----------------------------")

	}
}
func (m *CGroupsMap) GetAllGroups() []*model.Group {
	m.mu.RLock()
	defer m.mu.RUnlock()
	groupsArray := make([]*model.Group, 0, len(m.Groups))
	for id, g := range m.Groups {
		newGroup := model.Group{
			GroupID:     id,
			Name:        g.Name,
			Description: g.Description,
			CreatedBy:   g.CreatedBy,
			CreatedOn:   g.CreatedOn,
			ModifiedBy:  g.ModifiedBy,
			ModifiedOn:  g.ModifiedOn,
			OrderBy:     g.OrderBy,
			IsActive:    g.IsActive,
		}
		groupsArray = append(groupsArray, &newGroup)
	}
	return groupsArray
}
func (m *CGroupsMap) GetGroupById(id string) model.Group {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Groups[id]
}

//AREGLAR DESPUES
