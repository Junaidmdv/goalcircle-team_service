package team

import (
	"context"
	"errors"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TeamRepository interface {
	CreateTeam(context.Context, *entity.Team) (*entity.Team, error)
	DeleteTeam(context.Context, uuid.UUID) error
	GetTeamCode(context.Context, uuid.UUID) (string, error)
	UpdateTeamDetails(context.Context, uuid.UUID, *UpdateTeamReq) (*entity.Team, error)
	UpdateTeamContactDetails(context.Context, uuid.UUID, *UpdateTeamContact) (*entity.Team, error)
	SetCaptain(context.Context, uuid.UUID, uuid.UUID) error
	SetViceCaptain(context.Context, uuid.UUID, uuid.UUID) error
	ListTeams(context.Context, *ListTeamsReq) ([]entity.Team, int32, error)
	IsTeamExist(context.Context, uuid.UUID) (bool, error)
	UpdateLogoKey(context.Context, uuid.UUID, string) error
	IsJerseyNumOccupied(context.Context, uuid.UUID, int32) (bool, error)
	GetTeamDetails(context.Context, uuid.UUID) (*entity.Team, error)
}

type teamRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewTeamRepository(db *gorm.DB, logger logger.Logger) TeamRepository {
	return &teamRepository{
		db:     db,
		logger: logger,
	}
}

func (tr *teamRepository) CreateTeam(ctx context.Context, team *entity.Team) (*entity.Team, error) {
	if err := tr.db.WithContext(ctx).Create(team).Error; err != nil {
		tr.logger.Error("database error", "error", err, "method", "CreateTeam")
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return team, nil
}

func (tr *teamRepository) DeleteTeam(ctx context.Context, teamid uuid.UUID) error {
	if err := tr.db.WithContext(ctx).Unscoped().Where("id=?", teamid).Delete(&entity.Team{}).Error; err != nil {
		tr.logger.Error("database error", "error", err)
		return apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return nil
}

func (tr *teamRepository) GetTeamCode(ctx context.Context, teamid uuid.UUID) (string, error) {
	var code string
	if err := tr.db.WithContext(ctx).Select("code").First(&code, teamid).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperror.NewNotFoundError("team data is not found")
		}
		tr.logger.Error("databse error", "error", err)
		return "", apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return code, nil
}

func (tr *teamRepository) UpdateTeamDetails(ctx context.Context, teamID uuid.UUID, req *UpdateTeamReq) (*entity.Team, error) {

	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if req.ShortName != nil {
		updates["short_name"] = *req.ShortName
	}

	if req.City != nil {
		updates["city"] = *req.City
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if req.Email != nil {
		updates["email"] = *req.Email
	}

	if req.PhoneNum != nil {
		updates["phone_num"] = *req.PhoneNum
	}

	if len(updates) == 0 {
		return nil, apperror.NewInvalidArgumentError("no fields to update")
	}

	var team entity.Team
	result := tr.db.WithContext(ctx).
		Model(&team).
		Clauses(clause.Returning{}).
		Where("id = ?", teamID).
		Updates(updates)
	if result.Error != nil {
		tr.logger.Error("database error", "error", result.Error)
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, apperror.NewNotFoundError("team not found")
	}
	return &team, nil
}

func (tr *teamRepository) UpdateTeamContactDetails(ctx context.Context, teamID uuid.UUID, req *UpdateTeamContact) (*entity.Team, error) {
	updates := make(map[string]interface{})
	if req.ContactEmail != nil {
		updates["contact_email"] = *req.ContactEmail
	}

	if req.ContactPhoneNum != nil {
		updates["contact_phone_num"] = *req.ContactPhoneNum
	}
	var team entity.Team
	result := tr.db.WithContext(ctx).
		Model(&team).
		Clauses(clause.Returning{}).
		Where("id = ?", teamID).
		Updates(updates)
	if result.Error != nil {
		tr.logger.Error("database error", "error", result.Error)
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, apperror.NewNotFoundError("team not found")
	}
	return &team, nil

}

func (tm *teamRepository) SetCaptain(ctx context.Context, teamID uuid.UUID, playerID uuid.UUID) error {
	result := tm.db.WithContext(ctx).
		Model(&entity.Team{}).
		Where("id = ?", teamID).
		Update("captain_id", playerID)

	if result.Error != nil {
		tm.logger.Error("databaser error", "error", result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFoundError("team not found")
	}
	return nil
}

func (tm *teamRepository) SetViceCaptain(ctx context.Context, teamID uuid.UUID, playerID uuid.UUID) error {
	result := tm.db.WithContext(ctx).
		Model(&entity.Team{}).
		Where("id = ?", teamID).
		Update("vice_captain_id", playerID)

	if result.Error != nil {
		tm.logger.Error("databaser error", "error", result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFoundError("team not found")
	}
	return nil
}

func (tm *teamRepository) ListTeams(ctx context.Context, req *ListTeamsReq) ([]entity.Team, int32, error) {
	var total int64
	var teams []entity.Team

	query := tm.db.WithContext(ctx).Scopes(
		SearchTeam(req.Search),
		FilterByCity(req.City),
		FilterByTeamStatus(req.TeamStatus),
	)
	result := query.Count(&total)
	if result.Error != nil {
		tm.logger.Error("database error", "error", result.Error, "method", "repository.Player.GetTeamPlayers")
		return teams, -1, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	err := query.
		Scopes(Paginate(int(req.Page), int(req.Limit))).
		Order("created_at DESC").
		Find(&teams).Error

	if err != nil {
		return nil, -1, apperror.NewInternalError("Something went wrong please try again later", err)
	}
	return teams, -1, nil
}

func (tr *teamRepository) IsTeamExist(ctx context.Context, teamID uuid.UUID) (bool, error) {

	var team entity.Team

	if err := tr.db.WithContext(ctx).First(&team, "id=?", teamID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		tr.logger.Error("database error", "method", "teamrepo.IsTeamExist", "error", err)
		return false, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return true, nil
}

func (tr *teamRepository) UpdateLogoKey(ctx context.Context, teamID uuid.UUID, key string) error {
	result := tr.db.WithContext(ctx).
		Model(&entity.Team{}).
		Where("id = ?", teamID).
		Update("logo_key", key)

	if result.Error != nil {
		tr.logger.Error(
			"database failure",
			"method", "teamrepo.UpdateLogoKey",
			"error", result.Error,
		)
		return apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFoundError("team not found")
	}

	return nil
}

func (pr *teamRepository) IsJerseyNumOccupied(
	ctx context.Context,
	teamID uuid.UUID,
	jerseyNumber int32,
) (bool, error) {

	var player entity.Player

	err := pr.db.WithContext(ctx).
		Model(&entity.Player{}).
		Joins("JOIN team_members tm ON tm.id = players.team_member_id").
		Where(
			"tm.team_id = ? AND players.jersey_number = ?  ",
			teamID,
			jerseyNumber,
		).Not("players.status = ?", entity.PlayerStatusArchieved).
		First(&player).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		pr.logger.Error(
			"database error",
			"error", err,
			"method", "IsJerseyNumOccupied",
		)

		return false, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return true, nil
}

func (tr *teamRepository) GetTeamDetails(ctx context.Context, teamID uuid.UUID) (*entity.Team, error) {
	var team entity.Team

	if err := tr.db.WithContext(ctx).First(&team, teamID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("team not found in this id")
		}
		tr.logger.Error("database error", "error", err, "method", "teamRepo.GetTeamDetails")
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return &team, nil
}
