package bureaux

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/database"
	"github.com/kgermando/dentic-support-api/models"
	"github.com/kgermando/dentic-support-api/utils"
)

// GetPaginatedBureaux retourne les bureaux avec pagination
func GetPaginatedBureaux(c *fiber.Ctx) error {
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

	var bureaux []models.Bureau
	var totalRecords int64

	db.Model(&models.Bureau{}).
		Where("name ILIKE ?", "%"+search+"%").
		Count(&totalRecords)

	err = db.
		Where("name ILIKE ?", "%"+search+"%").
		Preload("Direction").
		Preload("Agents").
		Preload("TeamJoins").
		Offset(offset).
		Limit(limit).
		Order("updated_at DESC").
		Find(&bureaux).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch bureaux",
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
		"message":    "Bureaux retrieved successfully",
		"data":       bureaux,
		"pagination": pagination,
	})
}

// GetAllBureaux retourne tous les bureaux
func GetAllBureaux(c *fiber.Ctx) error {
	db := database.DB
	var bureaux []models.Bureau
	db.Preload("Direction").Preload("Agents").Preload("TeamJoins").Find(&bureaux)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All bureaux",
		"data":    bureaux,
	})
}

// GetBureauByDirection retourne les bureaux d'une direction
func GetBureauByDirection(c *fiber.Ctx) error {
	directionUUID := c.Params("direction_uuid")
	db := database.DB
	var bureaux []models.Bureau
	db.Where("direction_uuid = ?", directionUUID).
		Preload("Direction").
		Preload("Agents").
		Find(&bureaux)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Bureaux by direction",
		"data":    bureaux,
	})
}

// GetBureau retourne un bureau par UUID
func GetBureau(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB
	var bureau models.Bureau
	db.Where("uuid = ?", uuid).
		Preload("Direction").
		Preload("Agents").
		Preload("TeamJoins").
		First(&bureau)
	if bureau.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Bureau not found",
			"data":    nil,
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Bureau found",
		"data":    bureau,
	})
}

// CreateBureau crée un nouveau bureau
func CreateBureau(c *fiber.Ctx) error {
	p := &models.Bureau{}
	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	if p.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le nom du bureau est requis",
		})
	}

	bureau := &models.Bureau{
		UUID:          utils.GenerateUUID(),
		Name:          p.Name,
		Description:   p.Description,
		DirectionUUID: p.DirectionUUID,
	}

	if err := database.DB.Create(bureau).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create bureau",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Bureau created successfully",
		"data":    bureau,
	})
}

// UpdateBureau met à jour un bureau
func UpdateBureau(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateInput struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		DirectionUUID string `json:"direction_uuid"`
	}

	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
		})
	}

	var bureau models.Bureau
	db.Where("uuid = ?", uuid).First(&bureau)
	if bureau.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Bureau not found",
		})
	}

	bureau.Name = input.Name
	bureau.Description = input.Description
	bureau.DirectionUUID = input.DirectionUUID

	db.Save(&bureau)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Bureau updated successfully",
		"data":    bureau,
	})
}

// DeleteBureau supprime un bureau
func DeleteBureau(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var bureau models.Bureau
	db.Where("uuid = ?", uuid).First(&bureau)
	if bureau.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Bureau not found",
		})
	}

	db.Delete(&bureau)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Bureau deleted successfully",
		"data":    nil,
	})
}
