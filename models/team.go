package models

import (
	"time"

	"gorm.io/gorm"
)

// Team représente une équipe de travail au sein des differents bureaux dans une meme direction
type Team struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`

	TeamJoins []TeamJoin `gorm:"foreignKey:TeamUUID" json:"team_joins"`
}

// Permet de faire le lien entre les équipes, les agents et les bureaux
type TeamJoin struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	TeamUUID   string `gorm:"type:varchar(255);not null" json:"team_uuid"`
	Team       Team   `gorm:"foreignKey:TeamUUID" json:"team"`
	AgentUUID  string `gorm:"type:varchar(255);not null" json:"agent_uuid"`
	Agent      Agent  `gorm:"foreignKey:AgentUUID" json:"agent"`
	BureauUUID string `gorm:"type:varchar(255);not null" json:"bureau_uuid"`
	Bureau     Bureau `gorm:"foreignKey:BureauUUID" json:"bureau"`
}
