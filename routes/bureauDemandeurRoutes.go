package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/controlleurs/bureau_demandeurs"
)

// SetupBureauDemandeurRoutes configure les routes pour les bureaux demandeurs
func SetupBureauDemandeurRoutes(api fiber.Router) {
	b := api.Group("/bureau-demandeurs")
	b.Get("/all/paginate", bureau_demandeurs.GetPaginatedBureauDemandeurs)
	b.Get("/all", bureau_demandeurs.GetAllBureauDemandeurs)
	b.Get("/by-direction/:direction_uuid", bureau_demandeurs.GetBureauDemandeurByDirection)
	b.Post("/create", bureau_demandeurs.CreateBureauDemandeur)
	b.Get("/get/:uuid", bureau_demandeurs.GetBureauDemandeur)
	b.Put("/update/:uuid", bureau_demandeurs.UpdateBureauDemandeur)
	b.Delete("/delete/:uuid", bureau_demandeurs.DeleteBureauDemandeur)
}
