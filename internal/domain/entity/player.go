package entity

import (
	"time"

	"github.com/google/uuid"
)

type PlayerStatus string

const (
	PlayerStatusActive              PlayerStatus = "ACTIVE"
	PlayerStatusInjured             PlayerStatus = "INJURED"
	PlayerStatusSuspended           PlayerStatus = "SUSPENDED"
	PlayerStatusArchived            PlayerStatus = "ARCHIVED"
	PlayerStatusPendingActionvation PlayerStatus = "PENDING_ACTIVATION"
)

type Player struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	TeamMemberID uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	FullName     string
	DateOfBirth  time.Time
	JerseyNumber int32
	Position     string
	Height       float32
	Weight       float32
	Status       PlayerStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
