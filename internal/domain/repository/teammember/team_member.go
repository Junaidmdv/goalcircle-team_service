package teammember

import (
	"context"
	"errors"

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
	GetTeamMemeberRole(context.Context, uuid.UUID, uuid.UUID) (entity.TeamMemberRole, error)
	GetStaffDesignation(context.Context, string) (entity.StaffDesignation, error)
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

func (tm *teamMemberRepository) GetTeamMemeberRole(ctx context.Context, teamID, teamMemberID uuid.UUID) (entity.TeamMemberRole, error) {
	var role entity.TeamMemberRole
	result := tm.db.WithContext(ctx).
		Model(&entity.TeamMember{}).
		Select("role").
		Where("team_id = ? AND id = ?", teamID, teamMemberID).
		Scan(&role)

	if result.Error != nil {
		tm.logger.Error("database error", "error", result.Error)
		return "", apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return "", apperror.NewNotFoundError("team member not found")
	}

	return role, nil
}

func (tm *teamMemberRepository) GetStaffDesignation(
	ctx context.Context,
	userID string,
) (entity.StaffDesignation, error) {

	var result struct {
		Designation entity.StaffDesignation
	}

	err := tm.db.WithContext(ctx).
		Table("team_members").
		Select("staff.designation").
		Joins("JOIN staff ON staff.team_member_id = team_members.id").
		Where("team_members.user_id = ?", userID).
		Take(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperror.NewNotFoundError("staff not found")
		}

		tm.logger.Error("database error", "error", err)
		return "", apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return result.Designation, nil
}
