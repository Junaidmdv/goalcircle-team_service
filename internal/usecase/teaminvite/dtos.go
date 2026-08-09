package teaminvite

import (
	"time"

	"github.com/google/uuid"
)

type TeamInviteReq struct {
	UserID       string
	TeamMemberID string
}

type TeamInviteRes struct {
	TeamMemberId uuid.UUID
	Code         string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
