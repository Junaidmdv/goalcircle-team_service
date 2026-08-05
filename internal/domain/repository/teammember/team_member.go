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
	IsTeamMemberExist(context.Context, string) (bool, error)
	DeleteTeamMember(context.Context, uuid.UUID, uuid.UUID) error
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

func (tm *teamMemberRepository) IsTeamMemberExist(ctx context.Context, userID string) (bool, error) {
	var count int64
	err := tm.db.WithContext(ctx).
		Model(&entity.TeamMember{}).
		Where("user_id=? ", userID).
		Count(&count).Error
	if err != nil {
		tm.logger.Error("database error", "error", err, "method", "teamMemberRepo.IsTeamMemberExist")
		return false, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return count > 0, nil
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

func (tm *teamMemberRepository) GetTeamMemeberRole(
	ctx context.Context,
	teamID, userID uuid.UUID,
) (entity.TeamMemberRole, error) {

	var member entity.TeamMember

	result := tm.db.WithContext(ctx).
		Select("role").
		Where("team_id = ? AND user_id = ? AND status=?", teamID, userID, entity.TeamMemberStatusActive).
		Take(&member)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", apperror.NewNotFoundError("team member not found")
		}

		tm.logger.Error("database error", "error", result.Error)
		return "", apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	return member.Role, nil
}

func (tm *teamMemberRepository) GetStaffDesignation(ctx context.Context, userID string) (entity.StaffDesignation, error) {

	var result struct {
		Designation entity.StaffDesignation
	}

	err := tm.db.WithContext(ctx).
		Table("team_members").
		Select("staff.designation").
		Joins("JOIN staff ON staff.team_member_id = team_members.id").
		Where("team_members.user_id = ? AND team_members.status=?", userID, entity.TeamMemberStatusActive).
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

func (tr *teamMemberRepository) DeleteTeamMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {

	result := tr.db.WithContext(ctx).Where("user_id=? AND team_id=?", userID, teamID).Delete(&entity.TeamMember{})
	if result.Error != nil {
		tr.logger.Error("database error", "method", "teamMemberRepo.DeleteTeamMember", "error", result.Error)
		return apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFoundError("team member is not found")
	}

	return nil
}
