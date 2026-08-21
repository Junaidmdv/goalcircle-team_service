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

type UpdatePlayerReq struct {
	UserID       string
	TeamID       string
	PlayerID     string
	FullName     *string
	DateOfBirth  *time.Time
	JerseyNumber *int32
	Position     *entity.PlayerPosition
	Height       *float32
	Weight       *float32
	Status       *entity.PlayerStatus
}

type UpdatePlayersRes struct {
	PlayerID     uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	DateOfBirth  time.Time
	JerseyNumber int32
	Position     entity.PlayerPosition
	Height       float32
	Weight       float32
	Status       entity.PlayerStatus
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
	ImageUrl     string
}

type PaginateDetails struct {
	TotalPage int32
	Page      int32
	Limit     int32
	Total     int64
}

type GetPlayerReq struct {
	TeamID   string
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
	ID           uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	JerseyNumber int32
	Position     entity.PlayerPosition
	Status       entity.PlayerStatus
	ReleasedAt   time.Time
	JoinedAt     time.Time
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
