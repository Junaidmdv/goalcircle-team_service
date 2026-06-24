package entity

import (
	"time"

	"github.com/google/uuid"
)

type TeamInvite struct {
	ID        uuid.UUID
	TeamID    uuid.UUID
	TeamCode  string
	Code      string
	Role      UserRole
	IsUsed    bool
	UsedBy    *uuid.UUID
	ExpiresAt time.Time
	CreatedAt time.Time
}
