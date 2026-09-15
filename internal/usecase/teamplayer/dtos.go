package teamplayeruc

import (
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type ListTeamPlayersReq struct {
	TeamID       string
	Page         int32
	Limit        int32
	PlayerStatus entity.PlayerStatus
	Position     entity.PlayerPosition
	Search       string
}

type PaginateDetails struct {
	TotalPage int32
	Page      int32
	Limit     int32
	Total     int64
}

type PlayerRes struct {
	ID           uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	//JerseyNumber int32
	Position entity.PlayerPosition
	Status   entity.PlayerStatus
	ImageUrl string
}

type GetTeamPlayerReq struct {
	TeamID   string
	PlayerID string
}

type GetTeamPlayerRes struct {
	TeamID       string
	PlayerID     uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	DateOfBirth  time.Time
	//JerseyNumber int32
	Position     entity.PlayerPosition
	PreferedFoot entity.PreferredFoot
	Nationality  string
	Status       entity.PlayerStatus
	CreatedAt    time.Time
}

type ReleasePlayerReq struct {
	UserID   string
	TeamID   string
	PlayerID string
}

type ReleasePlayerRes struct {
	ID           uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	JerseyNumber int32
	Position     entity.PlayerPosition
	Status       entity.PlayerStatus
	ReleasedAt   time.Time
	JoinedAt     time.Time
}
