package cache

import (
	"backend-gql/graph/model"
	"database/sql"
	"log"
	"sync"
)

var GlobalMatrixcProfiles *CProfilesMap

type CProfilesMap struct {
	Profiles map[string]model.Profile
	mu       sync.RWMutex
}

func InitProfilesMap() *CProfilesMap {
	return &CProfilesMap{
		Profiles: make(map[string]model.Profile),
	}
}
func (p *CProfilesMap) LoadFromDB(db *sql.DB) error {

	query := `SELECT Profile_ID, NAME, DESCRIPTION, CREATED_ON, CREATED_BY, MODIFIED_ON, MODIFIED_BY, ORDER_BY, IS_ACTIVE FROM c_Profiles WHERE IS_ACTIVE =1 `

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	newProfiles := make(map[string]model.Profile)

	for rows.Next() {
		var profileId string
		var createdOn, name, orderBy, createdBy, modifiedOn, modifiedBy, description *string
		var isActive bool
		err := rows.Scan(
			&profileId,
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

		t := model.Profile{
			ProfileID:   profileId,
			Name:        NilValidate(name),
			OrderBy:     NilValidate(orderBy),
			CreatedOn:   createdOn,
			CreatedBy:   NilValidate(createdBy),
			ModifiedOn:  modifiedOn,
			ModifiedBy:  NilValidate(modifiedBy),
			Description: NilValidate(description),
			IsActive:    &isActive,
		}
		newProfiles[t.ProfileID] = t
	}

	p.mu.Lock()
	p.Profiles = newProfiles
	p.mu.Unlock()

	return nil

}
func (p *CProfilesMap) GetAllProfiles() []*model.Profile {
	p.mu.RLock()
	defer p.mu.RUnlock()
	profileArray := make([]*model.Profile, 0, len(p.Profiles))
	for id, g := range p.Profiles {
		newGroup := model.Profile{
			ProfileID:   id,
			Name:        g.Name,
			Description: g.Description,
			CreatedBy:   g.CreatedBy,
			CreatedOn:   g.CreatedOn,
			ModifiedBy:  g.ModifiedBy,
			ModifiedOn:  g.ModifiedOn,
			OrderBy:     g.OrderBy,
			IsActive:    g.IsActive,
		}
		profileArray = append(profileArray, &newGroup)
	}
	return profileArray
}

func (p *CProfilesMap) PrintMapTool() {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for id, t := range p.Profiles {

		log.Println("ID: ", id)
		log.Println("Name: ", *t.Name)
		log.Println("Description: ", Str(t.Description))
		log.Println("OrderBy: ", *t.OrderBy)
		log.Println("CreatedOn: ", Str(t.CreatedOn))
		log.Println("CreatedBy:", Str(t.CreatedBy))
		log.Println("ModifiedBy: ", Str(t.ModifiedBy))
		log.Println("Configuration: ", *t.OrderBy)
		log.Println("Is_Active: ", *t.IsActive)

		log.Println("----------------------------")

	}
}

func (m *CProfilesMap) GetcProfilesById(id string) model.Profile {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Profiles[id]
}

//AREGLAR DESPUES
