package taches

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/database"
	"github.com/kgermando/dentic-support-api/models"
	"github.com/kgermando/dentic-support-api/utils"
)

// GetPaginatedTaches retourne les tâches avec pagination
func GetPaginatedTaches(c *fiber.Ctx) error {
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

	var taches []models.Tache
	var totalRecords int64

	db.Model(&models.Tache{}).
		Where("name ILIKE ? OR statut ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.
		Where("name ILIKE ? OR statut ILIKE ?", "%"+search+"%", "%"+search+"%").
		Preload("Agent").
		Offset(offset).
		Limit(limit).
		Order("updated_at DESC").
		Find(&taches).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch tâches",
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
		"message":    "Tâches retrieved successfully",
		"data":       taches,
		"pagination": pagination,
	})
}

// GetAllTaches retourne toutes les tâches
func GetAllTaches(c *fiber.Ctx) error {
	db := database.DB
	var taches []models.Tache
	db.Preload("Agent").Find(&taches)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All tâches",
		"data":    taches,
	})
}

// GetTachesByTicket retourne les tâches d'un ticket
func GetTachesByTicket(c *fiber.Ctx) error {
	ticketUUID := c.Params("ticket_uuid")
	db := database.DB
	var taches []models.Tache
	db.Where("ticket_uuid = ?", ticketUUID).
		Preload("Agent").
		Order("updated_at DESC").
		Find(&taches)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Tâches by ticket",
		"data":    taches,
	})
}

// GetTachesByAgent retourne les tâches d'un agent
func GetTachesByAgent(c *fiber.Ctx) error {
	agentUUID := c.Params("agent_uuid")
	db := database.DB
	var taches []models.Tache
	db.Where("agent_uuid = ?", agentUUID).
		Preload("Agent").
		Order("updated_at DESC").
		Find(&taches)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Tâches by agent",
		"data":    taches,
	})
}

// GetTache retourne une tâche par UUID
func GetTache(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB
	var tache models.Tache
	db.Where("uuid = ?", uuid).
		Preload("Agent").
		First(&tache)
	if tache.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Tâche not found",
			"data":    nil,
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Tâche found",
		"data":    tache,
	})
}

// CreateTache crée une nouvelle tâche
func CreateTache(c *fiber.Ctx) error {
	p := &models.Tache{}
	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	if p.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le nom de la tâche est requis",
		})
	}
	if p.AgentUUID == "" || p.TicketUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "L'agent et le ticket sont requis",
		})
	}

	tache := &models.Tache{
		UUID:        utils.GenerateUUID(),
		Name:        p.Name,
		Description: p.Description,
		Statut:      "En attente",
		AgentUUID:   p.AgentUUID,
		TicketUUID:  p.TicketUUID,
	}

	if err := database.DB.Create(tache).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create tâche",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Tâche created successfully",
		"data":    tache,
	})
}

// UpdateTache met à jour une tâche
func UpdateTache(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateInput struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Statut      string `json:"statut"`
		AgentUUID   string `json:"agent_uuid"`
		TicketUUID  string `json:"ticket_uuid"`
	}

	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
		})
	}

	var tache models.Tache
	db.Where("uuid = ?", uuid).First(&tache)
	if tache.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Tâche not found",
		})
	}

	tache.Name = input.Name
	tache.Description = input.Description
	tache.Statut = input.Statut
	tache.AgentUUID = input.AgentUUID
	tache.TicketUUID = input.TicketUUID

	db.Save(&tache)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Tâche updated successfully",
		"data":    tache,
	})
}

// UpdateTacheStatut met à jour uniquement le statut d'une tâche
func UpdateTacheStatut(c *fiber.Ctx) error {
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

	var tache models.Tache
	db.Where("uuid = ?", uuid).First(&tache)
	if tache.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Tâche not found",
		})
	}

	tache.Statut = input.Statut
	db.Save(&tache)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Statut de la tâche updated successfully",
		"data":    tache,
	})
}

// DeleteTache supprime une tâche
func DeleteTache(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var tache models.Tache
	db.Where("uuid = ?", uuid).First(&tache)
	if tache.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Tâche not found",
		})
	}

	db.Delete(&tache)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Tâche deleted successfully",
		"data":    nil,
	})
}
