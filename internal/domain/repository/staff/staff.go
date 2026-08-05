package staffrepo

import (
	"context"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"gorm.io/gorm"
)

type StaffRepository interface {
	AddStaff(context.Context, *entity.Staff) (*entity.Staff, error)
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
