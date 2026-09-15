package player

import (
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type AddPlayerReq struct {
	UserID        string                `json:"user_id" validate:"required"`
	FullName      string                `json:"full_name" validate:"required"`
	DateOfBirth   time.Time             `json:"dob" validate:"required"`
	Position      entity.PlayerPosition `json:"player_position" validate:"required,player-position"`
	PreferredFoot entity.PreferredFoot  `json:"preferred_foot" validate:"required"`
	Nationality   string                `json:"nationality" validate:"required,len=2"`
	RefreshToken  string                `json:"refresh_token" validate:"required"`
}

type AddPlayerRes struct {
	PlayerID      uuid.UUID
	TeamMemberID  uuid.UUID
	FullName      string
	PlayerPositin entity.PlayerPosition
	PlayerStatus  entity.PlayerStatus
}

type UpdatePlayerReq struct {
	UserID        string                 `json:"user_id" validate:"required"`
	TeamID        string                 `json:"team_id" validate:"required"`
	PlayerID      string                 `json:"player_id" validate:"required"`
	FullName      *string                `json:"full_name,omitempty" validate:"omitempty,min=2,max=100"`
	DateOfBirth   *time.Time             `json:"date_of_birth,omitempty"`
	Position      *entity.PlayerPosition `json:"player_position" validate:"omitempty,player-position"`
	PreferredFoot *entity.PreferredFoot  `json:"preferred_foot" validate:"omitempty"`
	Nationality   *string                `json:"nationality" validate:"omitempty,len=2"`
}

type GetPlayerReq struct {
	PlayerID string `json:"player_id" validate:"required"`
}
