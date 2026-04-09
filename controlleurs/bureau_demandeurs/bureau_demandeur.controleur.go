package bureau_demandeurs

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/database"
	"github.com/kgermando/dentic-support-api/models"
	"github.com/kgermando/dentic-support-api/utils"
)

// GetPaginatedBureauDemandeurs retourne les bureaux demandeurs avec pagination
func GetPaginatedBureauDemandeurs(c *fiber.Ctx) error {
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

	var bureaux []models.BureauDemandeur
	var totalRecords int64

	db.Model(&models.BureauDemandeur{}).
		Where("name ILIKE ?", "%"+search+"%").
		Count(&totalRecords)

	err = db.
		Where("name ILIKE ?", "%"+search+"%").
		Preload("DirectionDemandeur").
		Preload("Demandeurs").
		Offset(offset).
		Limit(limit).
		Order("updated_at DESC").
		Find(&bureaux).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch bureaux demandeurs",
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
		"message":    "Bureaux demandeurs retrieved successfully",
		"data":       bureaux,
		"pagination": pagination,
	})
}

// GetAllBureauDemandeurs retourne tous les bureaux demandeurs
func GetAllBureauDemandeurs(c *fiber.Ctx) error {
	db := database.DB
	var bureaux []models.BureauDemandeur
	db.Preload("DirectionDemandeur").Preload("Demandeurs").Find(&bureaux)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All bureaux demandeurs",
		"data":    bureaux,
	})
}

// GetBureauDemandeurByDirection retourne les bureaux demandeurs d'une direction
func GetBureauDemandeurByDirection(c *fiber.Ctx) error {
	directionUUID := c.Params("direction_uuid")
	db := database.DB
	var bureaux []models.BureauDemandeur
	db.Where("direction_demandeur_uuid = ?", directionUUID).
		Preload("DirectionDemandeur").
		Preload("Demandeurs").
		Find(&bureaux)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Bureaux demandeurs by direction",
		"data":    bureaux,
	})
}

// GetBureauDemandeur retourne un bureau demandeur par UUID
func GetBureauDemandeur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB
	var bureau models.BureauDemandeur
	db.Where("uuid = ?", uuid).
		Preload("DirectionDemandeur").
		Preload("Demandeurs").
		First(&bureau)
	if bureau.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Bureau demandeur not found",
			"data":    nil,
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Bureau demandeur found",
		"data":    bureau,
	})
}

// CreateBureauDemandeur crée un nouveau bureau demandeur
func CreateBureauDemandeur(c *fiber.Ctx) error {
	p := &models.BureauDemandeur{}
	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	if p.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le nom du bureau demandeur est requis",
		})
	}

	bureau := &models.BureauDemandeur{
		UUID:                   utils.GenerateUUID(),
		Name:                   p.Name,
		DirectionDemandeurUUID: p.DirectionDemandeurUUID,
	}

	if err := database.DB.Create(bureau).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create bureau demandeur",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Bureau demandeur created successfully",
		"data":    bureau,
	})
}

// UpdateBureauDemandeur met à jour un bureau demandeur
func UpdateBureauDemandeur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateInput struct {
		Name                   string `json:"name"`
		DirectionDemandeurUUID string `json:"direction_demandeur_uuid"`
	}

	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
		})
	}

	var bureau models.BureauDemandeur
	db.Where("uuid = ?", uuid).First(&bureau)
	if bureau.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Bureau demandeur not found",
		})
	}

	bureau.Name = input.Name
	bureau.DirectionDemandeurUUID = input.DirectionDemandeurUUID
	db.Save(&bureau)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Bureau demandeur updated successfully",
		"data":    bureau,
	})
}

// DeleteBureauDemandeur supprime un bureau demandeur
func DeleteBureauDemandeur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var bureau models.BureauDemandeur
	db.Where("uuid = ?", uuid).First(&bureau)
	if bureau.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Bureau demandeur not found",
		})
	}

	db.Delete(&bureau)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Bureau demandeur deleted successfully",
		"data":    nil,
	})
}
