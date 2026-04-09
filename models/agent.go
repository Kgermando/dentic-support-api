package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Agent struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Fullname        string     `gorm:"not null" json:"fullname"`
	Email           string     `gorm:"unique; not null" json:"email"`
	Telephone       string     `gorm:"unique" json:"telephone"`
	TranchAge       string     `gorm:"default:'25-34'" json:"tranch_age"` // permettra de savoir la moyenne d'age des agents
	Password        string     `json:"password"`
	PasswordConfirm string     `json:"password_confirm" gorm:"-"`
	Role            string     `json:"role"` // 'Directeur', 'Secretaire', 'Chef du bureau', 'Agent', 'SuperAdmin'
	Permission      string     `json:"permission"`
	Status          bool       `gorm:"default:false" json:"status"`
	DirectionUUID   string     `json:"direction_uuid"`
	Direction       Direction  `gorm:"foreignKey:DirectionUUID" json:"direction"`
	BureauUUID      string     `json:"bureau_uuid"`
	Bureau          Bureau     `gorm:"foreignKey:BureauUUID" json:"bureau"`
	TeamJoins       []TeamJoin `gorm:"foreignKey:AgentUUID" json:"team_joins"`
}

type AgentResponse struct {
	UUID            string     `json:"uuid"`
	Fullname        string     `json:"fullname"`
	Email           string     `json:"email"`
	Telephone       string     `json:"telephone"`
	TranchAge       string     `json:"tranch_age"` // permettra de savoir la moyenne d'age des agents
	Password        string     `json:"password"`
	PasswordConfirm string     `json:"password_confirm"`
	Role            string     `json:"role"` // 'Directeur', 'Secretaire', 'Chef du bureau', 'Agent', 'SuperAdmin'
	Permission      string     `json:"permission"`
	Status          bool       `json:"status"`
	DirectionUUID   string     `json:"direction_uuid"`
	Direction       Direction  `json:"direction"`
	BureauUUID      string     `json:"bureau_uuid"`
	Bureau          Bureau     `json:"bureau"`
	TeamJoins       []TeamJoin `json:"team_joins"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Login struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

func (a *Agent) SetPassword(p string) {
	hp, _ := bcrypt.GenerateFromPassword([]byte(p), 14)
	a.Password = string(hp)
}

func (a *Agent) ComparePassword(p string) error {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(p))
	return err
}
