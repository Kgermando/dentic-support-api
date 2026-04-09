package models

import (
	"time"

	"gorm.io/gorm"
)

type Bureau struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name          string     `gorm:"not null" json:"name"`
	Description   string     `json:"description"`
	DirectionUUID string     `json:"direction_uuid"`
	Direction     Direction  `gorm:"foreignKey:DirectionUUID" json:"direction"`
	Agents        []Agent    `gorm:"foreignKey:BureauUUID" json:"agents"`
	TeamJoins     []TeamJoin `gorm:"foreignKey:BureauUUID" json:"team_joins"`
}
