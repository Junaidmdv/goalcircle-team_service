package teaminvite

import (
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type TeamInviteReq struct {
	TeamMemberID string
	Role         entity.UserRole
}

type TeamInviteRes struct {
	TeamMemberId uuid.UUID
	code         string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
