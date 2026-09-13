package staffrepo

import (
	"context"
	"errors"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StaffRepository interface {
	AddStaff(context.Context, *entity.Staff) (*entity.Staff, error)
	UpdateImageKey(context.Context, uuid.UUID, string) error
	GetStaff(ctx context.Context, teamID uuid.UUID, staffID uuid.UUID) (*entity.Staff, error)
	GetStaffImageKey(ctx context.Context, teamID, staffID uuid.UUID) (string, error)
	UpdateStaffImageURL(ctx context.Context, teamID, staffID uuid.UUID) error
	CountByDesination(ctx context.Context, teamID uuid.UUID, desig entity.StaffDesignation) (int32, error)
	UpdateStaff(ctx context.Context, staffID uuid.UUID, staff *entity.Staff) (*entity.Staff, error)
	ListStaff(ctx context.Context, teamID uuid.UUID, details *ListStaffDetails) ([]entity.Staff, int32, error)
	ReleaseStaff(ctx context.Context, teamID uuid.UUID, staffID uuid.UUID) error
	DeleteStaff(ctx context.Context, staffID uuid.UUID) error
	GetStaffDesignation(ctx context.Context, teamID, userID uuid.UUID) (entity.StaffDesignation, error)
	GetStaffByTeamMemberID(ctx context.Context, teamMemberID uuid.UUID) (*entity.Staff, error)
	IsStaffExistByUserID(ctx context.Context, userID uuid.UUID) (bool, error)
}

type staffRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewStaffRepository(db *gorm.DB, logger logger.Logger) StaffRepository {
	return &staffRepository{
		db:     db,
		logger: logger,
	}
}

func (sr *staffRepository) AddStaff(ctx context.Context, staff *entity.Staff) (*entity.Staff, error) {

	if err := sr.db.WithContext(ctx).Create(staff).Error; err != nil {
		sr.logger.Error("database error", "error", err)
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return staff, nil
}

func (sr *staffRepository) GetStaff(ctx context.Context, teamID uuid.UUID, staffID uuid.UUID) (*entity.Staff, error) {

	var staff entity.Staff
	if err := sr.db.WithContext(ctx).Preload("TeamMember").Joins("JOIN team_members tm ON staffs.team_member_id=tm.id").
		Where("staffs.id=? AND tm.team_id=?", staffID, teamID).First(&staff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("staff is not found in this id")
		}
		sr.logger.Error("database error", "error", err, "method", "staff.GetStaff")
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return &staff, nil
}

func (sr *staffRepository) UpdateImageKey(ctx context.Context, staffID uuid.UUID, key string) error {

	result := sr.db.WithContext(ctx).Model(&entity.Staff{}).Where("id=?", staffID).Update("image_key", key)

	if result.Error != nil {
		sr.logger.Error("database error", "error", result.Error, "method", "UpdateImageKey")
		return apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFoundError("failed to update image key")
	}

	return nil
}

func (sr *staffRepository) GetStaffImageKey(ctx context.Context, teamID uuid.UUID, staffID uuid.UUID) (string, error) {

	var key string
	if err := sr.db.WithContext(ctx).Model(&entity.Staff{}).Joins("JOIN team_members tm ON staffs.team_member_id=tm.id").Where("staffs.id=? AND tm.id=?", staffID, teamID).Pluck("staffs.image_key", &key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperror.NewNotFoundError("staff is not found")
		}
		return "", nil
	}
	return key, nil
}

func (sr *staffRepository) CountByDesination(ctx context.Context, teamID uuid.UUID, desig entity.StaffDesignation) (int32, error) {
	var count int64
	result := sr.db.WithContext(ctx).Model(&entity.Staff{}).Joins("JOIN team_members tm ON staffs.team_member_id=tm.id").Where("tm.team_id=? AND staffs.designation=?", teamID, desig).Count(&count)

	if result.Error != nil {
		return -1, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	return int32(count), nil
}

func (sr *staffRepository) UpdateStaff(ctx context.Context, staffID uuid.UUID, staff *entity.Staff) (*entity.Staff, error) {

	result := sr.db.WithContext(ctx).Where("id=?", staffID).Updates(staff)
	if result.Error != nil {
		sr.logger.Error("database error", "error", result.Error, "method", "staffRepo.UpdateStaff")
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, apperror.NewNotFoundError("failed found staff")
	}
	return staff, nil
}

func (sr *staffRepository) ListStaff(ctx context.Context, teamID uuid.UUID, details *ListStaffDetails) ([]entity.Staff, int32, error) {

	var total int64
	var staff []entity.Staff

	query := sr.db.WithContext(ctx).Model(&entity.Staff{}).Preload("TeamMember").Joins("JOIN team_members tm ON staffs.team_member_id=tm.id").Scopes(
		Search(details.Search),
		GetTeamStaff(teamID),
		filterByRole(details.Role),
		filterByDesignation(details.Designation),
	)

	result := query.Count(&total)
	if result.Error != nil {
		sr.logger.Error("database error", "error", result.Error, "method", "repository.Staff.ListStaff")
		return nil, -1, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	err := query.
		Scopes(Paginate(details.Page, details.Limit)).
		Order("created_at DESC").
		Find(&staff).Error

	if err != nil {
		return nil, -1, apperror.NewInternalError("Something went wrong please try again later", err)

	}

	return staff, int32(total), nil
}

func (sr *staffRepository) ReleaseStaff(ctx context.Context, teamID uuid.UUID, staffID uuid.UUID) error {
	result := sr.db.WithContext(ctx).Model(&entity.Staff{}).Where("staffs.id=? AND tm.team_id=?", staffID).Updates(map[string]interface{}{})

	if result.Error == nil {
		sr.logger.Error("database error", "error", result.Error, "method", "repository.Staff.ReleaseStaff")
		return apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFoundError("staff not found in this id")
	}
	return nil
}

func (sr *staffRepository) DeleteStaff(ctx context.Context, staffID uuid.UUID) error {
	result := sr.db.WithContext(ctx).Delete(&entity.Staff{}, staffID)
	if result.Error != nil {
		sr.logger.Error("database error", "method", "staffRepo.DeleteStaff", "error", result.Error)
		return apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFoundError("staff is not found")
	}

	return nil

}

func (tm *staffRepository) GetStaffDesignation(ctx context.Context, teamID, userID uuid.UUID) (entity.StaffDesignation, error) {

	var result struct {
		Designation entity.StaffDesignation
	}

	err := tm.db.WithContext(ctx).
		Model(&entity.Staff{}).
		Select("staffs.designation").
		Joins("JOIN team_members tm ON staffs.team_member_id=tm.id").
		Where("tm.user_id =? AND tm.team_id=? AND tm.status=?", userID, teamID, entity.StaffStatusActive).
		Take(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperror.NewPermissionDeniedError("staff member not found or is not an active member in the team")
		}
		tm.logger.Error("database error", "error", err)
		return "", apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return result.Designation, nil
}

func (tm *staffRepository) GetStaffByTeamMemberID(ctx context.Context, teamMemberID uuid.UUID) (*entity.Staff, error) {
	var staff entity.Staff
	if err := tm.db.WithContext(ctx).Preload("TeamMember").
		Where("staffs.team_member_id=?", teamMemberID).Take(&staff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("staff is not exist")
		}
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return &staff, nil
}

func (r *staffRepository) UpdateStaffImageURL(
	ctx context.Context,
	teamID, staffID uuid.UUID,
) error {
	result := r.db.WithContext(ctx).
		Model(&entity.Staff{}).
		Where(`
			staffs.id = ?
			AND EXISTS (
				SELECT 1
				FROM team_members tm
				WHERE tm.id = staffs.team_member_id
				AND tm.team_id = ?
			)
		`, staffID, teamID).
		Update("image_key", "")

	if result.Error != nil {
		return apperror.NewInternalError(
			"failed to remove staff image",
			result.Error,
		)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFoundError("staff not found")
	}

	return nil
}

func (r *staffRepository) IsStaffExistByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {

	var count int64

	err := r.db.WithContext(ctx).
		Model(&entity.Staff{}).
		Joins("JOIN team_members tm ON staffs.team_member_id = tm.id").
		Where(
			"tm.user_id = ? AND staffs.status = ?",
			userID,
			entity.StaffStatusActive,
		).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
