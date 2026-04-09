package models

import (
	"time"

	"gorm.io/gorm"
)

type Direction struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name        string   `gorm:"not null" json:"name"`
	Description string   `json:"description"`
	Bureau      []Bureau `gorm:"foreignKey:DirectionUUID" json:"bureaux"`
	Agents      []Agent  `gorm:"foreignKey:DirectionUUID" json:"agents"`
}
