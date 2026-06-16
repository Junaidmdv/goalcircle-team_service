package entity

import (
	"time"

	"github.com/google/uuid"
)

type Player struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	TeamMemberID uuid.UUID `gorm:"type:uuid;uniqueIndex"` 
	FullName     string
	DateOfBirth  time.Time
	JerseyNumber int
	Position     string
	Height       float64
	Weight       float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}




