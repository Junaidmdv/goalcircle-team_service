package permission

import "github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"

var SquadPermission = map[entity.StaffDesignation]map[Permission]struct{}{
	entity.StaffDesignationOwner: {
		PermissionUpdateTeamDetails:  {},
		PermissionAddPlayer:          {},
		PermissionAddStaff:           {},
		PermissionAddManager:         {},
		PermissionAddPlayerImage:     {},
		PermissionCreateInvite:       {},
		PermissionAppointCaptain:     {},
		PermissionAppointViceCaptain: {},
	},

	entity.StaffDesignationManager: {
		PermissionUpdateTeamDetails:  {},
		PermissionAddPlayer:          {},
		PermissionAddStaff:           {},
		PermissionAddPlayerImage:     {},
		PermissionCreateInvite:       {},
		PermissionAppointCaptain:     {},
		PermissionAppointViceCaptain: {},
	},

	entity.StaffDesignationHeadCoach: {
		PermissionAppointCaptain:     {},
		PermissionAppointViceCaptain: {},
	},

	entity.StaffDesignationAssistantCoach: {},

}



func HasPermission(desig entity.StaffDesignation, perm Permission) bool {
	permite, ok := SquadPermission[desig]
	if !ok {
		return false
	}

	_, ok = permite[perm]
	return ok
}
