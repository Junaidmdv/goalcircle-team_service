package staffrepo

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Search(search string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if search == "" {
			return db
		}

		if id, err := uuid.Parse(search); err == nil {
			db = db.Where(
				"staffs.id = ? OR staffs.full_name ILIKE ?",
				id,
				"%"+search+"%",
			)
		} else {
			db = db.Where(
				"staffs.full_name ILIKE ?",
				"%"+search+"%",
			)
		}

		return db
	}
}

func Paginate(pageNum, pageLimit int32) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageNum - 1) / pageLimit
		return db.Offset(int(offset)).Limit(int(pageLimit))
	}
}

func filterByRole(role entity.StaffRole) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if role == entity.StaffRoleUnspecified {
			return db
		}
		return db.Where("players.position=?", role)
	}
}

func GetTeamStaff(teamID uuid.UUID) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("tm.team_id=?", teamID)
	}
}

func filterByDesignation(desig entity.StaffDesignation) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if desig == entity.StaffDesignationUnspecified {
			return db
		}
		return db.Where("players.position=?", desig)
	}
}
