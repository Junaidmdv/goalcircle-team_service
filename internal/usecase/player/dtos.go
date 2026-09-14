package playeruc

import (
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type AddPlayerReq struct {
	UserID        string
	FullName      string
	DOB           time.Time
	Nationality   string
	PreferredFoot entity.PreferredFoot
	Position      entity.PlayerPosition
	ImageBytes    []byte
}

type AddPlayerRes struct {
	ID            uuid.UUID
	FullName      string
	DOB           time.Time
	Nationality   string
	PreferredFoot entity.PreferredFoot
	Position      entity.PlayerPosition
	Status        entity.PlayerStatus
	PresignedUrl  string
}

type UpdatePlayerReq struct {
	PlayerID      string
	UserID        string
	FullName      *string
	DOB           *time.Time
	Nationality   *string
	PreferredFoot *entity.PreferredFoot
	Position      *entity.PlayerPosition
}

type UpdatePlayersRes struct {
	PlayerID      uuid.UUID
	Status        entity.PlayerStatus
	FullName      string
	DOB           time.Time
	Nationality   string
	PreferredFoot entity.PreferredFoot
	Position      entity.PlayerPosition
}

type PlayerRes struct {
	ID           uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	Position     entity.PlayerPosition
	Status       entity.PlayerStatus
	ImageUrl     string
}

type GetPlayerReq struct {
	UserID   string
	PlayerID string
}

type GetPlayerRes struct {
	PlayerID     uuid.UUID
	FullName     string
	DateOfBirth  time.Time
	Position     entity.PlayerPosition
	PreferedFoot entity.PreferredFoot
	Nationality  string
	Status       entity.PlayerStatus
	CreatedAt    time.Time
}

type UpdatePlayerImageReq struct {
	TeamID    string
	UserID    string
	PlayerID  string
	ImageData []byte
}

type UpdatePlayerImageRes struct {
	TeamID       string
	PlayerID     string
	PresignedUrl string
}

type GetPlayerPresignedUrlReq struct {
	UserID   string
	PlayerID string
	TeamID   string
}

type GetPlayerPresignedUrlRes struct {
	TeamID       string
	PlayerId     string
	PresignedUrl string
}

type RemovePlayerImageReq struct {
	UserID   string
	TeamID   string
	PlayerID string
}

type RemovePlayerImageRes struct {
	Success bool
}

type ListPlayerReq struct {
	Page         int32
	Limit        int32
	PlayerStatus entity.PlayerStatus
	Position     entity.PlayerPosition
	Search       string
}

type ListPlayerDetails struct {
	ID           uuid.UUID
	FullName     string
	Position     entity.PlayerPosition
	PreferedFoot entity.PreferredFoot
	Status       entity.PlayerStatus
	ImageUrl     string
}

type PaginateDetails struct {
	TotalPage int32
	Page      int32
	Limit     int32
	Total     int64
}

type ListPlayerRes struct {
	Players    []ListPlayerDetails
	Pagination PaginateDetails
}
