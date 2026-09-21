package entity

import (
	"time"

	"github.com/google/uuid"
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
	PreferredFootLeft   PreferredFoot = "left"
	PreferredFootRight  PreferredFoot = "right"
	PreferredFootBoth   PreferredFoot = "both"
	PreferredFootInvalid PreferredFoot = "invalid"
)

type Player struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	Code            string    `gorm:"not null;uniqueIndex"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	FullName        string
	DateOfBirth     time.Time `gorm:"type:date"`
	Nationality     string    `gorm:"type:varchar(100)"`
	Position        PlayerPosition
	PreferredFoot   PreferredFoot
	ProfileImageKey *string `gorm:"type:varchar(500)"`
	Height          float32 `gorm:"type:numeric(5,2)"`
	Weight          float32 `gorm:"type:numeric(5,2)"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
