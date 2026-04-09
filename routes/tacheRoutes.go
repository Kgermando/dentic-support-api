package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/controlleurs/taches"
)

// SetupTacheRoutes configure les routes pour les tâches
func SetupTacheRoutes(api fiber.Router) {
	t := api.Group("/taches")
	t.Get("/all/paginate", taches.GetPaginatedTaches)
	t.Get("/all", taches.GetAllTaches)
	t.Get("/by-ticket/:ticket_uuid", taches.GetTachesByTicket)
	t.Get("/by-agent/:agent_uuid", taches.GetTachesByAgent)
	t.Post("/create", taches.CreateTache)
	t.Get("/get/:uuid", taches.GetTache)
	t.Put("/update/:uuid", taches.UpdateTache)
	t.Patch("/statut/:uuid", taches.UpdateTacheStatut)
	t.Delete("/delete/:uuid", taches.DeleteTache)
}
