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

type PlayerPosition string

const (
	PositionGK PlayerPosition = "gk"

	PositionCB  PlayerPosition = "cb"
	PositionLB  PlayerPosition = "lb"
	PositionRB  PlayerPosition = "rb"
	PositionLWB PlayerPosition = "lwb"
	PositionRWB PlayerPosition = "rwb"

	PositionCDM PlayerPosition = "cdm"
	PositionCM  PlayerPosition = "cm"
	PositionCAM PlayerPosition = "cam"
	PositionLM  PlayerPosition = "lm"
	PositionRM  PlayerPosition = "rm"

	PositionLW PlayerPosition = "lw"
	PositionRW PlayerPosition = "rw"
	PositionST PlayerPosition = "st"
	PositionCF PlayerPosition = "cf"
)

type Player struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	TeamMemberID uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	FullName     string
	DateOfBirth  time.Time
	JerseyNumber int32
	Position     PlayerPosition
	Height       float32
	Weight       float32
	Status       PlayerStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
