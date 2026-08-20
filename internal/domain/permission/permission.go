package permission

type Permission string

const (
	PermissionUpdateTeamDetails Permission = "UPDATE_TEAM_DETAILS"
	PermissionAddPlayer  Permission = "ADD_PLAYER"
	PermissionAddStaff   Permission = "ADD_STAFF"
	PermissionAddManager Permission = "ADD_MANAGER"

	PermissionAppointCaptain     Permission = "APPOINT_CAPTAIN"
	PermissionAppointViceCaptain Permission = "APPOINT_VICE_CAPTAIN"

	PermissionAddPlayerImage Permission = "ADD_PLAYER_IMAGE"

	PermissionCreateInvite Permission = "CREATE_INVITE"  

	PermissionUpdatePlayerStatus Permission="UPDATE_PLAYER_STAT"
)
