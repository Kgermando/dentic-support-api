package routes

import (
	"github.com/gofiber/fiber/v2"
	bureau_demandeurs "github.com/kgermando/dentic-support-api/controlleurs/bureau_demandeurs"
	"github.com/kgermando/dentic-support-api/controlleurs/demandeurs"
	direction_demandeurs "github.com/kgermando/dentic-support-api/controlleurs/direction_demandeurs"
)

// SetupDemandeurRoutes configure les routes pour les demandeurs et leurs structures
func SetupDemandeurRoutes(api fiber.Router) {
	// Direction demandeurs
	dd := api.Group("/direction-demandeurs")
	dd.Get("/all/paginate", direction_demandeurs.GetPaginatedDirectionDemandeurs)
	dd.Get("/all", direction_demandeurs.GetAllDirectionDemandeurs)
	dd.Post("/create", direction_demandeurs.CreateDirectionDemandeur)
	dd.Get("/get/:uuid", direction_demandeurs.GetDirectionDemandeur)
	dd.Put("/update/:uuid", direction_demandeurs.UpdateDirectionDemandeur)
	dd.Delete("/delete/:uuid", direction_demandeurs.DeleteDirectionDemandeur)

	// Bureau demandeurs
	bd := api.Group("/bureau-demandeurs")
	bd.Get("/all/paginate", bureau_demandeurs.GetPaginatedBureauDemandeurs)
	bd.Get("/all", bureau_demandeurs.GetAllBureauDemandeurs)
	bd.Get("/by-direction/:direction_uuid", bureau_demandeurs.GetBureauDemandeurByDirection)
	bd.Post("/create", bureau_demandeurs.CreateBureauDemandeur)
	bd.Get("/get/:uuid", bureau_demandeurs.GetBureauDemandeur)
	bd.Put("/update/:uuid", bureau_demandeurs.UpdateBureauDemandeur)
	bd.Delete("/delete/:uuid", bureau_demandeurs.DeleteBureauDemandeur)

	// Demandeurs
	dem := api.Group("/demandeurs")
	dem.Get("/all/paginate", demandeurs.GetPaginatedDemandeurs)
	dem.Get("/all", demandeurs.GetAllDemandeurs)
	dem.Post("/create", demandeurs.CreateDemandeur)
	dem.Get("/get/:uuid", demandeurs.GetDemandeur)
	dem.Put("/update/:uuid", demandeurs.UpdateDemandeur)
	dem.Delete("/delete/:uuid", demandeurs.DeleteDemandeur)
}
