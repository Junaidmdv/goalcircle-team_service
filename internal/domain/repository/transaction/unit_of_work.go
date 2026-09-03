package transaction

import (
	"context"

	staffrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staff"
	teamrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
	teammemberrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"gorm.io/gorm"
)

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(Repositories) error) error
}

type unitOfWork struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewUnitOfWork(db *gorm.DB, logger logger.Logger) UnitOfWork {
	return &unitOfWork{
		db:     db,
		logger: logger,
	}
}

func (u *unitOfWork) Execute(ctx context.Context, fn func(Repositories) error) error {

	return u.db.WithContext(ctx).Transaction(func(db *gorm.DB) error {

		repos := &repositories{
			teamRepo:       teamrepo.NewTeamRepository(db, u.logger),
			staffRepo:      staffrepo.NewStaffRepository(db, u.logger),
			teamMemberRepo: teammemberrepo.NewTeamMemberRepository(db, u.logger),
		}

		return fn(repos)
	})

}
