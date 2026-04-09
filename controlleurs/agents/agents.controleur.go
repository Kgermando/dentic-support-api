package agents

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/database"
	"github.com/kgermando/dentic-support-api/models"
	"github.com/kgermando/dentic-support-api/utils"
)

// Paginate
func GetPaginatedAgents(c *fiber.Ctx) error {
	db := database.DB

	// Parse query parameters for pagination
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit", "15"))
	if err != nil || limit <= 0 {
		limit = 15
	}
	offset := (page - 1) * limit

	// Parse search query
	search := c.Query("search", "")

	var agents []models.Agent
	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Agent{}).
		Where("fullname ILIKE ? OR role ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.
		Where("fullname ILIKE ? OR role ILIKE ?", "%"+search+"%", "%"+search+"%").
		Preload("Direction").
		Preload("Bureau").
		Offset(offset).
		Limit(limit).
		Order("agents.updated_at DESC").
		Find(&agents).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch Agents",
			"error":   err.Error(),
		})
	}

	// Calculate total pages
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))

	//  Prepare pagination metadata
	pagination := map[string]interface{}{
		"total_records": totalRecords,
		"total_pages":   totalPages,
		"current_page":  page,
		"page_size":     limit,
	}

	// Return response
	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "Agents retrieved successfully",
		"data":       agents,
		"pagination": pagination,
	})
}

// query all data
func GetAllAgents(c *fiber.Ctx) error {
	db := database.DB
	var agents []models.Agent
	db.Preload("Direction").Preload("Bureau").Find(&agents)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All agents",
		"data":    agents,
	})
}

// Get one data
func GetAgent(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB
	var agent models.Agent
	db.Where("uuid = ?", uuid).Preload("Direction").Preload("Bureau").First(&agent)
	if agent.Fullname == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No Agent name found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Agent found",
			"data":    agent,
		},
	)
}

// Create data
func CreateAgent(c *fiber.Ctx) error {
	p := &models.Agent{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	if p.Fullname == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Form not complete",
				"data":    nil,
			},
		)
	}

	if p.Password != p.PasswordConfirm {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	agent := &models.Agent{
		Fullname:      p.Fullname,
		Email:         p.Email,
		Telephone:     p.Telephone,
		TranchAge:     p.TranchAge,
		Role:          p.Role,
		Permission:    p.Permission,
		Status:        p.Status,
		DirectionUUID: p.DirectionUUID,
		BureauUUID:    p.BureauUUID,
	}

	agent.SetPassword(p.Password)

	agent.UUID = utils.GenerateUUID()

	database.DB.Create(agent)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Agent Created success",
			"data":    agent,
		},
	)
}

// Update data
func UpdateAgent(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateDataInput struct {
		Fullname        string `gorm:"not null" json:"fullname"`
		Email           string `gorm:"unique; not null" json:"email"`
		Telephone       string `gorm:"unique; not null" json:"telephone"`
		TranchAge       string `json:"tranch_age"`
		Password        string `json:"password" validate:"required"`
		PasswordConfirm string `json:"password_confirm" gorm:"-"`
		Role            string `json:"role"`
		Permission      string `json:"permission"`
		Status          bool   `json:"status"`
		DirectionUUID   string `json:"direction_uuid"`
		BureauUUID      string `json:"bureau_uuid"`
	}

	var updateData UpdateDataInput

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"data":    nil,
			},
		)
	}

	agent := new(models.Agent)

	db.Where("uuid = ?", uuid).First(&agent)
	agent.Fullname = updateData.Fullname
	agent.Email = updateData.Email
	agent.Telephone = updateData.Telephone
	agent.TranchAge = updateData.TranchAge
	agent.Role = updateData.Role
	agent.Permission = updateData.Permission
	agent.Status = updateData.Status
	agent.DirectionUUID = updateData.DirectionUUID
	agent.BureauUUID = updateData.BureauUUID

	db.Save(&agent)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Agent updated success",
			"data":    agent,
		},
	)
}

// Delete data
func DeleteAgent(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var agent models.Agent
	db.Where("uuid = ?", uuid).First(&agent)
	if agent.Fullname == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No Agent name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&agent)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Agent deleted success",
			"data":    nil,
		},
	)
}
