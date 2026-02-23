package cache

import (
	"backend-gql/graph/model"
	"database/sql"
	"log"
	"sync"
)

var GlobalMatrixcProfile *CProfileMap

type CProfileMap struct {
	profiles map[string]model.Profile
	mu       sync.RWMutex
}

func InitProfilesMap() *CProfileMap {
	return &CProfileMap{
		profiles: make(map[string]model.Profile),
	}
}
func (p *CProfileMap) LoadFromDB(db *sql.DB) error {

	query := `SELECT PROFILE_ID, NAME, DESCRIPTION, CREATED_ON, CREATED_BY, MODIFIED_ON, MODIFIED_BY, ORDER_BY, IS_ACTIVE FROM c_profiles WHERE IS_ACTIVE =1 `

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	newProfiles := make(map[string]model.Profile)

	for rows.Next() {
		// Instancia del modelo Tool (GQLGen genera este struct)
		t := model.Profile{}

		// 3. Scan exacto según el orden del SELECT
		// Nota: Asegúrate de que los tipos en model.Tool coincidan (string, int, etc)

		err := rows.Scan(
			&t.ProfileID,
			&t.Name,
			&t.Description,
			&t.CreatedOn,
			&t.CreatedBy,
			&t.ModifiedOn,
			&t.ModifiedBy,
			&t.OrderBy,
			&t.IsActive,
		)
		if err != nil {
			return err
		}

		newProfiles[t.ProfileID] = t
	}

	p.mu.Lock()
	p.profiles = newProfiles
	p.mu.Unlock()

	return nil

}

func (p *CProfileMap) PrintMapTool() {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for id, t := range p.profiles {

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

func (m *CProfileMap) LookIncProfileById(id string) model.Profile {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.profiles[id]
}

//AREGLAR DESPUES
