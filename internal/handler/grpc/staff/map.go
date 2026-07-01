package staff

import (
	teamv1 "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
)

var protoToDomainTeamMemberRole = map[teamv1.TeamMemberRole]entity.TeamMemberRole{
	teamv1.TeamMemberRole_TEAM_MEMBER_ROLE_OWNER:   entity.TeamMemberRoleOwner,
	teamv1.TeamMemberRole_TEAM_MEMBER_ROLE_MANAGER: entity.TeamMemberRoleManager,
	teamv1.TeamMemberRole_TEAM_MEMBER_ROLE_STAFF:   entity.TeamMemberRoleStaff,
	teamv1.TeamMemberRole_TEAM_MEMBER_ROLE_PLAYER:  entity.TeamMemberRolePlayer,
}

var domainToProtoTeamMemberRole = map[entity.TeamMemberRole]teamv1.TeamMemberRole{
	entity.TeamMemberRoleOwner:   teamv1.TeamMemberRole_TEAM_MEMBER_ROLE_OWNER,
	entity.TeamMemberRoleManager: teamv1.TeamMemberRole_TEAM_MEMBER_ROLE_MANAGER,
	entity.TeamMemberRoleStaff:   teamv1.TeamMemberRole_TEAM_MEMBER_ROLE_STAFF,
	entity.TeamMemberRolePlayer:  teamv1.TeamMemberRole_TEAM_MEMBER_ROLE_PLAYER,
}

func MapTeamMemberRole(role teamv1.TeamMemberRole) entity.TeamMemberRole {
	if r, ok := protoToDomainTeamMemberRole[role]; ok {
		return r
	}
	return entity.TeamMemberRole("")
}

func MapProtoTeamMemberRole(role entity.TeamMemberRole) teamv1.TeamMemberRole {
	if r, ok := domainToProtoTeamMemberRole[role]; ok {
		return r
	}
	return teamv1.TeamMemberRole_TEAM_MEMBER_ROLE_UNSPECIFIED
}

var protoToDomainStaffDesignation = map[teamv1.StaffDesignation]entity.StaffDesignation{
	teamv1.StaffDesignation_STAFF_DESIGNATION_HEAD_COACH:        entity.StaffDesignationHeadCoach,
	teamv1.StaffDesignation_STAFF_DESIGNATION_ASSISTANT_COACH:   entity.StaffDesignationAssistantCoach,
	teamv1.StaffDesignation_STAFF_DESIGNATION_GOALKEEPER_COACH:  entity.StaffDesignationGoalkeeperCoach,
	teamv1.StaffDesignation_STAFF_DESIGNATION_FITNESS_COACH:     entity.StaffDesignationFitnessCoach,
	teamv1.StaffDesignation_STAFF_DESIGNATION_PHYSIO:            entity.StaffDesignationPhysio,
	teamv1.StaffDesignation_STAFF_DESIGNATION_DOCTOR:            entity.StaffDesignationDoctor,
	teamv1.StaffDesignation_STAFF_DESIGNATION_ANALYST:           entity.StaffDesignationAnalyst,
	teamv1.StaffDesignation_STAFF_DESIGNATION_TEAM_COORDINATOR:  entity.StaffDesignationTeamCoordinator,
	teamv1.StaffDesignation_STAFF_DESIGNATION_EQUIPMENT_MANAGER: entity.StaffDesignationEquipmentManager,
	teamv1.StaffDesignation_STAFF_DESIGNATION_OTHER:              entity.StaffDesignationOther,
}

var domainToProtoStaffDesignation = map[entity.StaffDesignation]teamv1.StaffDesignation{
	entity.StaffDesignationHeadCoach:        teamv1.StaffDesignation_STAFF_DESIGNATION_HEAD_COACH,
	entity.StaffDesignationAssistantCoach:   teamv1.StaffDesignation_STAFF_DESIGNATION_ASSISTANT_COACH,
	entity.StaffDesignationGoalkeeperCoach:  teamv1.StaffDesignation_STAFF_DESIGNATION_GOALKEEPER_COACH,
	entity.StaffDesignationFitnessCoach:     teamv1.StaffDesignation_STAFF_DESIGNATION_FITNESS_COACH,
	entity.StaffDesignationPhysio:           teamv1.StaffDesignation_STAFF_DESIGNATION_PHYSIO,
	entity.StaffDesignationDoctor:           teamv1.StaffDesignation_STAFF_DESIGNATION_DOCTOR,
	entity.StaffDesignationAnalyst:          teamv1.StaffDesignation_STAFF_DESIGNATION_ANALYST,
	entity.StaffDesignationTeamCoordinator:  teamv1.StaffDesignation_STAFF_DESIGNATION_TEAM_COORDINATOR,
	entity.StaffDesignationEquipmentManager: teamv1.StaffDesignation_STAFF_DESIGNATION_EQUIPMENT_MANAGER,
	entity.StaffDesignationOther:            teamv1.StaffDesignation_STAFF_DESIGNATION_OTHER,
}

func MapStaffDesignation(designation teamv1.StaffDesignation) entity.StaffDesignation {
	if d, ok := protoToDomainStaffDesignation[designation]; ok {
		return d
	}
	return entity.StaffDesignationOther
}

func MapProtoStaffDesignation(designation entity.StaffDesignation) teamv1.StaffDesignation {
	if d, ok := domainToProtoStaffDesignation[designation]; ok {
		return d
	}
	return teamv1.StaffDesignation_STAFF_DESIGNATION_UNSPECIFIED
}
