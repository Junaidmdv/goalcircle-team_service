package entity

import (
	"time"

	"github.com/google/uuid"
)



type PlayerInvitationStatus string

const (
	InvitationPending   PlayerInvitationStatus = "pending"
	InvitationAccepted  PlayerInvitationStatus = "accepted"
	InvitationRejected  PlayerInvitationStatus = "rejected"
	InvitationCancelled PlayerInvitationStatus = "canelled"
	InvitationExpired   PlayerInvitationStatus = "expired"
	InvitationInvalid   PlayerInvitationStatus = "invalid"
)

type PlayerTeamInvite struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	TeamID           uuid.UUID `gorm:"type:uuid;not null;index"`
	PlayerID         uuid.UUID `gorm:"type:uuid;not null;index"`
	OfferedJerseyNum *int32
	Message          *string 
	RejectReason     *string
	InvitedByUserID  uuid.UUID
	InvitedUserID    uuid.UUID
	Status           PlayerInvitationStatus `gorm:"default:pending"`
	ExpiresAt        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	RespondedAt      time.Time 

	Team Team `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE"` 
	Player Player `gorm:"foreignKey:PlayerID;constraint:OnDelete:CASCADE"`
}
