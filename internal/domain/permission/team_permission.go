package permission

import (
	"fmt"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
)

var TeamPermissions = map[entity.TeamMemberRole]map[Permission]bool{
	entity.TeamMemberRoleOwner: {
		PermissionUpdateTeamDetails:        true,
		PermissionAddTeamPlayerImage:       true,
		PermissionAddPlayer:                true,
		PermissionAddStaff:                 true,
		PermissionUpdateTeamContactDetails: true,
		PermissionCreateInvite:             true,
	},
	entity.TeamMemberRoleManager: {
		PermissionAddPlayer:          true,
		PermissionAddTeamPlayerImage: true,
		PermissionAddStaff:           true,
		PermissionUpdateTeamDetails:  true,
		PermissionCreateInvite:       true,
	},
}

func HasPermissionTeam(role entity.TeamMemberRole, perm Permission) bool {
	fmt.Printf("Role: %q\n", role)
	member, ok := TeamPermissions[role]
	if !ok {
		return false
	}
	return member[perm]
}
