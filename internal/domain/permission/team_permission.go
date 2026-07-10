package permission

import "github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"

var TeamPermissions = map[entity.TeamMemberRole]map[Permission]bool{
	entity.TeamMemberRoleOwner: {
		PermissionUpdateTeamDetails:        true,
		PermissionAddPlayer:                true,
		PermissionAddStaff:                 true,
		PermissionUpdateTeamContactDetails: true,
	},
	entity.TeamMemberRoleManager: {
		PermissionAddPlayer:         true,
		PermissionAddStaff:          true,
		PermissionUpdateTeamDetails: true,
	},
}

func HasPermissionTeam(role entity.TeamMemberRole, perm Permission) bool {
	member, ok := TeamPermissions[role]
	if !ok {
		return false
	}
	return member[perm]
}
