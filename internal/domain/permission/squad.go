package permission

import "github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"

var SquadPermission = map[entity.StaffDesignation]map[Permission]bool{
	entity.StaffDesignationHeadCoach: {
		PermissionApointCaptain:     true,
		PermissionApointVicecaptain: true,
	},
}

func HasPermissionSquad(desig entity.StaffDesignation, perm Permission) bool {
	staff, ok := SquadPermission[desig]
	if !ok {
		return false
	}
	return staff[perm]
}
