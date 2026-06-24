package player

import (
	"context"
	"errors"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlayerRepository interface {
	CreatePlayer(context.Context, *entity.Player) (*entity.Player, error)
	UpdatePlayerStatus(context.Context, *uuid.UUID, *entity.PlayerStatus) error
	GetTeamPlayers(context.Context, *ListUserReq) ([]entity.Player, int64, error)
	GetPlayer(context.Context, *uuid.UUID) (*entity.Player, error)
}

type playerRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewPlayerRepository(db *gorm.DB, logger logger.Logger) PlayerRepository {
	return &playerRepository{
		db:     db,
		logger: logger,
	}
}

func (pr *playerRepository) CreatePlayer(ctx context.Context, player *entity.Player) (*entity.Player, error) {

	if err := pr.db.WithContext(ctx).Create(player).Error; err != nil {
		pr.logger.Error("failed create player", "error", err, "method", "player.CreatePlayer")
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return player, nil
}

func (pr *playerRepository) UpdatePlayerStatus(ctx context.Context, playerID *uuid.UUID, status *entity.PlayerStatus) error {

	if err := pr.db.WithContext(ctx).Where("id=?", playerID).Update("status", status).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NewNotFoundError("player data is not exist")
		}
		pr.logger.Error("database error", "error", err, "method", "repository.Player.UpdatePlayerStatus")
		return apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return nil
}

func (pr *playerRepository) GetTeamPlayers(ctx context.Context, details *ListUserReq) ([]entity.Player, int64, error) {
	var total int64
	var users []entity.Player

	query := pr.db.WithContext(ctx).Scopes(
		Search(details.Search),
		GetTeamPlayers(details.TeamID),
		filterByPlayerStatus(details.PlayerStatus),
		filterByPosition(details.Position),
	)

	result := query.Count(&total)
	if result.Error != nil {
		pr.logger.Error("database error", "error", result.Error, "method", "repository.Player.GetTeamPlayers")
		return users, -1, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	return users, -1, nil
}

func (pr *playerRepository) GetPlayer(ctx context.Context, playerID *uuid.UUID) (*entity.Player, error) {
	var player entity.Player
	if err := pr.db.WithContext(ctx).First(&player, playerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("player is not found")
		}
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return &player, nil
}
