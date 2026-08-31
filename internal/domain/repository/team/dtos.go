package team

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
)



type ListTeamsReq struct {
	Page       int32
	Limit      int32
	City       string
	TeamStatus entity.TeamStatus
	Search     string
}
