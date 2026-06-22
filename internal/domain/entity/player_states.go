package entity

import (
	"time"

	"github.com/google/uuid"
)

type PlayerStats struct {
	PlayerID      uuid.UUID `gorm:"primaryKey"`
	Appearances   int32
	MinutesPlayed int32
	Goals         int32
	Assists       int32
	YellowCards   int32
	RedCards      int32
	Saves         int32
	CleanSheets   int32
	GoalsConceded int32
	UpdatedAt     time.Time
}


