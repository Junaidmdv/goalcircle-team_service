package entity

import (
	"time"

	"github.com/google/uuid"
)

type TeamPlayerStatus string

const (
	TeamPlayerStatusActive    TeamPlayerStatus = "active"
	TeamPlayerStatusInactive  TeamPlayerStatus = "inactive"
	TeamPlayerStatusSuspended TeamPlayerStatus = "suspended"
	TeamPlayerStatusLeft      TeamPlayerStatus = "left"
	TeamPlayerStatusInvalid   TeamPlayerStatus = "invalid"
)

type TeamPlayer struct {
	ID           uuid.UUID `json:"id"`
	PlayerID     uuid.UUID
	TeamMemberID uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	JerseyNumber int32
	JoinedAt     time.Time
	LeftAt       time.Time
	Status       TeamPlayerStatus `json:"status"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	TeamMember   TeamMember       `gorm:"foreignKey:TeamMemberID;references:ID"`
	Player       Player           `gorm:"foreignKey:PlayerID;references:ID"`
}
