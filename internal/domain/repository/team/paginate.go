package team

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"gorm.io/gorm"
)

func FilterByTeamStatus(status entity.TeamStatus) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if status == "" {
			return db
		}

		return db.Where("team_status = ?", status)
	}
}

func FilterByCity(city string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if city == "" {
			return db
		}

		return db.Where("city ILIKE ?", "%"+city+"%")
	}
}

func SearchTeam(search string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if search == "" {
			return db
		}

		search = "%" + search + "%"

		return db.Where(
			"name ILIKE ? OR short_name ILIKE ? OR team_code ILIKE ?",
			search,
			search,
			search,
		)
	}
}



func Paginate(page, pageSize int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
	

		offset := (page - 1) * pageSize

		return db.Offset(offset).Limit(pageSize)
	}
}
