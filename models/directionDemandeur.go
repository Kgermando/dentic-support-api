package models

import (
	"time"

	"gorm.io/gorm"
)

type DirectionDemandeur struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name string `gorm:"not null" json:"name"`

	BureauDemandeurs []BureauDemandeur `gorm:"foreignKey:DirectionDemandeurUUID" json:"bureau_demandeurs"`
	Demandeurs       []Demandeur       `gorm:"foreignKey:DirectionDemandeurUUID" json:"demandeurs"`
}
