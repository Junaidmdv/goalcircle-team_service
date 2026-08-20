package permission

// import (
// 	"fmt"

// 	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
// )

// var TeamPermissions = map[entity.TeamMemberRole]map[Permission]bool{
// 	entity.TeamMemberRoleOwner: {
// 	
// 	},
// 	entity.TeamMemberRoleManager: {
// 		
// 	},
// }

// func HasPermissionTeam(role entity.TeamMemberRole, perm Permission) bool {
// 	fmt.Printf("Role: %q\n", role)
// 	member, ok := TeamPermissions[role]
// 	if !ok {
// 		return false
// 	}
// 	return member[perm]
// }
