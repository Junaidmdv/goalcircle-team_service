package postgres

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
)

func (db *postgresDB) Migration() error {
	if err := db.DB.AutoMigrate(
		&entity.Team{},
		&entity.TeamMember{},
		&entity.Player{},
		&entity.Staff{},
		&entity.StaffInvite{}, 
		&entity.TeamPlayer{}, 
		&entity.Player{}, 
		&entity.PlayerTeamInvite{}, 
	); err != nil {
		return err
	}
	return nil
}
