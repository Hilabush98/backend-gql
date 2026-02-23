package cache

import (
	"backend-gql/graph/model"
	"database/sql"
	"log"
	"sync"
)

var GlobalMatrixcTool *CToolsMap

type CToolsMap struct {
	tools map[string]model.Tool
	mu    sync.RWMutex
}

func InitToolsMap() *CToolsMap {
	return &CToolsMap{
		tools: make(map[string]model.Tool),
	}
}
func (p *CToolsMap) LoadFromDB(db *sql.DB) error {

	query := `SELECT TOOL_ID, TOOL_ID_FATHER, NAME, DESCRIPTION, ORDER_BY, CREATED_ON, CREATED_BY, MODIFIED_ON, MODIFIED_BY, IS_ACTIVE, CONFIGURATION, PATH FROM C_TOOLS WHERE IS_ACTIVE =1 `

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	newTools := make(map[string]model.Tool)

	for rows.Next() {
		// Instancia del modelo Tool (GQLGen genera este struct)
		t := model.Tool{}

		// 3. Scan exacto según el orden del SELECT
		// Nota: Asegúrate de que los tipos en model.Tool coincidan (string, int, etc)

		err := rows.Scan(
			&t.ToolID,
			&t.ToolIDFather,
			&t.Name,
			&t.Description,
			&t.OrderBy,
			&t.CreatedOn,
			&t.CreatedBy,
			&t.ModifiedOn,
			&t.ModifiedBy,
			&t.IsActive,
			&t.Configuration,
			&t.Path,
		)
		if err != nil {
			return err
		}

		newTools[t.ToolID] = t
	}

	p.mu.Lock()
	p.tools = newTools
	p.mu.Unlock()

	return nil

}

func (p *CToolsMap) PrintMapTool() {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for id, t := range p.tools {

		log.Println("ID: ", id)
		log.Println("Tool_id_Father: ", Str(t.ToolIDFather))
		log.Println("Name: ", *t.Name)
		log.Println("Description: ", Str(t.Description))
		log.Println("OrderBy: ", *t.OrderBy)
		log.Println("CreatedOn: ", Str(t.CreatedOn))
		log.Println("CreatedBy:", Str(t.CreatedBy))
		log.Println("ModifiedBy: ", Str(t.ModifiedBy))
		log.Println("Is_Active: ", *t.IsActive)
		log.Println("Configuration: ", *t.Configuration)
		log.Println("Path: ", Str(t.Path))

		log.Println("----------------------------")

	}
}

func (m *CToolsMap) LookIncToolsById(id string) model.Tool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tools[id]
}

//AREGLAR DESPUES
