package teaminvite

import (
	"context"
	"errors"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamInviteRepository interface {
	CreateInvitation(context.Context, *entity.TeamInvite) (*entity.TeamInvite, error)
	GetInvitation(context.Context, *uuid.UUID) (*entity.TeamInvite, bool, error)
}

type teamInviteRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewTeamMemberInviteRepository(db *gorm.DB) TeamInviteRepository {
	return &teamInviteRepository{
		db: db,
	}
}

func (tr *teamInviteRepository) CreateInvitation(ctx context.Context, invite *entity.TeamInvite) (*entity.TeamInvite, error) {
	if err := tr.db.WithContext(ctx).Create(invite).Error; err != nil {
		tr.logger.Error("databaser error", "eroror", err)
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return invite, nil
}

func (tr *teamInviteRepository) GetInvitation(ctx context.Context, teamMemberID *uuid.UUID) (*entity.TeamInvite, bool, error) {

	var invite entity.TeamInvite
	result := tr.db.WithContext(ctx).Where("team_member_id=?", teamMemberID).Take(&invite)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		tr.logger.Error("database error", "error", result.Error, "method", "GetInvitation")
		return nil, false, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)

	}

	return &invite, true, nil

}
