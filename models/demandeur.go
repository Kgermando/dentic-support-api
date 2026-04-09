package models

import (
	"time"

	"gorm.io/gorm"
)

type Demandeur struct {
	UUID      string         `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Fullname  string `gorm:"not null" json:"fullname"`
	Email     string `gorm:"unique;not null" json:"email"`
	Telephone string `gorm:"type:varchar(255)" json:"telephone"`

	Site                   string             `gorm:"type:varchar(255)" json:"site"` // Site de travail du demandeur
	DirectionDemandeurUUID string             `gorm:"type:varchar(255)" json:"direction_demandeur_uuid"`
	DirectionDemandeur     DirectionDemandeur `gorm:"foreignKey:DirectionDemandeurUUID" json:"direction_demandeur"`
	BureauDemandeurUUID    string             `gorm:"type:varchar(255)" json:"bureau_demandeur_uuid"`
	BureauDemandeur        BureauDemandeur    `gorm:"foreignKey:BureauDemandeurUUID" json:"bureau_demandeur"`
}
