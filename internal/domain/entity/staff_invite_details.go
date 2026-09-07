package entity

import (
	"time"

	"github.com/google/uuid"
)

type StaffInviteDetails struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	TeamInviteID uuid.UUID
	FullName     string
	Role         StaffRole
	Designation  StaffDesignation
	PhoneNume    *string
	ImageKey     *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
