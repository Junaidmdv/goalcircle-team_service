package staffuc

import (
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type AddStaffReq struct {
	TeamMemberID string
	TeamID       string
	FullName     string
	Designation  entity.StaffDesignation
	DOB          time.Time
	Role         entity.TeamMemberRole
}

type AddStaffRes struct {
	StaffID      uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	Age          int32
	Designation  entity.StaffDesignation
}

type ChangeStaffStatusRes struct {
}

type ChangeStaffStatusReq struct {
}

type ListTeamStaffReq struct {
}

type ListTeamStaffRes struct {
}
