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
	Weight        float32
	Height        float32
}

type AddPlayerRes struct {
	ID            uuid.UUID
	Code          string
	FullName      string
	DOB           time.Time
	Nationality   string
	PreferredFoot entity.PreferredFoot
	Position      entity.PlayerPosition
	PresignedUrl  string
	Weight        float32
	Height        float32
}

type UpdatePlayerReq struct {
	UserID        string
	FullName      *string
	DOB           *time.Time
	Nationality   *string
	PreferredFoot *entity.PreferredFoot
	Position      *entity.PlayerPosition
	Weight        *float32
	Height        *float32
}

type UpdatePlayersRes struct {
	PlayerID      uuid.UUID
	UserID        uuid.UUID
	Code          string
	FullName      string
	DOB           time.Time
	Nationality   string
	PreferredFoot entity.PreferredFoot
	Position      entity.PlayerPosition
	Weight        float32
	Height        float32
}

type PlayerRes struct {
	ID           uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	Position     entity.PlayerPosition
	ImageUrl     string
}

type GetPlayerReq struct {
	UserID   string
	PlayerID string
}

type GetPlayerRes struct {
	PlayerID     uuid.UUID
	PresignedUrl string
	Code         string
	FullName     string
	DateOfBirth  time.Time
	Position     entity.PlayerPosition
	PreferedFoot entity.PreferredFoot
	Nationality  string
	Height       float32
	Weight       float32
	CreatedAt    time.Time
}

type UpdatePlayerImageReq struct {
	UserID    string
	ImageData []byte
}

type UpdatePlayerImageRes struct {
    UserID uuid.UUID 
	PresignedUrl string
}

type GetPlayerPresignedUrlReq struct {
	PlayerID string
}

type GetPlayerPresignedUrlRes struct {
	PlayerId     string
	PresignedUrl string
}

type RemovePlayerImageReq struct {
	UserID string
}

type RemovePlayerImageRes struct {
	Success bool
}

type ListPlayerReq struct {
	Page     int32
	Limit    int32
	Position entity.PlayerPosition
	Search   *string
}

type ListPlayerDetails struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	FullName     string
	Position     entity.PlayerPosition
	PreferedFoot entity.PreferredFoot
	ImageUrl     string
	Code         string
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

type PlayerProfileReq struct {
	UserID string
}

type PlayerProfileRes struct {
	UserID          uuid.UUID
	PlayerID        uuid.UUID
	FullName        string
	ProfileImageUrl string
	DateOfBirth     time.Time
	Position        entity.PlayerPosition
	PreferedFoot    entity.PreferredFoot
	Nationality     string
	Height          float32
	Weight          float32
	CreatedAt       time.Time
	Code            string
}

type GetMyTeamsReq struct {
}

type GetMyTeamsRes struct {
}
