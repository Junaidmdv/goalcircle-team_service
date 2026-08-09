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

type UpdatePlayerStatusReq struct {
	PlayerID     string              `json:"full_name" validate:"required"`
	PlayerStatus entity.PlayerStatus `json:"player_status" validate:"required,player-status"`
}

type GetPlayerReq struct {
	PlayerID string `json:"player_id" validate:"required"`
}
