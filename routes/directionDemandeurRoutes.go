package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/controlleurs/direction_demandeurs"
)

// SetupDirectionDemandeurRoutes configure les routes pour les directions demandeurs
func SetupDirectionDemandeurRoutes(api fiber.Router) {
	d := api.Group("/direction-demandeurs")
	d.Get("/all/paginate", direction_demandeurs.GetPaginatedDirectionDemandeurs)
	d.Get("/all", direction_demandeurs.GetAllDirectionDemandeurs)
	d.Post("/create", direction_demandeurs.CreateDirectionDemandeur)
	d.Get("/get/:uuid", direction_demandeurs.GetDirectionDemandeur)
	d.Put("/update/:uuid", direction_demandeurs.UpdateDirectionDemandeur)
	d.Delete("/delete/:uuid", direction_demandeurs.DeleteDirectionDemandeur)
}
