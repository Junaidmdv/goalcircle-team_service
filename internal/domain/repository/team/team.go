package team

import (
	"context"
	"errors"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository interface {
	CreateTeam(context.Context, *entity.Team) (*entity.Team, error)
	DeleteTeam(context.Context, uuid.UUID) error
	GetTeamCode(context.Context, uuid.UUID) (string, error)
}

type teamRepository struct {
	db      *gorm.DB
	loggger logger.Logger
}

func NewTeamRepository(db *gorm.DB, logger logger.Logger) TeamRepository {
	return &teamRepository{
		db:      db,
		loggger: logger,
	}
}

func (tr *teamRepository) CreateTeam(ctx context.Context, team *entity.Team) (*entity.Team, error) {
	if err := tr.db.WithContext(ctx).Create(team).Error; err != nil {
		tr.loggger.Error("database error", "error", err, "method", "CreateTeam")
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return team, nil
}

func (tr *teamRepository) DeleteTeam(ctx context.Context, teamid uuid.UUID) error {
	if err := tr.db.WithContext(ctx).Unscoped().Where("id=?", teamid).Delete(&entity.Team{}).Error; err != nil {
		tr.loggger.Error("database error", "error", err)
		return apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return nil
}

func (tr *teamRepository) GetTeamCode(ctx context.Context, teamid uuid.UUID) (string, error) {
	var code string
	if err := tr.db.WithContext(ctx).Select("code").First(&code, teamid).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperror.NewNotFoundError("team data is not found")
		}
		tr.loggger.Error("databse error", "error", err)
		return "", apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return code, nil
}



func(tr *teamRepository)UpdateTeamDetails(ctx context.Context,teamID uuid.UUID,)