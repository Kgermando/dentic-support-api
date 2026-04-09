package teams

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/database"
	"github.com/kgermando/dentic-support-api/models"
	"github.com/kgermando/dentic-support-api/utils"
)

// GetPaginatedTeams retourne les équipes avec pagination
func GetPaginatedTeams(c *fiber.Ctx) error {
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

	var teams []models.Team
	var totalRecords int64

	db.Model(&models.Team{}).
		Where("name ILIKE ?", "%"+search+"%").
		Count(&totalRecords)

	err = db.
		Where("name ILIKE ?", "%"+search+"%").
		Preload("TeamJoins").
		Offset(offset).
		Limit(limit).
		Order("updated_at DESC").
		Find(&teams).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch équipes",
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
		"message":    "Équipes retrieved successfully",
		"data":       teams,
		"pagination": pagination,
	})
}

// GetAllTeams retourne toutes les équipes
func GetAllTeams(c *fiber.Ctx) error {
	db := database.DB
	var teams []models.Team
	db.Preload("TeamJoins").Find(&teams)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All équipes",
		"data":    teams,
	})
}

// GetTeam retourne une équipe par UUID
func GetTeam(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB
	var team models.Team
	db.Where("uuid = ?", uuid).
		Preload("TeamJoins.Agent").
		Preload("TeamJoins.Bureau").
		First(&team)
	if team.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Équipe not found",
			"data":    nil,
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Équipe found",
		"data":    team,
	})
}

// CreateTeam crée une nouvelle équipe
func CreateTeam(c *fiber.Ctx) error {
	p := &models.Team{}
	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	if p.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le nom de l'équipe est requis",
		})
	}

	team := &models.Team{
		UUID:        utils.GenerateUUID(),
		Name:        p.Name,
		Description: p.Description,
	}

	if err := database.DB.Create(team).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create équipe",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Équipe created successfully",
		"data":    team,
	})
}

// UpdateTeam met à jour une équipe
func UpdateTeam(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateInput struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
		})
	}

	var team models.Team
	db.Where("uuid = ?", uuid).First(&team)
	if team.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Équipe not found",
		})
	}

	team.Name = input.Name
	team.Description = input.Description
	db.Save(&team)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Équipe updated successfully",
		"data":    team,
	})
}

// DeleteTeam supprime une équipe
func DeleteTeam(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var team models.Team
	db.Where("uuid = ?", uuid).First(&team)
	if team.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Équipe not found",
		})
	}

	db.Delete(&team)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Équipe deleted successfully",
		"data":    nil,
	})
}

// ==================== TEAM JOINS ====================

// GetTeamJoinsByTeam retourne les membres d'une équipe
func GetTeamJoinsByTeam(c *fiber.Ctx) error {
	teamUUID := c.Params("team_uuid")
	db := database.DB
	var joins []models.TeamJoin
	db.Where("team_uuid = ?", teamUUID).
		Preload("Team").
		Preload("Agent").
		Preload("Bureau").
		Find(&joins)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Team members retrieved",
		"data":    joins,
	})
}

// AddTeamJoin ajoute un agent à une équipe
func AddTeamJoin(c *fiber.Ctx) error {
	p := &models.TeamJoin{}
	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	if p.TeamUUID == "" || p.AgentUUID == "" || p.BureauUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "L'équipe, l'agent et le bureau sont requis",
		})
	}

	join := &models.TeamJoin{
		UUID:       utils.GenerateUUID(),
		TeamUUID:   p.TeamUUID,
		AgentUUID:  p.AgentUUID,
		BureauUUID: p.BureauUUID,
	}

	if err := database.DB.Create(join).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to add team member",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Team member added successfully",
		"data":    join,
	})
}

// RemoveTeamJoin retire un agent d'une équipe
func RemoveTeamJoin(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var join models.TeamJoin
	db.Where("uuid = ?", uuid).First(&join)
	if join.UUID == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Team member not found",
		})
	}

	db.Delete(&join)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Team member removed successfully",
		"data":    nil,
	})
}
