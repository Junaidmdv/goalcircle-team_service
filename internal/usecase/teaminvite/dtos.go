package teaminvite

import (
	"time"

	"github.com/google/uuid"
)

type TeamInviteReq struct {
	TeamMemberID string
	Role         string
}

type TeamInviteRes struct {
	TeamMemberId uuid.UUID
	Code         string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
