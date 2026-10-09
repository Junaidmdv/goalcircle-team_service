package playerinviterepo

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateStatus struct {
	Status entity.PlayerInvitationStatus 
	OfferedJerseyNum int32 
}

type PaginateDetails struct {
	TeamID uuid.UUID
	Page   int
	Limit  int
	Status entity.PlayerInvitationStatus
}

func filterByStatus(status entity.PlayerInvitationStatus) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if status == entity.InvitationInvalid {
			return db
		}
		return db.Where("team_invites.status=?", status)
	}
}

func Paginate(pageNum, pageLimit int32) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageNum - 1) / pageLimit
		return db.Offset(int(offset)).Limit(int(pageLimit))
	}
}
