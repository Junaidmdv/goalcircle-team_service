package team

import (
	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
)

var entityTeamStatusToProto = map[entity.TeamStatus]pb.TeamStatus{
	entity.TeamStatusActive:    pb.TeamStatus_TEAM_STATUS_ACTIVE,
	entity.TeamStatusInactive:  pb.TeamStatus_TEAM_STATUS_INACTIVE,
	entity.TeamStatusSuspended: pb.TeamStatus_TEAM_STATUS_SUSPENDED,
	entity.TeamStatusArchived:  pb.TeamStatus_TEAM_STATUS_ARCHIVED,
}

var protoTeamStatusToEntity = map[pb.TeamStatus]entity.TeamStatus{
	pb.TeamStatus_TEAM_STATUS_ACTIVE:    entity.TeamStatusActive,
	pb.TeamStatus_TEAM_STATUS_INACTIVE:  entity.TeamStatusInactive,
	pb.TeamStatus_TEAM_STATUS_SUSPENDED: entity.TeamStatusSuspended,
	pb.TeamStatus_TEAM_STATUS_ARCHIVED:  entity.TeamStatusArchived,
}

func MapTeamStatus(status pb.TeamStatus) entity.TeamStatus {
	if s, ok := protoTeamStatusToEntity[status]; ok {
		return s
	}
	return entity.TeamStatusActive // or another default that fits your business logic
}

func MapProtoTeamStatus(status entity.TeamStatus) pb.TeamStatus {
	if s, ok := entityTeamStatusToProto[status]; ok {
		return s
	}
	return pb.TeamStatus_TEAM_STATUS_UNSPECIFIED
}
