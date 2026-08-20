package staff

import (
	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	teamv1 "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
)

var protoToDomainStaffDesignation = map[teamv1.StaffDesignation]entity.StaffDesignation{
	teamv1.StaffDesignation_STAFF_DESIGNATION_UNSPECIFIED:      entity.StaffDesignationUnspecified,
	teamv1.StaffDesignation_STAFF_DESIGNATION_OWNER:            entity.StaffDesignationOwner,
	teamv1.StaffDesignation_STAFF_DESIGNATION_MANAGER:          entity.StaffDesignationManager,
	teamv1.StaffDesignation_STAFF_DESIGNATION_HEAD_COACH:       entity.StaffDesignationHeadCoach,
	teamv1.StaffDesignation_STAFF_DESIGNATION_ASSISTANT_COACH:  entity.StaffDesignationAssistantCoach,
	teamv1.StaffDesignation_STAFF_DESIGNATION_GOALKEEPER_COACH: entity.StaffDesignationGoalkeeperCoach,
	teamv1.StaffDesignation_STAFF_DESIGNATION_FITNESS_COACH:    entity.StaffDesignationFitnessCoach,
	teamv1.StaffDesignation_STAFF_DESIGNATION_PHYSIO:           entity.StaffDesignationPhysio,
	teamv1.StaffDesignation_STAFF_DESIGNATION_DOCTOR:           entity.StaffDesignationDoctor,
	teamv1.StaffDesignation_STAFF_DESIGNATION_ANALYST:          entity.StaffDesignationAnalyst,
	teamv1.StaffDesignation_STAFF_DESIGNATION_OTHER:            entity.StaffDesignationOther,
}

var domainToProtoStaffDesignation = map[entity.StaffDesignation]teamv1.StaffDesignation{
	entity.StaffDesignationUnspecified:     teamv1.StaffDesignation_STAFF_DESIGNATION_UNSPECIFIED,
	entity.StaffDesignationOwner:           teamv1.StaffDesignation_STAFF_DESIGNATION_OWNER,
	entity.StaffDesignationManager:         teamv1.StaffDesignation_STAFF_DESIGNATION_MANAGER,
	entity.StaffDesignationHeadCoach:       teamv1.StaffDesignation_STAFF_DESIGNATION_HEAD_COACH,
	entity.StaffDesignationAssistantCoach:  teamv1.StaffDesignation_STAFF_DESIGNATION_ASSISTANT_COACH,
	entity.StaffDesignationGoalkeeperCoach: teamv1.StaffDesignation_STAFF_DESIGNATION_GOALKEEPER_COACH,
	entity.StaffDesignationFitnessCoach:    teamv1.StaffDesignation_STAFF_DESIGNATION_FITNESS_COACH,
	entity.StaffDesignationPhysio:          teamv1.StaffDesignation_STAFF_DESIGNATION_PHYSIO,
	entity.StaffDesignationDoctor:          teamv1.StaffDesignation_STAFF_DESIGNATION_DOCTOR,
	entity.StaffDesignationAnalyst:         teamv1.StaffDesignation_STAFF_DESIGNATION_ANALYST,
	entity.StaffDesignationOther:           teamv1.StaffDesignation_STAFF_DESIGNATION_OTHER,
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

var protoToDomainStaffRole = map[teamv1.StaffRole]entity.StaffRole{
	teamv1.StaffRole_STAFF_ROLE_UNSPECIFIED: entity.StaffRoleUnspecified,
	teamv1.StaffRole_STAFF_ROLE_MANAGEMENT:  entity.StaffRoleManagement,
	teamv1.StaffRole_STAFF_ROLE_COACH:       entity.StaffRoleCoach,
	teamv1.StaffRole_STAFF_ROLE_MEDICAL:     entity.StaffRoleMedical,
	teamv1.StaffRole_STAFF_ROLE_OTHER:       entity.StaffRoleOther,
}

var domainToProtoStaffRole = map[entity.StaffRole]teamv1.StaffRole{
	entity.StaffRoleUnspecified: teamv1.StaffRole_STAFF_ROLE_UNSPECIFIED,
	entity.StaffRoleManagement:  teamv1.StaffRole_STAFF_ROLE_MANAGEMENT,
	entity.StaffRoleCoach:       teamv1.StaffRole_STAFF_ROLE_COACH,
	entity.StaffRoleMedical:     teamv1.StaffRole_STAFF_ROLE_MEDICAL,
	entity.StaffRoleOther:       teamv1.StaffRole_STAFF_ROLE_OTHER,
}

func MapStaffRole(role teamv1.StaffRole) entity.StaffRole {
	if r, ok := protoToDomainStaffRole[role]; ok {
		return r
	}

	return entity.StaffRoleOther
}

func MapProtoStaffRole(role entity.StaffRole) teamv1.StaffRole {
	if r, ok := domainToProtoStaffRole[role]; ok {
		return r
	}

	return teamv1.StaffRole_STAFF_ROLE_UNSPECIFIED
}

func ToCreateStaff(req *pb.StaffDetails) *AddStaffReq {
	role := MapStaffRole(req.Role)
	desig := MapStaffDesignation(req.Designation)

	return &AddStaffReq{
		TeamID:      req.TeamId,
		UserID:      req.UserId,
		FullName:    req.FullName,
		Role:        role,
		Designation: desig,
		DOB:         req.Dob.AsTime(),
	}
}

func ToUpdateStaff(req *pb.UpdateStaffReq) *UpdateStaffReq {
	dob := req.Dob.AsTime()
	role := MapStaffRole(*req.Role)
	desig := MapStaffDesignation(*req.Designation)
	return &UpdateStaffReq{
		UserID:      req.UserId,
		TeamID:      req.TeamId,
		StaffID:     req.StaffId,
		FullName:    req.FullName,
		DateOfBirth: &dob,
		Role:        &role,
		Designation: &desig,
		PhoneNum:    req.PhoneNum,
	}
}
