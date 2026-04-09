package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/controlleurs/agents"
)

// SetupAgentRoutes configure les routes agents 
func SetupAgentRoutes(api fiber.Router) {
	// ============================================================
	// AGENTS ROUTES
	// ============================================================
	ag := api.Group("/agents")
	ag.Get("/all/paginate", agents.GetPaginatedAgents)
	ag.Get("/all", agents.GetAllAgents) 
	ag.Post("/create", agents.CreateAgent)
	ag.Get("/get/:uuid", agents.GetAgent)
	ag.Put("/update/:uuid", agents.UpdateAgent)
	ag.Delete("/delete/:uuid", agents.DeleteAgent)

}
