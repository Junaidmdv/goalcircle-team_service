package teamsaga

import (
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type AddUserRoleRes struct {
	SessionID          string
	UserID             string
	Email              string
	AccessToken        string
	AccessTokenExpiry  time.Time
	RefreshToken       string
	RefreshTokenExpiry time.Time
}

type CreateTeamRes struct {
	ID        uuid.UUID
	Name      string
	ShortName string
	City      string
	Code      string
}

type TeamMemberRes struct {
	TeamMemberID uuid.UUID
	TeamID       uuid.UUID
	UserID       string
	FullName     string
	Role         entity.UserRole
}
