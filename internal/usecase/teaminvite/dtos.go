package teaminvite

import (
	"time"

	"github.com/google/uuid"
)

type TeamInviteReq struct { 
	TeamID       string 
	UserID       string
	TeamMemberID string
}

type TeamInviteRes struct {
	TeamMemberId uuid.UUID
	Code         string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
