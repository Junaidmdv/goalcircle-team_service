package entity

import (
	"time"

	"github.com/google/uuid"
)

type InviteType string

const (
	StaffInvite  InviteType = "stafff"
	PlayerInvite InviteType = "player"
)

type InvitationStatus string

const (
	Pending   InvitationStatus = "pending"
	Accepted  InvitationStatus = "accepted"
	Rejected  InvitationStatus = "rejected"
	Expired   InvitationStatus = "expired"
	cancelled InvitationStatus = "canelled"
)

type TeamInvite struct {
	ID              uuid.UUID `gorm:"primaryKey"`
	TeamID          uuid.UUID
	InvitedByUserID uuid.UUID
	InvitedUserID   uuid.UUID
	InviteType      InviteType
	Status          InvitationStatus `gorm:"default:pending"`
	InviteCode      string `gorm:"unique"`
	ExpiresAt       time.Time
	CreatedAt       time.Time 
	UpdatedAt       time.Time 
}
