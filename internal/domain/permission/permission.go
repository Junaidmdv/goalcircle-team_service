package permission

type Permission string

const (
	PermissionUpdateTeamDetails        Permission = "UPDATE_TEAM_DETAILS"
	PermissionAddPlayer                Permission = "ADD_PLAYER"
	PermissionAddStaff                 Permission = "ADD_STAFF"
	PermissionAddManager               Permission = "ADD_MANAGER"
	PermissionUpdateTeamContactDetails Permission = "UPDATE_TEAM_CONTACT_DETAILS"
	PermissionApointCaptain            Permission = "APOINT_CAPTAIN"
	PermissionApointVicecaptain        Permission = "APOINT_VICECAPTAIN"
	PermissionAddTeamPlayerImage       Permission = "ADD_PLAYER_IMAGE"
	PermissionCreateInvite             Permission = "CREATE_INVITE"
)
