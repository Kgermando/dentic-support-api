package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kgermando/dentic-support-api/middlewares"
)

// Setup configure toutes les routes de l'application
func Setup(app *fiber.App) {
	// Groupe API principal avec middleware logger
	api := app.Group("/api", logger.New())

	// Routes publiques (sans authentification)
	SetupAuthRoutes(api)

	// Routes protégées — nécessitent un token JWT valide
	protected := api.Group("", middlewares.IsAuthenticated)
	SetupAgentRoutes(protected)              // Routes agents utilisateurs
	SetupDirectionRoutes(protected)          // Routes directions
	SetupBureauRoutes(protected)             // Routes bureaux
	SetupDemandeurRoutes(protected)          // Routes demandeurs
	SetupDirectionDemandeurRoutes(protected) // Routes directions demandeurs
	SetupBureauDemandeurRoutes(protected)    // Routes bureaux demandeurs
	SetupTicketRoutes(protected)             // Routes tickets
	SetupTacheRoutes(protected)              // Routes tâches
	SetupTeamRoutes(protected)               // Routes équipes
	SetupDashboardRoutes(protected)          // Dashboard (filtré par rôle)
}
