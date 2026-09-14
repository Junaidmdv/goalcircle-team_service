package entity

import (
	"time"

	"github.com/google/uuid"
)

type PlayerStatus string

const (
	PlayerStatusActive    PlayerStatus = "active"
	PlayerStatusInactive  PlayerStatus = "inactive"
	PlayerStatusSuspended PlayerStatus = "suspended"
	PlayerStatusLeft      PlayerStatus = "left"
	PlayerStatusInvalid   PlayerStatus = "invalid"
)

type PlayerPosition string

const (
	PlayerPositionGoalkeeper  PlayerPosition = "goalkeeper"
	PlayerPositionDefender    PlayerPosition = "defender"
	PlayerPositionMidfielder  PlayerPosition = "midfielder"
	PlayerPositionForward     PlayerPosition = "forward"
	PlayerPositionUnspecified PlayerPosition = "unspecified"
)

type PreferredFoot string

const (
	PreferredFootLeft  PreferredFoot = "left"
	PreferredFootRight PreferredFoot = "right"
	PreferredFootBoth  PreferredFoot = "both"
)

type Player struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	FullName        string
	DateOfBirth     time.Time      `gorm:"type:date"`
	Nationality     string         `gorm:"type:varchar(100)"`
	Position        PlayerPosition 
	PreferredFoot   PreferredFoot  
	ProfileImageKey *string        `gorm:"type:varchar(500)"`
	Status          PlayerStatus   
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
