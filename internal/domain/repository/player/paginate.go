package player

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetTeamPlayers(teamID uuid.UUID) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("team_id=?", teamID)
	}
}

func filterByPlayerStatus(status entity.PlayerStats) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status=?", status)
	}
}

func filterByPosition(position entity.PlayerPosition) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("position=?", position)
	}
}

func Paginate(pageNum, pageLimit int32) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageNum - 1) / pageLimit
		return db.Offset(int(offset)).Limit(int(pageLimit))
	}
}

func Search(search string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if search == "" {
			return db
		}

		search = "%" + search + "%"

		return db.Where(
			"id ILIKE ? OR full_name ILIKE ?",
			search,
			search,
		)
	}
}
