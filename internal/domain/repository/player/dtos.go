package player

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type ListUserReq struct {
	TeamID       uuid.UUID
	Page         int32
	Limit        int32
	PlayerStatus entity.PlayerStatus
	Position     entity.PlayerPosition
	Search       string
}
