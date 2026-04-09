package models

import (
	"time"

	"gorm.io/gorm"
)

type BureauDemandeur struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name                   string             `gorm:"not null" json:"name"`
	DirectionDemandeurUUID string             `gorm:"type:varchar(255)" json:"direction_demandeur_uuid"`
	DirectionDemandeur     DirectionDemandeur `gorm:"foreignKey:DirectionDemandeurUUID" json:"direction_demandeur"`
	Demandeurs             []Demandeur        `gorm:"foreignKey:BureauDemandeurUUID" json:"demandeurs"`
}
