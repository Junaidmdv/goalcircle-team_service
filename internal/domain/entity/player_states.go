package entity

import (
	"time"

	"github.com/google/uuid"
)

type PlayerStats struct {
	ID            uuid.UUID `gorm:"primaryKey"`
	PlayerID      uuid.UUID
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


