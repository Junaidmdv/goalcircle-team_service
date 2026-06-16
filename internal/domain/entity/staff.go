package entity

import (
	"time"
	"github.com/google/uuid"
)

type Staff struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	TeamMemberID uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	FullName    string
	Designation string
	CreatedAt time.Time
	UpdatedAt time.Time
}
