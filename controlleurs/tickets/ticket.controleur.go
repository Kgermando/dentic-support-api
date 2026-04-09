package tickets

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/database"
	"github.com/kgermando/dentic-support-api/models"
	"github.com/kgermando/dentic-support-api/utils"
)

// GetPaginatedTickets retourne les tickets avec pagination
func GetPaginatedTickets(c *fiber.Ctx) error {
	db := database.DB

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit", "15"))
	if err != nil || limit <= 0 {
		limit = 15
	}
	offset := (page - 1) * limit
	search := c.Query("search", "")

	var tickets []models.Ticket
	var totalRecords int64

	db.Model(&models.Ticket{}).
		Where("category ILIKE ? OR statut ILIKE ? OR probleme ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.
		Where("category ILIKE ? OR statut ILIKE ? OR probleme ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%").
		Preload("Demandeur").
		Preload("Bureau").
		Offset(offset).
		Limit(limit).
		Order("updated_at DESC").
		Find(&tickets).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch tickets",
			"error":   err.Error(),
		})
	}

	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))
	pagination := map[string]interface{}{
		"total_records": totalRecords,
		"total_pages":   totalPages,
		"current_page":  page,
		"page_size":     limit,
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "Tickets retrieved successfully",
		"data":       tickets,
		"pagination": pagination,
	})
}

// GetAllTickets retourne tous les tickets
func GetAllTickets(c *fiber.Ctx) error {
	db := database.DB
	var tickets []models.Ticket
	db.Preload("Demandeur").Preload("Bureau").Find(&tickets)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All tickets",
		"data":    tickets,
	})
}

// GetTicketsByBureau retourne les tickets d'un bureau
func GetTicketsByBureau(c *fiber.Ctx) error {
	bureauUUID := c.Params("bureau_uuid")
	db := database.DB
	var tickets []models.Ticket
	db.Where("bureau_uuid = ?", bureauUUID).
		Preload("Demandeur").
		Preload("Bureau").
		Order("updated_at DESC").
		Find(&tickets)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Tickets by bureau",
		"data":    tickets,
	})
}

// GetTicketsByDemandeur retourne les tickets d'un demandeur
func GetTicketsByDemandeur(c *fiber.Ctx) error {
	demandeurUUID := c.Params("demandeur_uuid")
	db := database.DB
	var tickets []models.Ticket
	db.Where("demandeur_uuid = ?", demandeurUUID).
		Preload("Demandeur").
		Preload("Bureau").
		Order("updated_at DESC").
		Find(&tickets)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Tickets by demandeur",
		"data":    tickets,
	})
}

// GetTicket retourne un ticket par UUID
func GetTicket(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB
	var ticket models.Ticket
	db.Where("uuid = ?", uuid).
		Preload("Demandeur").
		Preload("Bureau").
		First(&ticket)
	if ticket.Probleme == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Ticket not found",
			"data":    nil,
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Ticket found",
		"data":    ticket,
	})
}

// CreateTicket crée un nouveau ticket
func CreateTicket(c *fiber.Ctx) error {
	p := &models.Ticket{}
	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	if p.Probleme == "" || p.Category == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le problème et la catégorie sont requis",
		})
	}
	if p.DemandeurUUID == "" || p.BureauUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le demandeur et le bureau sont requis",
		})
	}

	ticket := &models.Ticket{
		UUID:            utils.GenerateUUID(),
		Category:        p.Category,
		Probleme:        p.Probleme,
		Statut:          "Ouvert",
		TempsResolution: p.TempsResolution,
		DemandeurUUID:   p.DemandeurUUID,
		BureauUUID:      p.BureauUUID,
	}

	if err := database.DB.Create(ticket).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create ticket",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Ticket created successfully",
		"data":    ticket,
	})
}

// UpdateTicket met à jour un ticket
func UpdateTicket(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateInput struct {
		Category        string `json:"category"`
		Probleme        string `json:"probleme"`
		Statut          string `json:"statut"`
		TempsResolution string `json:"temps_resolution"`
		DemandeurUUID   string `json:"demandeur_uuid"`
		BureauUUID      string `json:"bureau_uuid"`
	}

	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
		})
	}

	var ticket models.Ticket
	db.Where("uuid = ?", uuid).First(&ticket)
	if ticket.Probleme == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Ticket not found",
		})
	}

	ticket.Category = input.Category
	ticket.Probleme = input.Probleme
	ticket.Statut = input.Statut
	ticket.TempsResolution = input.TempsResolution
	ticket.DemandeurUUID = input.DemandeurUUID
	ticket.BureauUUID = input.BureauUUID

	db.Save(&ticket)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Ticket updated successfully",
		"data":    ticket,
	})
}

// UpdateTicketStatut met à jour uniquement le statut d'un ticket
func UpdateTicketStatut(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type StatutInput struct {
		Statut string `json:"statut"`
	}

	var input StatutInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
		})
	}

	var ticket models.Ticket
	db.Where("uuid = ?", uuid).First(&ticket)
	if ticket.Probleme == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Ticket not found",
		})
	}

	ticket.Statut = input.Statut
	db.Save(&ticket)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Ticket statut updated successfully",
		"data":    ticket,
	})
}

// DeleteTicket supprime un ticket
func DeleteTicket(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var ticket models.Ticket
	db.Where("uuid = ?", uuid).First(&ticket)
	if ticket.Probleme == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Ticket not found",
		})
	}

	db.Delete(&ticket)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Ticket deleted successfully",
		"data":    nil,
	})
}
