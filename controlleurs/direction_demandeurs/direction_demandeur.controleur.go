package direction_demandeurs

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/database"
	"github.com/kgermando/dentic-support-api/models"
	"github.com/kgermando/dentic-support-api/utils"
)

// GetPaginatedDirectionDemandeurs retourne les directions demandeurs avec pagination
func GetPaginatedDirectionDemandeurs(c *fiber.Ctx) error {
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

	var directions []models.DirectionDemandeur
	var totalRecords int64

	db.Model(&models.DirectionDemandeur{}).
		Where("name ILIKE ?", "%"+search+"%").
		Count(&totalRecords)

	err = db.
		Where("name ILIKE ?", "%"+search+"%").
		Preload("BureauDemandeurs").
		Preload("Demandeurs").
		Offset(offset).
		Limit(limit).
		Order("updated_at DESC").
		Find(&directions).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch directions demandeurs",
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
		"message":    "Directions demandeurs retrieved successfully",
		"data":       directions,
		"pagination": pagination,
	})
}

// GetAllDirectionDemandeurs retourne toutes les directions demandeurs
func GetAllDirectionDemandeurs(c *fiber.Ctx) error {
	db := database.DB
	var directions []models.DirectionDemandeur
	db.Preload("BureauDemandeurs").Preload("Demandeurs").Find(&directions)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All directions demandeurs",
		"data":    directions,
	})
}

// GetDirectionDemandeur retourne une direction demandeur par UUID
func GetDirectionDemandeur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB
	var direction models.DirectionDemandeur
	db.Where("uuid = ?", uuid).
		Preload("BureauDemandeurs").
		Preload("Demandeurs").
		First(&direction)
	if direction.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Direction demandeur not found",
			"data":    nil,
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Direction demandeur found",
		"data":    direction,
	})
}

// CreateDirectionDemandeur crée une nouvelle direction demandeur
func CreateDirectionDemandeur(c *fiber.Ctx) error {
	p := &models.DirectionDemandeur{}
	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	if p.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le nom de la direction demandeur est requis",
		})
	}

	direction := &models.DirectionDemandeur{
		UUID: utils.GenerateUUID(),
		Name: p.Name,
	}

	if err := database.DB.Create(direction).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create direction demandeur",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Direction demandeur created successfully",
		"data":    direction,
	})
}

// UpdateDirectionDemandeur met à jour une direction demandeur
func UpdateDirectionDemandeur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateInput struct {
		Name string `json:"name"`
	}

	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
		})
	}

	var direction models.DirectionDemandeur
	db.Where("uuid = ?", uuid).First(&direction)
	if direction.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Direction demandeur not found",
		})
	}

	direction.Name = input.Name
	db.Save(&direction)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Direction demandeur updated successfully",
		"data":    direction,
	})
}

// DeleteDirectionDemandeur supprime une direction demandeur
func DeleteDirectionDemandeur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var direction models.DirectionDemandeur
	db.Where("uuid = ?", uuid).First(&direction)
	if direction.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Direction demandeur not found",
		})
	}

	db.Delete(&direction)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Direction demandeur deleted successfully",
		"data":    nil,
	})
}
