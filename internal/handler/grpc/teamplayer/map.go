package teamplayerhandler

import (
	"errors"

	teampb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
)

func toPlayerStatus(status teampb.PlayerStatus) (entity.TeamPlayerStatus, error) {
	switch status {
	case teampb.PlayerStatus_PLAYER_STATUS_ACTIVE:
		return entity.TeamPlayerStatusActive, nil

	case teampb.PlayerStatus_PLAYER_STATUS_INACTIVE:
		return entity.TeamPlayerStatusInactive, nil

	case teampb.PlayerStatus_PLAYER_STATUS_SUSPENDED:
		return entity.TeamPlayerStatusSuspended, nil

	case teampb.PlayerStatus_PLAYER_STATUS_LEFT:
		return entity.TeamPlayerStatusLeft, nil

	default:
		return "", errors.New("invalid player status")
	}
}


func toPlayerPosition(position teampb.PlayerPosition) (entity.PlayerPosition, error) {
	switch position {
	case teampb.PlayerPosition_PLAYER_POSITION_GOALKEEPER:
		return entity.PlayerPositionGoalkeeper, nil

	case teampb.PlayerPosition_PLAYER_POSITION_DEFENDER:
		return entity.PlayerPositionDefender, nil

	case teampb.PlayerPosition_PLAYER_POSITION_MIDFIELDER:
		return entity.PlayerPositionMidfielder, nil

	case teampb.PlayerPosition_PLAYER_POSITION_FORWARD:
		return entity.PlayerPositionForward, nil

	default:
		return "", errors.New("invalid player position")
	}
}

func toPreferredFoot(foot teampb.PreferredFoot) (entity.PreferredFoot, error) {
	switch foot {
	case teampb.PreferredFoot_PREFERRED_FOOT_LEFT:
		return entity.PreferredFootLeft, nil

	case teampb.PreferredFoot_PREFERRED_FOOT_RIGHT:
		return entity.PreferredFootRight, nil

	case teampb.PreferredFoot_PREFERRED_FOOT_BOTH:
		return entity.PreferredFootBoth, nil

	default:
		return "", errors.New("invalid preferred foot")
	}
}