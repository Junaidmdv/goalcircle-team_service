package player

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
)

func GetTeamMemberStatusFromPlayerStatus(
	status entity.PlayerStatus,
) ( entity.TeamMemberStatus,error) {

	switch status {
	case entity.PlayerStatusActive:
		return  entity.TeamMemberStatusActive,nil

	case entity.PlayerStatusInjured,
		entity.PlayerStatusSuspended:
		return entity.TeamMemberStatusInactive,nil

	default:
		return "",apperror.NewInvalidArgumentError("invalid player status")
	}
}
