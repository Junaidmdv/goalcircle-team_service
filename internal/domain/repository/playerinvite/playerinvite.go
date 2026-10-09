package playerinviterepo

import (
	"context"
	"errors"
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlayerInviteRepository interface {
	AddPlayerInviteDetails(ctx context.Context, playerInvite *entity.PlayerTeamInvite) (*entity.PlayerTeamInvite, error)
	GetTeamPlayerInvites(ctx context.Context, paginate *PaginateDetails) ([]entity.PlayerTeamInvite, int64, error)
	GetTeamPlayerInvite(ctx context.Context, inviteID uuid.UUID) (*entity.PlayerTeamInvite, error)
	ExistPendingInvite(ctx context.Context, playerID uuid.UUID) (bool, error)
	UpdatePlayerInviteRes(ctx context.Context, rqID uuid.UUID, status entity.PlayerInvitationStatus) error 
	PlayerInvitationList(ctx context.Context, userID uuid.UUID) ([]entity.PlayerTeamInvite, error)
}

type playerInviteRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewPlayerInviteRepository(db *gorm.DB, logger logger.Logger) PlayerInviteRepository {
	return &playerInviteRepository{
		db:     db,
		logger: logger,
	}
}

func (r *playerInviteRepository) AddPlayerInviteDetails(ctx context.Context, playerinvite *entity.PlayerTeamInvite) (*entity.PlayerTeamInvite, error) {
	if err := r.db.WithContext(ctx).Create(playerinvite).Error; err != nil {
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return playerinvite, nil
}

func (r *playerInviteRepository) GetTeamPlayerInvites(ctx context.Context, paginate *PaginateDetails) ([]entity.PlayerTeamInvite, int64, error) {
	var teamInvites []entity.PlayerTeamInvite
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.PlayerTeamInvite{}).Where("team_id=?", paginate.TeamID).Scopes(
		filterByStatus(paginate.Status),
	)

	result := query.Count(&total)
	if result.Error != nil {
		r.logger.Error("database error", "error", result.Error, "method", "repository.Player.GetTeamPlayers")
		return nil, -1, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	err := query.
		Scopes(Paginate(int32(paginate.Page), int32(paginate.Limit))).
		Order("created_at DESC").
		Find(&teamInvites).Error

	if err != nil {
		return nil, -1, apperror.NewInternalError("Something went wrong please try again later", err)

	}

	return teamInvites, total, nil

}

func (r *playerInviteRepository) GetTeamPlayerInvite(ctx context.Context, inviteID uuid.UUID) (*entity.PlayerTeamInvite, error) {
	var playerInvite entity.PlayerTeamInvite
	if err := r.db.WithContext(ctx).Preload("team_invites").
		Where("ti.team_invite_id=?", inviteID).
		Take(&playerInvite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("player invites is not exist")
		}
		r.logger.Error("database error", "error", err)
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return &playerInvite, nil
}

func (tr *playerInviteRepository) ExistPendingInvite(ctx context.Context, playerID uuid.UUID) (bool, error) {
	var exists bool

	err := tr.db.WithContext(ctx).Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM player_team_invites pi
			WHERE pi.player_id = ?
			  AND pi.status = ?
			  AND pi.expired_at > NOW()
		)`, playerID, entity.InvitationPending,
	).Scan(&exists).Error
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *playerInviteRepository) UpdatePlayerInviteRes(ctx context.Context, rqID uuid.UUID, status entity.PlayerInvitationStatus) error {
	if err := r.db.WithContext(ctx).Where("id=?", rqID).Updates(map[string]interface{}{
		"status":       status,
		"responded_at": time.Now(),
	}).Error; err != nil {
		return apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return nil
}

func (r *playerInviteRepository) PlayerInvitationList(ctx context.Context, userID uuid.UUID) ([]entity.PlayerTeamInvite, error) {
	var playerInvite []entity.PlayerTeamInvite

	result := r.db.WithContext(ctx).Where("user_id=? AND status=?", userID, entity.InvitationPending).Find(&playerInvite)
	if result.Error != nil {
		r.logger.Error("database error", "error", result.Error)
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	return playerInvite, nil
}
