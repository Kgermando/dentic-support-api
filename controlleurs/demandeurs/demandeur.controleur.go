package demandeurs

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/database"
	"github.com/kgermando/dentic-support-api/models"
	"github.com/kgermando/dentic-support-api/utils"
)

// GetPaginatedDemandeurs retourne les demandeurs avec pagination
func GetPaginatedDemandeurs(c *fiber.Ctx) error {
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

	var demandeurs []models.Demandeur
	var totalRecords int64

	db.Model(&models.Demandeur{}).
		Where("fullname ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.
		Where("fullname ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%").
		Preload("DirectionDemandeur").
		Preload("BureauDemandeur").
		Offset(offset).
		Limit(limit).
		Order("updated_at DESC").
		Find(&demandeurs).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch demandeurs",
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
		"message":    "Demandeurs retrieved successfully",
		"data":       demandeurs,
		"pagination": pagination,
	})
}

// GetAllDemandeurs retourne tous les demandeurs
func GetAllDemandeurs(c *fiber.Ctx) error {
	db := database.DB
	var demandeurs []models.Demandeur
	db.Preload("DirectionDemandeur").Preload("BureauDemandeur").Find(&demandeurs)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All demandeurs",
		"data":    demandeurs,
	})
}

// GetDemandeur retourne un demandeur par UUID
func GetDemandeur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB
	var demandeur models.Demandeur
	db.Where("uuid = ?", uuid).
		Preload("DirectionDemandeur").
		Preload("BureauDemandeur").
		First(&demandeur)
	if demandeur.Fullname == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Demandeur not found",
			"data":    nil,
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Demandeur found",
		"data":    demandeur,
	})
}

// CreateDemandeur crée un nouveau demandeur
func CreateDemandeur(c *fiber.Ctx) error {
	p := &models.Demandeur{}
	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	if p.Fullname == "" || p.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le nom et l'email du demandeur sont requis",
		})
	}

	demandeur := &models.Demandeur{
		UUID:                   utils.GenerateUUID(),
		Fullname:               p.Fullname,
		Email:                  p.Email,
		Telephone:              p.Telephone,
		Site:                   p.Site,
		DirectionDemandeurUUID: p.DirectionDemandeurUUID,
		BureauDemandeurUUID:    p.BureauDemandeurUUID,
	}

	if err := database.DB.Create(demandeur).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create demandeur",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Demandeur created successfully",
		"data":    demandeur,
	})
}

// UpdateDemandeur met à jour un demandeur
func UpdateDemandeur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateInput struct {
		Fullname               string `json:"fullname"`
		Email                  string `json:"email"`
		Telephone              string `json:"telephone"`
		Site                   string `json:"site"`
		DirectionDemandeurUUID string `json:"direction_demandeur_uuid"`
		BureauDemandeurUUID    string `json:"bureau_demandeur_uuid"`
	}

	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
		})
	}

	var demandeur models.Demandeur
	db.Where("uuid = ?", uuid).First(&demandeur)
	if demandeur.Fullname == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Demandeur not found",
		})
	}

	demandeur.Fullname = input.Fullname
	demandeur.Email = input.Email
	demandeur.Telephone = input.Telephone
	demandeur.Site = input.Site
	demandeur.DirectionDemandeurUUID = input.DirectionDemandeurUUID
	demandeur.BureauDemandeurUUID = input.BureauDemandeurUUID

	db.Save(&demandeur)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Demandeur updated successfully",
		"data":    demandeur,
	})
}

// DeleteDemandeur supprime un demandeur
func DeleteDemandeur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var demandeur models.Demandeur
	db.Where("uuid = ?", uuid).First(&demandeur)
	if demandeur.Fullname == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Demandeur not found",
		})
	}

	db.Delete(&demandeur)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Demandeur deleted successfully",
		"data":    nil,
	})
}
