package team

import (
	"context"

	team_repo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
)

type TeamUsecase interface {
}


type teamUsecase struct {
	teamRepo team_repo.TeamRepository
}

func NewTeamUsecase(teamrepo team_repo.TeamRepository) TeamUsecase {
	return &teamUsecase{
        teamRepo: teamrepo,
	}
}

func (tu *teamUsecase) CreateTeam(ctx context.Context, dt *CreateTeamReq) (*CreateTeamRes, error) { 
  return nil,nil
}


