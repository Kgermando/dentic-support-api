package routes

import (
	"github.com/gofiber/fiber/v2" 
	"github.com/kgermando/dentic-support-api/controlleurs/demandeurs" 
)

// SetupDemandeurRoutes configure les routes pour les demandeurs et leurs structures
func SetupDemandeurRoutes(api fiber.Router) { 
	// Demandeurs
	dem := api.Group("/demandeurs")
	dem.Get("/all/paginate", demandeurs.GetPaginatedDemandeurs)
	dem.Get("/all", demandeurs.GetAllDemandeurs)
	dem.Post("/create", demandeurs.CreateDemandeur)
	dem.Get("/get/:uuid", demandeurs.GetDemandeur)
	dem.Put("/update/:uuid", demandeurs.UpdateDemandeur)
	dem.Delete("/delete/:uuid", demandeurs.DeleteDemandeur)
}
