package player

import (
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type AddPlayerReq struct {
	TeamID       string
	UserID       string
	FullName     string
	DOB          time.Time
	JerseyNumber int32
	Postion      entity.PlayerPosition
	Height       float32
	Weight       float32
	ImageBytes   []byte
	ContentType  string
}

type AddPlayerRes struct {
	ID           uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	JerseyNumber int32
	Position     entity.PlayerPosition
	Status       entity.PlayerStatus
	PresignedUrl string
}

type UpdatPlayerStatusReq struct {
	PlayerID string
	Status   entity.PlayerStatus
}

type UpdatePlayerStatusRes struct {
	PlayerID string
	Status   entity.PlayerStatus
}

type ListTeamPlayersReq struct {
	TeamID       string
	Page         int32
	Limit        int32
	PlayerStatus entity.PlayerStatus
	Position     entity.PlayerPosition
	Search       string
}

type PlayerRes struct {
	ID           uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	JerseyNumber int32
	Position     entity.PlayerPosition
	Status       entity.PlayerStatus
}

type PaginateDetails struct {
	TotalPage int32
	Page      int32
	Limit     int32
	Total     int64
}

type GetPlayerReq struct {
	PlayerID string
}

type GetPlayerRes struct {
	PlayerID     uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	DateOfBirth  time.Time
	JerseyNumber int32
	Position     entity.PlayerPosition
	Height       float32
	Weight       float32
	Status       entity.PlayerStatus
	CreatedAt    time.Time
}

type ReleasePlayerReq struct {
	UserID   string
	TeamID   string
	PlayerID string
}

type ReleasePlayerRes struct {
	Success bool
}

type UpdatePlayerImageReq struct {
	UserID      string
	TeamID      string
	PlayerID    string
	ImageData   []byte
	ContentType string
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
