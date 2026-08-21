package teammemberuc

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)


type AddTeamOwnerRes struct {
	TeamMemberID uuid.UUID `gorm:"primaryKey"`
	TeamID       uuid.UUID
	UserID       string
	FullName     string
	Role         entity.TeamMemberRole
}

type RegisterTeamMemberReq struct {
	UserID string
	Code   string
}

type RegisterTeamMemberRes struct {
	InvitationID uuid.UUID
}

type CompensateRegisterTeamMemberReq struct {
	InvitationID uuid.UUID
	TeamMemberID uuid.UUID
}
