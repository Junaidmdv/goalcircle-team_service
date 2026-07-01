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
	GetInvitationByCode(context.Context, string) (*entity.TeamInvite, error)
	UpdateTeamInvite(context.Context, uuid.UUID, *UpdateTeamInviteReq) error
}

type teamInviteRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewTeamMemberInviteRepository(db *gorm.DB, logger logger.Logger) TeamInviteRepository {
	return &teamInviteRepository{
		db:     db,
		logger: logger,
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

func (tr *teamInviteRepository) GetInvitationByCode(ctx context.Context, code string) (*entity.TeamInvite, error) {

	var invite entity.TeamInvite

	result := tr.db.WithContext(ctx).Where("code=?", code).Take(&invite)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("invalid invite code")
		}
		tr.logger.Error("database error", "error", result.Error, "method", "GetInvitation")
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)

	}

	return &invite, nil
}

func (tr *teamInviteRepository) UpdateTeamInvite(ctx context.Context, id uuid.UUID, req *UpdateTeamInviteReq) error {
	if err := tr.db.WithContext(ctx).Where("id=?", id).Updates(map[string]any{
		"used_id": req.UserID,
		"is_used": req.IsUsed,
	}).Error; err != nil {
		tr.logger.Error("database error", "error", err, "method", "repo.updateTeamInvite")
		return apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return nil
}
