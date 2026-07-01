package teammember

import (
	"context"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamMemberRepository interface {
	AddTeamMember(context.Context, *entity.TeamMember) (*entity.TeamMember, error)
	RemoveTeamMember(context.Context, *uuid.UUID) error
	UpdateUserID(context.Context, *uuid.UUID, string) error
}

type teamMemberRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewTeamMemberRepository(db *gorm.DB, logger logger.Logger) TeamMemberRepository {
	return &teamMemberRepository{
		db:     db,
		logger: logger,
	}
}

func (tm *teamMemberRepository) AddTeamMember(ctx context.Context, input *entity.TeamMember) (*entity.TeamMember, error) {
	if err := tm.db.WithContext(ctx).Create(input).Error; err != nil {
		tm.logger.Error("database error", "error", err)
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return input, nil
}

func (tm *teamMemberRepository) RemoveTeamMember(ctx context.Context, teamteamMemberID *uuid.UUID) error {
	if err := tm.db.WithContext(ctx).Unscoped().Delete(&entity.TeamMember{}, teamteamMemberID).Error; err != nil {
		tm.logger.Error("database error", "error", err)
		return apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return nil
}

func (tm *teamMemberRepository) UpdateUserID(ctx context.Context, teamMemberID *uuid.UUID, userId string) error {

	if err := tm.db.WithContext(ctx).Where("team_member_id=?", teamMemberID).Update("user_id", userId).Error; err != nil {
		tm.logger.Error("database error", "error", err)
		return apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return nil
}
