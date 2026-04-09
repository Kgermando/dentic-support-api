package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/controlleurs/teams"
)

// SetupTeamRoutes configure les routes pour les équipes
func SetupTeamRoutes(api fiber.Router) {
	// Équipes
	t := api.Group("/teams")
	t.Get("/all/paginate", teams.GetPaginatedTeams)
	t.Get("/all", teams.GetAllTeams)
	t.Post("/create", teams.CreateTeam)
	t.Get("/get/:uuid", teams.GetTeam)
	t.Put("/update/:uuid", teams.UpdateTeam)
	t.Delete("/delete/:uuid", teams.DeleteTeam)

	// Team joins (membres)
	tj := api.Group("/team-joins")
	tj.Get("/by-team/:team_uuid", teams.GetTeamJoinsByTeam)
	tj.Post("/add", teams.AddTeamJoin)
	tj.Delete("/remove/:uuid", teams.RemoveTeamJoin)
}
