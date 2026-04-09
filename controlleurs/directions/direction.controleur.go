package directions

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/database"
	"github.com/kgermando/dentic-support-api/models"
	"github.com/kgermando/dentic-support-api/utils"
)

// GetPaginatedDirections retourne les directions avec pagination
func GetPaginatedDirections(c *fiber.Ctx) error {
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

	var directions []models.Direction
	var totalRecords int64

	db.Model(&models.Direction{}).
		Where("name ILIKE ?", "%"+search+"%").
		Count(&totalRecords)

	err = db.
		Where("name ILIKE ?", "%"+search+"%").
		Preload("Bureau").
		Preload("Agents").
		Offset(offset).
		Limit(limit).
		Order("updated_at DESC").
		Find(&directions).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch directions",
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
		"message":    "Directions retrieved successfully",
		"data":       directions,
		"pagination": pagination,
	})
}

// GetAllDirections retourne toutes les directions
func GetAllDirections(c *fiber.Ctx) error {
	db := database.DB
	var directions []models.Direction
	db.Preload("Bureau").Preload("Agents").Find(&directions)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All directions",
		"data":    directions,
	})
}

// GetDirection retourne une direction par UUID
func GetDirection(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB
	var direction models.Direction
	db.Where("uuid = ?", uuid).
		Preload("Bureau").
		Preload("Agents").
		First(&direction)
	if direction.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Direction not found",
			"data":    nil,
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Direction found",
		"data":    direction,
	})
}

// CreateDirection crée une nouvelle direction
func CreateDirection(c *fiber.Ctx) error {
	p := &models.Direction{}
	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	if p.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le nom de la direction est requis",
		})
	}

	direction := &models.Direction{
		UUID:        utils.GenerateUUID(),
		Name:        p.Name,
		Description: p.Description,
	}

	if err := database.DB.Create(direction).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create direction",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Direction created successfully",
		"data":    direction,
	})
}

// UpdateDirection met à jour une direction
func UpdateDirection(c *fiber.Ctx) error {
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

	var direction models.Direction
	db.Where("uuid = ?", uuid).First(&direction)
	if direction.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Direction not found",
		})
	}

	direction.Name = input.Name
	direction.Description = input.Description

	db.Save(&direction)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Direction updated successfully",
		"data":    direction,
	})
}

// DeleteDirection supprime une direction
func DeleteDirection(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var direction models.Direction
	db.Where("uuid = ?", uuid).First(&direction)
	if direction.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Direction not found",
		})
	}

	db.Delete(&direction)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Direction deleted successfully",
		"data":    nil,
	})
}
