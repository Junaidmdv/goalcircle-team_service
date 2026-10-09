package entity

import (
	"time"

	"github.com/google/uuid"
)

type StaffInviteStatus string

type StaffInvite struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	TeamID      uuid.UUID
	InviteCode  string `gorm:"unique"`
	FullName    string
	Role        StaffRole
	Designation StaffDesignation
	PhoneNume   *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpiresAt   time.Time
	IsUsed      bool

	Team Team `gorm:"foreignKey:TeamID;constraint:OnDelete;CASCADE"`
}
