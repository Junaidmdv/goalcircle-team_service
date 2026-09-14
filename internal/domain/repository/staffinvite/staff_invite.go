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
	AddStaffDetails(ctx context.Context, staff *entity.StaffInviteDetails) error
	GetStaffDetails(ctx context.Context, teamInviteID uuid.UUID) (*entity.StaffInviteDetails, error)
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

func (r *staffInviteRepository) AddStaffDetails(ctx context.Context, staff *entity.StaffInviteDetails) error {

	if err := r.db.WithContext(ctx).Create(staff).Error; err != nil {
		r.logger.Error("database error", "error", err, "details", "failed to add temp staff details")
		return apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return nil
}

func (r *staffInviteRepository) GetStaffDetails(ctx context.Context, teamInviteID uuid.UUID) (*entity.StaffInviteDetails, error) {

	var details entity.StaffInviteDetails

	err := r.db.WithContext(ctx).
		Model(&entity.StaffInviteDetails{}).
		Joins("JOIN team_invites ti ON ti.id = staff_invite_details.invite_id").
		Where("ti.id = ?", teamInviteID).
		First(&details).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("staff invite details not found")
		}

		return nil, err
	}
	return &details, nil
}
