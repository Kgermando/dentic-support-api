package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/controlleurs/dashboard"
)

// SetupDashboardRoutes configure la route unique du dashboard avec filtre par rôle
func SetupDashboardRoutes(api fiber.Router) {
	// Point d'entrée unique — le rôle est détecté automatiquement via le token JWT
	api.Get("/dashboard", dashboard.GetDashboard)
}
