package team

import (
	"context"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	team_repo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
)

type TeamUsecase interface {
	CreateTeam(context.Context, *CreateTeamReq) (*CreateTeamRes, error)
	DeleteTeam(context.Context, uuid.UUID) error
}

type teamUsecase struct {
	teamRepo team_repo.TeamRepository
	logger   logger.Logger
}

func NewTeamUsecase(teamrepo team_repo.TeamRepository, logger logger.Logger) TeamUsecase {
	return &teamUsecase{
		teamRepo: teamrepo,
		logger:   logger,
	}
}

func (tu *teamUsecase) CreateTeam(ctx context.Context, dt *CreateTeamReq) (*CreateTeamRes, error) {

	code, err := GenerateCode(entity.TeamCodeLenght)
	if err != nil {
		tu.logger.Error("Failed to generate code", "error", err, "method", "teamUsercase")
	}

	res, err := tu.teamRepo.CreateTeam(ctx, &entity.Team{
		ID:          uuid.New(),
		Name:        dt.Name,
		ShortName:   GenerateShortName(dt.Name),
		City:        dt.City,
		Description: dt.Description,
		TeamCode:    code,
	})
	if err != nil {
		return nil, err
	}

	return &CreateTeamRes{
		TeamID:    res.ID,
		Name:      res.Name,
		ShortName: res.Name,
		City:      res.City,
		Code:      res.TeamCode,
	}, nil
}

func (tu *teamUsecase) DeleteTeam(ctx context.Context, teamId uuid.UUID) error {
	return tu.teamRepo.DeleteTeam(ctx, teamId)
}

