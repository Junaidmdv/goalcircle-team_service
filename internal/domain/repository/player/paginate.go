package player

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetTeamPlayers(teamID uuid.UUID) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("tm.team_id=?", teamID)
	}
}

func filterByPlayerStatus(status entity.PlayerStatus) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if status == entity.PlayerStatusInvalid {
			return db
		}
		return db.Where("players.status=?", status)
	}
}

func filterByPosition(position entity.PlayerPosition) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if position == entity.PositionUnspecified {
			return db
		}
		return db.Where("players.position=?", position)
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

		if id, err := uuid.Parse(search); err == nil {
			db = db.Where(
				"players.id = ? OR players.full_name ILIKE ?",
				id,
				"%"+search+"%",
			)
		} else {
			db = db.Where(
				"players.full_name ILIKE ?",
				"%"+search+"%",
			)
		}

		return db
	}
}
