package team

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type UpdateTeamReq struct {
	TeamID      uuid.UUID
	Name        *string
	ShortName   *string
	City        *string
	Description *string
}

type UpdateTeamContact struct {
	ContactEmail    *string
	ContactPhoneNum *string
}

type ListTeamsReq struct {
	Page       int32
	Limit      int32
	City       string
	TeamStatus entity.TeamStatus
	Search     string
}
