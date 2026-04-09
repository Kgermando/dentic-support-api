package database

import (
	"fmt"
	"log"

	"github.com/kgermando/dentic-support-api/models"
	"github.com/kgermando/dentic-support-api/utils"
)

// SeedSuperAdmin crée un agent SuperAdmin uniquement si la table agents est vide.
func SeedSuperAdmin() {
	var count int64
	DB.Model(&models.Agent{}).Count(&count)
	if count > 0 {
		fmt.Println("SuperAdmin déjà existant, seed ignoré.")
		return
	}

	superAdmin := models.Agent{
		UUID:      utils.GenerateUUID(),
		Fullname:  "Super Administrateur",
		Email:     "superadmin@dentic.app",
		Telephone: "+000000000000",
		Role:      "SuperAdmin",
		Status:    true,
	}

	rawPassword := "SuperAdmin@2026!"
	superAdmin.SetPassword(rawPassword)

	if err := DB.Create(&superAdmin).Error; err != nil {
		log.Printf("Échec de la création du SuperAdmin : %v\n", err)
		return
	}

	fmt.Println("SuperAdmin créé avec succès ✅")
	fmt.Printf("  Email    : %s\n", superAdmin.Email)
	fmt.Printf("  Password : %s\n", rawPassword)
	fmt.Println("  ⚠️  Changez ce mot de passe après la première connexion !")
}
