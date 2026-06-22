package teamowner

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type AddTeamOwnerReq struct {
	TeamID   uuid.UUID
	UserId   string
	FullName string
	Role     entity.UserRole
}

type AddTeamOwnerRes struct {
	TeamMemberID uuid.UUID `gorm:"primaryKey"`
	TeamID       uuid.UUID
	UserID       string
	FullName     string
	Role         entity.UserRole
}

type AddPlayerRes struct {
}

type AddPlayerReq struct {
}

type AddStaffReq struct {
}

type AddStaffRes struct {
}

type AddManagerReq struct {
}

type AddManagerRes struct {
}
