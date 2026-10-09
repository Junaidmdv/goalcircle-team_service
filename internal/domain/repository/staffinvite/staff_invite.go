package staff_invite_repo

import (
	"context"
	"errors"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StaffInviteRepository interface {
	AddStaffDetails(ctx context.Context, staff *entity.StaffInvite) (*entity.StaffInvite, error)
	GetStaffDetails(ctx context.Context, staffInviteID uuid.UUID) (*entity.StaffInvite, error)
	GetStaffInvitationByCode(ctx context.Context, code string) (*entity.StaffInvite, error)
	MarkUsedStaffInvitation(ctx context.Context, staffinviteID uuid.UUID) error
}

type staffInviteRepository struct {
	db     gorm.DB
	logger logger.Logger
}

func NewStaffInviteRepository(db gorm.DB, logger logger.Logger) StaffInviteRepository {
	return &staffInviteRepository{
		db:     db,
		logger: logger,
	}
}

func (r *staffInviteRepository) AddStaffDetails(ctx context.Context, staff *entity.StaffInvite) (*entity.StaffInvite, error) {

	if err := r.db.WithContext(ctx).Create(staff).Error; err != nil {
		r.logger.Error("database error", "error", err, "details", "failed to add temp staff details")
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return staff, nil
}

func (r *staffInviteRepository) GetStaffDetails(ctx context.Context, staffInviteID uuid.UUID) (*entity.StaffInvite, error) {

	var details entity.StaffInvite

	err := r.db.WithContext(ctx).
		Model(&entity.StaffInvite{}).
		Where("id = ?", staffInviteID).
		First(&details).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("staff invite details not found")
		}

		return nil, err
	}
	return &details, nil
}

func (r *staffInviteRepository) GetStaffInvitationByCode(ctx context.Context, code string) (*entity.StaffInvite, error) {

	var staffinvite entity.StaffInvite

	if err := r.db.WithContext(ctx).Where("code=?", code).Take(&staffinvite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("code is not exist")
		}
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return &staffinvite, nil
}

func (r *staffInviteRepository) MarkUsedStaffInvitation(ctx context.Context, staffinviteID uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id=?", staffinviteID).Update("is_used", true)
	if result.Error != nil {
		return apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFoundError("failed to update staff")
	}
	return nil
}
