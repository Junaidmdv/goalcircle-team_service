package entity

import (
	"time"

	"github.com/google/uuid"
)

type TeamStatus string

const (
	TeamStatusActive    TeamStatus = "ACTIVE"
	TeamStatusInactive  TeamStatus = "INACTIVE"
	TeamStatusSuspended TeamStatus = "SUSPENDED"
	TeamStatusArchived  TeamStatus = "ARCHIVED"
)

type Team struct {
	ID              uuid.UUID `gorm:"primaryKey"`
	Name            string
	ShortName       string
	City            string
	LogoUrl         string
	TeamCode        string
	Description     string
	ContactEmail    string
	ContactPhoneNum string
	TeamStatus      TeamStatus
	PlayerCount     int32
	CaptainID       uuid.UUID
	ViceCaptainID   uuid.UUID
	CreatedAt       time.Time
	TeamMember      []TeamMember `gorm:"foreignKey:TeamID"`
	TeamStats       TeamStats    `gorm:"foreignKey:TeamID"`
}
