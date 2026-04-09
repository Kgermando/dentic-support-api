package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Setup configure toutes les routes de l'application
func Setup(app *fiber.App) {
	// Groupe API principal avec middleware logger
	api := app.Group("/api", logger.New())

	// Configuration de toutes les sous-routes par domaine
	SetupAuthRoutes(api)               // Routes d'authentification
	SetupAgentRoutes(api)              // Routes agents utilisateurs
	SetupDirectionRoutes(api)          // Routes directions
	SetupBureauRoutes(api)             // Routes bureaux
	SetupDemandeurRoutes(api)          // Routes demandeurs
	SetupDirectionDemandeurRoutes(api) // Routes directions demandeurs
	SetupBureauDemandeurRoutes(api)    // Routes bureaux demandeurs
	SetupTicketRoutes(api)             // Routes tickets
	SetupTacheRoutes(api)              // Routes tâches
	SetupTeamRoutes(api)               // Routes équipes
	SetupDashboardRoutes(api)          // Dashboard (filtré par rôle)
}
