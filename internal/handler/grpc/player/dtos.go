package player

import (
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type AddPlayerReq struct {
	TeamID       string                `json:"team_id" validate:"required"`
	UserID       string                `json:"user_id" validate:"required"`
	FullName     string                `json:"full_name" validate:"required"`
	DateOfBirth  time.Time             `json:"dob" validate:"required"`
	JerseyNumber int32                 `json:"jersey_number" validate:"required"`
	Position     entity.PlayerPosition `json:"player_position" validate:"required,player-position"`
	Height       float32               `json:"height" validate:"required"`
	Weight       float32               `json:"weight" validate:"required"`
}

type AddPlayerRes struct {
	PlayerID      uuid.UUID
	TeamMemberID  uuid.UUID
	FullName      string
	JerseyNumber  int32
	PlayerPositin entity.PlayerPosition
	PlayerStatus  entity.PlayerStatus
}

type UpdatePlayerReq struct {
	UserID       string                 `json:"user_id" validate:"required"`
	TeamID       string                 `json:"team_id" validate:"required"`
	PlayerID     string                 `json:"player_id" validate:"required"`
	FullName     *string                `json:"full_name,omitempty" validate:"omitempty,min=2,max=100"`
	DateOfBirth  *time.Time             `json:"date_of_birth,omitempty"`
	JerseyNumber *int32                 `json:"jersey_number,omitempty" validate:"omitempty,min=1,max=99"`
	Position     *entity.PlayerPosition `json:"position,omitempty"`
	Height       *float32               `json:"height,omitempty" validate:"omitempty,gt=0"`
	Weight       *float32               `json:"weight,omitempty" validate:"omitempty,gt=0"`
	Status       *entity.PlayerStatus   `json:"status,omitempty" validate:"omitempty,player-status"`
}

type GetPlayerReq struct {
	PlayerID string `json:"player_id" validate:"required"`
}
