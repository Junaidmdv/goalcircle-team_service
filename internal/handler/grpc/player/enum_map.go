package player

import (
	"errors"

	teampb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
)

func toPlayerStatus(status *teampb.PlayerStatus) (entity.PlayerStatus, error) {
	if status == nil {
		return "", nil
	}

	switch status {
	case teampb.PlayerStatus_PLAYER_STATUS_ACTIVE.Enum():
		return entity.PlayerStatusActive, nil

	case teampb.PlayerStatus_PLAYER_STATUS_INACTIVE.Enum():
		return entity.PlayerStatusInactive, nil

	case teampb.PlayerStatus_PLAYER_STATUS_SUSPENDED.Enum():
		return entity.PlayerStatusSuspended, nil

	case teampb.PlayerStatus_PLAYER_STATUS_LEFT.Enum():
		return entity.PlayerStatusLeft, nil

	default:
		return "", errors.New("invalid player status")
	}
}

func toPlayerPosition(position *teampb.PlayerPosition) (entity.PlayerPosition, error) {
	if position == nil {
		return "", nil
	}
	switch position {
	case teampb.PlayerPosition_PLAYER_POSITION_GOALKEEPER.Enum():
		return entity.PlayerPositionGoalkeeper, nil

	case teampb.PlayerPosition_PLAYER_POSITION_DEFENDER.Enum():
		return entity.PlayerPositionDefender, nil

	case teampb.PlayerPosition_PLAYER_POSITION_MIDFIELDER.Enum():
		return entity.PlayerPositionMidfielder, nil

	case teampb.PlayerPosition_PLAYER_POSITION_FORWARD.Enum():
		return entity.PlayerPositionForward, nil

	default:
		return "", errors.New("invalid player position")
	}
}

func toPreferredFoot(foot *teampb.PreferredFoot) (entity.PreferredFoot, error) {

	if foot == nil {
		return "", nil
	}

	switch foot {
	case teampb.PreferredFoot_PREFERRED_FOOT_LEFT.Enum():
		return entity.PreferredFootLeft, nil

	case teampb.PreferredFoot_PREFERRED_FOOT_RIGHT.Enum():
		return entity.PreferredFootRight, nil

	case teampb.PreferredFoot_PREFERRED_FOOT_BOTH.Enum():
		return entity.PreferredFootBoth, nil

	default:
		return "", errors.New("invalid preferred foot")
	}
}
