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

	query := `SELECT TOOL_ID, TOOL_ID_FATHER, NAME, DESCRIPTION, ORDER_BY, CREATED_ON, CREATED_BY, MODIFIED_ON, MODIFIED_BY, IS_ACTIVE, CONFIGURATION, PATH FROM APP.C_TOOLS WHERE IS_ACTIVE =true `

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	newTools := make(map[string]model.Tool)

	for rows.Next() {
		var toolId, configuration string
		var createdOn, name, path, orderBy, createdBy, modifiedOn, modifiedBy, description, toolIdFather *string
		var isActive bool
		err := rows.Scan(
			&toolId,
			&toolIdFather,
			&name,
			&description,
			&orderBy,
			&createdOn,
			&createdBy,
			&modifiedOn,
			&modifiedBy,
			&isActive,
			&configuration,
			&path,
		)
		if err != nil {
			return err
		}
		t := model.Tool{
			ToolID:        toolId,
			ToolIDFather:  NilValidate(toolIdFather),
			Name:          NilValidate(name),
			Description:   NilValidate(description),
			OrderBy:       NilValidate(orderBy),
			CreatedOn:     NilValidate(createdOn),
			CreatedBy:     NilValidate(createdBy),
			ModifiedOn:    NilValidate(modifiedOn),
			ModifiedBy:    NilValidate(modifiedBy),
			IsActive:      &isActive,
			Configuration: &configuration,
			Path:          NilValidate(path),
		}

		newTools[toolId] = t
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

func (m *CToolsMap) GetToolsById(id string) model.Tool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tools[id]
}
func (t *CToolsMap) GetAllTools() []*model.Tool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	toolArray := make([]*model.Tool, 0, len(t.tools))
	for id, g := range t.tools {
		newTool := model.Tool{
			ToolID:        id,
			ToolIDFather:  g.ToolIDFather,
			Name:          g.Name,
			Description:   g.Description,
			OrderBy:       g.OrderBy,
			CreatedOn:     g.CreatedOn,
			CreatedBy:     g.CreatedBy,
			ModifiedOn:    g.ModifiedOn,
			ModifiedBy:    g.ModifiedBy,
			IsActive:      g.IsActive,
			Configuration: g.Configuration,
			Path:          g.Path,
		}
		toolArray = append(toolArray, &newTool)
	}
	return toolArray

}

//AREGLAR DESPUES
