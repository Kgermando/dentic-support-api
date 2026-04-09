package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/controlleurs/bureaux"
)

// SetupBureauRoutes configure les routes pour les bureaux
func SetupBureauRoutes(api fiber.Router) {
	b := api.Group("/bureaux")
	b.Get("/all/paginate", bureaux.GetPaginatedBureaux)
	b.Get("/all", bureaux.GetAllBureaux)
	b.Get("/by-direction/:direction_uuid", bureaux.GetBureauByDirection)
	b.Post("/create", bureaux.CreateBureau)
	b.Get("/get/:uuid", bureaux.GetBureau)
	b.Put("/update/:uuid", bureaux.UpdateBureau)
	b.Delete("/delete/:uuid", bureaux.DeleteBureau)
}
