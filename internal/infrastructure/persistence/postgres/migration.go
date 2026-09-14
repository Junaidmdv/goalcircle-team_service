package postgres

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
)

func (db *postgresDB) Migration() error {
	if err := db.DB.AutoMigrate(
		&entity.Team{},
		&entity.TeamMember{},
		&entity.Player{},
		&entity.TeamInvite{},
		&entity.Staff{},
		&entity.StaffInviteDetails{}, 
		&entity.TeamPlayer{},
	); err != nil {
		return err
	}
	return nil
}
