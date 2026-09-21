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
	GetPlayers(context.Context, *ListUserReq) ([]entity.Player, int64, error)
	GetPlayer(context.Context, uuid.UUID) (*entity.Player, error)
	GetPlayerByUserID(ctx context.Context, userID uuid.UUID) (*entity.Player, error)
	UpdateImageKey(context.Context, uuid.UUID, string) error
	IsPlayerExist(context.Context, uuid.UUID) (bool, error)
	GetPlayerImageKey(ctx context.Context, playerID uuid.UUID) (string, error)
	UpdatePlayerDetails(ctx context.Context, playerID uuid.UUID, player *entity.Player) error
	DeletePlayer(ctx context.Context, playerID uuid.UUID) error
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


func (pr *playerRepository) GetPlayers(ctx context.Context, details *ListUserReq) ([]entity.Player, int64, error) {
	var total int64
	var players []entity.Player

	query := pr.db.WithContext(ctx).Model(&entity.Player{}).Scopes(
		Search(details.Search),
		filterByPosition(details.Position),
	)

	result := query.Count(&total)
	if result.Error != nil {
		pr.logger.Error("database error", "error", result.Error, "method", "repository.Player.GetTeamPlayers")
		return nil, -1, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	err := query.
		Scopes(Paginate(details.Page, details.Limit)).
		Order("created_at DESC").
		Find(&players).Error

	if err != nil {
		return nil, -1, apperror.NewInternalError("Something went wrong please try again later", err)

	}

	return players, total, nil
}

func (pr *playerRepository) GetPlayer(ctx context.Context, playerID uuid.UUID) (*entity.Player, error) {
	var player entity.Player

	if err := pr.db.WithContext(ctx).First(&player, playerID).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, apperror.NewNotFoundError("player not exist")
		}
	}

	return &player, nil
}

func (pr *playerRepository) GetPlayerByUserID(ctx context.Context, userID uuid.UUID) (*entity.Player, error) {
	var player entity.Player

	if err := pr.db.WithContext(ctx).Where("user_id=?", userID).Take(&player).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("player is not found")
		}

		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return &player, nil
}



func (pr *playerRepository) UpdateImageKey(ctx context.Context, playerID uuid.UUID, key string) error {

	result := pr.db.WithContext(ctx).Model(&entity.Player{}).Where("id=?", playerID).Update("profile_image_key", key)

	if result.Error != nil {
		pr.logger.Error(
			"database failure",
			"method", "teamrepo.UpdateImageKey",
			"error", result.Error,
		)
		return apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFoundError("player not found")
	}

	return nil
}

func (pr *playerRepository) IsPlayerExist(ctx context.Context, userID uuid.UUID) (bool, error) {
	var player entity.Player

	if err := pr.db.WithContext(ctx).First(&player, "user_id=?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		pr.logger.Error("database error", "method", "teamrepo.IsTeamExist", "error", err)
		return false, apperror.NewInternalError(apperror.InternalErrorMsg, err)

	}

	return true, nil
}

func (pr *playerRepository) GetPlayerImageKey(ctx context.Context, playerID uuid.UUID) (string, error) {
	var imageKey string

	result := pr.db.WithContext(ctx).
		Model(&entity.Player{}).
		Where("id = ?", playerID).
		Pluck("image_key", &imageKey)

	if result.Error != nil {
		return "", result.Error
	}

	if result.RowsAffected == 0 {
		return "", gorm.ErrRecordNotFound
	}

	return imageKey, nil
}

func (pr *playerRepository) UpdatePlayerDetails(ctx context.Context, playerID uuid.UUID, player *entity.Player) error {
	result := pr.db.WithContext(ctx).
		Model(&entity.Player{}).
		Where("id = ?", playerID).
		Updates(player)

	if result.Error != nil {
		return apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFoundError("failed to find player id")
	}
	return nil
}

func (pr *playerRepository) DeletePlayer(ctx context.Context, playerID uuid.UUID) error {

	result := pr.db.WithContext(ctx).Delete(&entity.Player{}, playerID)

	if result.Error != nil {
		return apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFoundError("player not found")
	}

	return nil
}
