package teammemberuc

import (
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
) 



type AddOwnerReq struct {
	UserID           string
	TeamID           string
	FullName         string
	Role             entity.StaffRole
	Designation      entity.StaffDesignation
	DOB              time.Time
	TeamMemberStatus entity.TeamMemberStatus
}

type AddOwnerRes struct {
	StaffID      uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	Age          int32
	Designation  entity.StaffDesignation
	Role         entity.StaffRole
}



type RegisterTeamMemberReq struct {
	UserID string
	Code   string
}

type RegisterTeamMemberRes struct {
	InvitationID uuid.UUID
	TeamMemberID uuid.UUID
	TeamID       uuid.UUID
	UserID       uuid.UUID
	FullName     string
	Role         string
	Status       string
	JoinedAt     time.Time
}

type CompensateRegisterTeamMemberReq struct {
	InvitationID uuid.UUID
	TeamMemberID uuid.UUID
	UserID       uuid.UUID
}
