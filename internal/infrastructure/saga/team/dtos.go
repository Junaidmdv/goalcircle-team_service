package teamsaga

import (
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type AddUserRoleReq struct {
	UserID       string
	Role         string
	RefreshToken string
}

type AddUserRoleRes struct {
	SessionID          string
	UserID             string
	Email              string
	AccessToken        string
	AccessTokenExpiry  time.Time
	RefreshToken       string
	RefreshTokenExpiry time.Time
}

type TeamDetails struct {
	ID        uuid.UUID
	Name      string
	ShortName string
	City      string
	Code      string
}

type OwnerDetails struct {
	StaffID      uuid.UUID
	TeamMemberID uuid.UUID
	TeamID       uuid.UUID
	UserID       string
	FullName     string
	Age          int32
	Role         entity.StaffRole
	Designation  entity.StaffDesignation
}

type CreateTeamRes struct {
	Team  *TeamDetails
	Owner *OwnerDetails
}

type TeamMemberRes struct {
}
