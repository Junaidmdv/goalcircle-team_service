package teamplayer

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type ListPlayerDetailsReq struct {
	TeamID       uuid.UUID
	Page         int32
	Limit        int32
	PlayerStatus entity.TeamPlayerStatus
	Position     entity.PlayerPosition
	Search       string
}
