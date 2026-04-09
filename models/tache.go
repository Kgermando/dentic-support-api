package models

import (
	"time"

	"gorm.io/gorm"
)

type Tache struct {
	UUID      string         `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name        string `gorm:"not null" json:"name"`                          // Nom de la tâche
	Description string `json:"description"`                                   // Description de la tâche
	Statut      string `gorm:"not null" json:"statut"`                        // Statut de la tâche (ex: "En cours", "Terminé", "En attente")
	AgentUUID   string `gorm:"type:varchar(255);not null" json:"agent_uuid"`  // UUID de l'agent assigné à la tâche
	Agent       Agent  `gorm:"foreignKey:AgentUUID" json:"agent"`             // Relation avec l'agent
	TicketUUID  string `gorm:"type:varchar(255);not null" json:"ticket_uuid"` // UUID du ticket associé à la tâche
}
