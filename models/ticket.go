package models

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	UUID      string         `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Category string `gorm:"not null" json:"category"` // Ex: "Matériel", "Logiciel", "Réseau", etc.

	Probleme string `gorm:"not null" json:"probleme"` // Description détaillée du problème rencontré par le demandeur

	Statut string `gorm:"not null" json:"statut"` // Ex: "Ouvert", "En cours", "Résolu", "Fermé"

	TempsResolution string `json:"temps_resolution"` // Estimation du temps nécessaire pour résoudre le problème

	// Partie demandeur
	DemandeurUUID string    `gorm:"type:varchar(255);not null" json:"demandeur_uuid"`
	Demandeur     Demandeur `gorm:"foreignKey:DemandeurUUID" json:"demandeur"`

	// Partie affectation directement dans le burauu concerneé selon les prerogatives de chaque bureau
	BureauUUID string `gorm:"type:varchar(255);not null" json:"bureau_uuid"`
	Bureau     Bureau `gorm:"foreignKey:BureauUUID" json:"bureau"`
}
