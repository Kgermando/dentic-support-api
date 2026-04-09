package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/controlleurs/tickets"
)

// SetupTicketRoutes configure les routes pour les tickets
func SetupTicketRoutes(api fiber.Router) {
	t := api.Group("/tickets")
	t.Get("/all/paginate", tickets.GetPaginatedTickets)
	t.Get("/all", tickets.GetAllTickets)
	t.Get("/by-bureau/:bureau_uuid", tickets.GetTicketsByBureau)
	t.Get("/by-demandeur/:demandeur_uuid", tickets.GetTicketsByDemandeur)
	t.Post("/create", tickets.CreateTicket)
	t.Get("/get/:uuid", tickets.GetTicket)
	t.Put("/update/:uuid", tickets.UpdateTicket)
	t.Patch("/statut/:uuid", tickets.UpdateTicketStatut)
	t.Delete("/delete/:uuid", tickets.DeleteTicket)
}
