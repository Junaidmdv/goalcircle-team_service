package teamplayer

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

type TeamPlayerRepository interface {
	AddTeamPlayer(ctx context.Context, player *entity.TeamPlayer) (*entity.TeamPlayer, error)
	IsTeamPlayerExist(ctx context.Context, userID uuid.UUID) (bool, error)
	GetTeamPlayers(ctx context.Context, details *ListPlayerDetailsReq) ([]entity.TeamPlayer, int64, error)
	GetTeamPlayer(ctx context.Context, teamID uuid.UUID, playerID uuid.UUID) (*entity.TeamPlayer, error)
	ReleaseTeamPlayer(ctx context.Context, teamID uuid.UUID, plyerID uuid.UUID) error
	IsPlayerPartOfTeam(ctx context.Context, teamID, playerID uuid.UUID) (bool, error)
	IsJerseyNumberTaken(ctx context.Context, teamID uuid.UUID, jerseyNumber *int32) (bool, error)
}

type teamPlayerRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewTeamPlayerRepository(db *gorm.DB, logger logger.Logger) TeamPlayerRepository {
	return &teamPlayerRepository{
		db:     db,
		logger: logger,
	}
}

func (tp *teamPlayerRepository) IsTeamPlayerExist(ctx context.Context, userID uuid.UUID) (bool, error) {

	var count int64

	err := tp.db.WithContext(ctx).
		Model(&entity.TeamPlayer{}).
		Joins("JOIN team_members tm ON team_players.team_member_id=tm.id").
		Where(
			"tm.user_id = ? AND team_players.status <> ? AND team_players.left_at IS NULL",
			userID,
			entity.TeamPlayerStatusLeft,
		).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (pr *teamPlayerRepository) GetTeamPlayers(ctx context.Context, details *ListPlayerDetailsReq) ([]entity.TeamPlayer, int64, error) {
	var total int64
	var users []entity.TeamPlayer

	query := pr.db.WithContext(ctx).Model(&entity.TeamPlayer{}).Preload("Player").Joins("JOIN players p ON team_players.player_id=id").
		Joins("JOIN team_members tm ON team_players.team_member_id=tm.id").
		Scopes(
			Search(details.Search),
			GetTeamPlayers(details.TeamID),
			filterByPlayerStatus(details.PlayerStatus),
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
		Find(&users).Error

	if err != nil {
		return nil, -1, apperror.NewInternalError("Something went wrong please try again later", err)

	}

	return users, total, nil
}

func (pr *teamPlayerRepository) GetTeamPlayer(
	ctx context.Context,
	teamID uuid.UUID,
	playerID uuid.UUID,
) (*entity.TeamPlayer, error) {

	var teamPlayer entity.TeamPlayer

	err := pr.db.WithContext(ctx).
		Preload("TeamMember").
		Preload("Player").
		Joins("JOIN team_members tm ON team_players.team_member_id = tm.id").
		Where(
			"team_players.player_id = ? AND tm.team_id = ?",
			playerID,
			teamID,
		).
		Take(&teamPlayer).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("player not found")
		}

		return nil, apperror.NewInternalError(
			apperror.InternalErrorMsg,
			err,
		)
	}

	return &teamPlayer, nil
}

func (tp *teamPlayerRepository) ReleaseTeamPlayer(
	ctx context.Context,
	teamID uuid.UUID,
	playerID uuid.UUID,
) error {

	result := tp.db.WithContext(ctx).
		Joins("JOIN team_members tm ON team_players.team_member_id = tm.id").
		Where(
			"tm.team_id = ? AND team_players.player_id = ?",
			teamID,
			playerID,
		).
		Updates(map[string]interface{}{
			"left_at": time.Now(),
			"status":  entity.TeamPlayerStatusLeft,
		})

	if result.Error != nil {
		return apperror.NewInternalError(
			apperror.InternalErrorMsg,
			result.Error,
		)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFoundError("player not found in team")
	}

	return nil
}

func (tp *teamPlayerRepository) IsPlayerPartOfTeam(ctx context.Context, teamID, playerID uuid.UUID) (bool, error) {

	var exists bool
	result := tp.db.Raw("SELECT EXISTS(SELECT 1 FROM team_players WHERE player_id = ? AND team_id=?)", playerID, teamID).Scan(&exists)

	if result.Error != nil {
		return false, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	return exists, nil
}

func (tp *teamPlayerRepository) IsJerseyNumberTaken(ctx context.Context, teamID uuid.UUID, jerseyNumber *int32) (bool, error) {
	var exists bool
	result := tp.db.Raw("SELECT EXISTS(SELECT 1 FROM team_players WHERE team_id = ? AND jersey_number=?)", teamID, jerseyNumber).Scan(&exists)

	if result.Error != nil {
		return false, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	return exists, nil
}

func (r *teamPlayerRepository) AddTeamPlayer(ctx context.Context, teamPlayer *entity.TeamPlayer) (*entity.TeamPlayer, error) {
	return nil, nil
}
