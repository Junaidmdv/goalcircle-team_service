package entity

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	Name          string
	City          string
	Logo          string
	TeamStatus    string
	PlayerCount   string
	CaptionID     uuid.UUID
	ViceCaptionID uuid.UUID
	CreatedAt     time.Time
}
