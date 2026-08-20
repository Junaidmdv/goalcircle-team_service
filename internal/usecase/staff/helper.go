package staffuc

import (
	"fmt"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
)

var staffLimits = map[entity.StaffDesignation]int32{
	entity.StaffDesignationOwner:          1,
	entity.StaffDesignationManager:        1,
	entity.StaffDesignationHeadCoach:      1,
	entity.StaffDesignationAssistantCoach: 5,
	entity.StaffDesignationOther:          20,
}

func StaffMaxLimit(desig entity.StaffDesignation, count int32) error {
	limit, ok := staffLimits[desig]
	if !ok {
		return apperror.NewBadRequestError("invalid staff role")
	}

	if count >= limit {
		return apperror.NewBadRequestError(fmt.Sprintf("maximum limit of %s is %d", desig, limit))
	}
	return nil

}  



var validDesignations = map[entity.StaffRole]map[entity.StaffDesignation]struct{}{
	entity.StaffRoleManagement: {
		entity.StaffDesignationOwner:   {},
		entity.StaffDesignationManager: {},
	},

	entity.StaffRoleCoach: {
		entity.StaffDesignationHeadCoach:       {},
		entity.StaffDesignationAssistantCoach:  {},
		entity.StaffDesignationGoalkeeperCoach: {},
		entity.StaffDesignationFitnessCoach:    {},
	},

	entity.StaffRoleMedical: {
		entity.StaffDesignationPhysio: {},
		entity.StaffDesignationDoctor: {},
	},

	entity.StaffRoleOther: {
		entity.StaffDesignationAnalyst: {},
		entity.StaffDesignationOther:   {},
	},
}



func IsValidStaffRoleDesignation(
	role entity.StaffRole,
	designation entity.StaffDesignation,
) bool {
	designations, ok := validDesignations[role]
	if !ok {
		return false
	}

	_, ok = designations[designation]
	return ok
}

