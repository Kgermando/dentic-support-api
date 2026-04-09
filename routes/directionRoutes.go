package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/controlleurs/directions"
)

// SetupDirectionRoutes configure les routes pour les directions
func SetupDirectionRoutes(api fiber.Router) {
	d := api.Group("/directions")
	d.Get("/all/paginate", directions.GetPaginatedDirections)
	d.Get("/all", directions.GetAllDirections)
	d.Post("/create", directions.CreateDirection)
	d.Get("/get/:uuid", directions.GetDirection)
	d.Put("/update/:uuid", directions.UpdateDirection)
	d.Delete("/delete/:uuid", directions.DeleteDirection)
}
