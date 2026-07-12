package entity

import (
	"time"

	"github.com/google/uuid"
)

type TeamInvite struct {
	ID             uuid.UUID
	TeamMemberID   uuid.UUID
	TeamCode       string
	Code           string
	TeamMemberType TeamMemberRole
	Role           string
	IsUsed         bool
	UserID         string
	ExpiresAt      time.Time
	CreatedAt      time.Time
	
}
