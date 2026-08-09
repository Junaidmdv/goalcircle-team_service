package player

import (
	teamv1 "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
)

var protoToDomainPlayerPosition = map[teamv1.PlayerPosition]entity.PlayerPosition{
	teamv1.PlayerPosition_PLAYER_POSITION_UNSPECIFIED: entity.PositionUnspecified,

	teamv1.PlayerPosition_PLAYER_POSITION_GK:  entity.PositionGK,
	teamv1.PlayerPosition_PLAYER_POSITION_CB:  entity.PositionCB,
	teamv1.PlayerPosition_PLAYER_POSITION_LB:  entity.PositionLB,
	teamv1.PlayerPosition_PLAYER_POSITION_RB:  entity.PositionRB,
	teamv1.PlayerPosition_PLAYER_POSITION_LWB: entity.PositionLWB,
	teamv1.PlayerPosition_PLAYER_POSITION_RWB: entity.PositionRWB,

	teamv1.PlayerPosition_PLAYER_POSITION_CDM: entity.PositionCDM,
	teamv1.PlayerPosition_PLAYER_POSITION_CM:  entity.PositionCM,
	teamv1.PlayerPosition_PLAYER_POSITION_CAM: entity.PositionCAM,
	teamv1.PlayerPosition_PLAYER_POSITION_RM:  entity.PositionRM,
	teamv1.PlayerPosition_PLAYER_POSITION_LM:  entity.PositionLM,

	teamv1.PlayerPosition_PLAYER_POSITION_CF: entity.PositionCF,
	teamv1.PlayerPosition_PLAYER_POSITION_LW: entity.PositionLW,
	teamv1.PlayerPosition_PLAYER_POSITION_RW: entity.PositionRW,
	teamv1.PlayerPosition_PLAYER_POSITION_ST: entity.PositionST,
}

var domainToProtoPlayerPosition = map[entity.PlayerPosition]teamv1.PlayerPosition{
	entity.PositionUnspecified: teamv1.PlayerPosition_PLAYER_POSITION_UNSPECIFIED,

	entity.PositionGK:  teamv1.PlayerPosition_PLAYER_POSITION_GK,
	entity.PositionCB:  teamv1.PlayerPosition_PLAYER_POSITION_CB,
	entity.PositionLB:  teamv1.PlayerPosition_PLAYER_POSITION_LB,
	entity.PositionRB:  teamv1.PlayerPosition_PLAYER_POSITION_RB,
	entity.PositionLWB: teamv1.PlayerPosition_PLAYER_POSITION_LWB,
	entity.PositionRWB: teamv1.PlayerPosition_PLAYER_POSITION_RWB,

	entity.PositionCDM: teamv1.PlayerPosition_PLAYER_POSITION_CDM,
	entity.PositionCM:  teamv1.PlayerPosition_PLAYER_POSITION_CM,
	entity.PositionCAM: teamv1.PlayerPosition_PLAYER_POSITION_CAM,
	entity.PositionRM:  teamv1.PlayerPosition_PLAYER_POSITION_RM,
	entity.PositionLM:  teamv1.PlayerPosition_PLAYER_POSITION_LM,

	entity.PositionCF: teamv1.PlayerPosition_PLAYER_POSITION_CF,
	entity.PositionLW: teamv1.PlayerPosition_PLAYER_POSITION_LW,
	entity.PositionRW: teamv1.PlayerPosition_PLAYER_POSITION_RW,
	entity.PositionST: teamv1.PlayerPosition_PLAYER_POSITION_ST,
}

func MapPlayerPosition(position teamv1.PlayerPosition) entity.PlayerPosition {
	if p, ok := protoToDomainPlayerPosition[position]; ok {
		return p
	}
	return entity.PositionUnspecified
}

func MapProtoPlayerPosition(position entity.PlayerPosition) teamv1.PlayerPosition {
	if p, ok := domainToProtoPlayerPosition[position]; ok {
		return p
	}
	return teamv1.PlayerPosition_PLAYER_POSITION_UNSPECIFIED
}

var protoToDomainPlayerStatus = map[teamv1.PlayerStatus]entity.PlayerStatus{
	teamv1.PlayerStatus_STATUS_ACTIVE:      entity.PlayerStatusActive,
	teamv1.PlayerStatus_STATUS_INACTIVE:    entity.PlayerStatusInjured,
	teamv1.PlayerStatus_STATUS_SUSPENDED:   entity.PlayerStatusSuspended,
	teamv1.PlayerStatus_STATUS_RELEASED:    entity.PlayerStatusReleased,
}

var domainToProtoPlayerStatus = map[entity.PlayerStatus]teamv1.PlayerStatus{
	entity.PlayerStatusActive:              teamv1.PlayerStatus_STATUS_ACTIVE,
	entity.PlayerStatusInjured:             teamv1.PlayerStatus_STATUS_INACTIVE,
	entity.PlayerStatusSuspended:           teamv1.PlayerStatus_STATUS_SUSPENDED,
	entity.PlayerStatusReleased:            teamv1.PlayerStatus_STATUS_RELEASED,
}

func MapPlayerStatus(status teamv1.PlayerStatus) entity.PlayerStatus {
	if s, ok := protoToDomainPlayerStatus[status]; ok {
		return s
	}
	return entity.PlayerStatusInvalid
}

func MapProtoPlayerStatus(status entity.PlayerStatus) teamv1.PlayerStatus {
	if s, ok := domainToProtoPlayerStatus[status]; ok {
		return s
	}
	return teamv1.PlayerStatus_STATUS_UNSPECIFIED
}
