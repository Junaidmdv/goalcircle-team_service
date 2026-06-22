package player

import (
	"context"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/player"
	teamRepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammemberinvite"
	"github.com/google/uuid"
)

type PlayerUsecase interface {
}

type playerUsecase struct {
	playerRepository player.PlayerRepository
	teamMemberRepo   teammember.TeamMemberRepository
	teamRepository   teamRepo.TeamRepository
	teamInviteRepo   teammemberinvite.TeamMemberInviteRepository
}

func NewPlayerUsecase(player player.PlayerRepository, tmRepo teammember.TeamMemberRepository, teamRepo teamRepo.TeamRepository, tmr teammemberinvite.TeamMemberInviteRepository) PlayerUsecase {
	return &playerUsecase{
		playerRepository: player,
		teamMemberRepo:   tmRepo,
		teamRepository:   teamRepo,
		teamInviteRepo:   tmr,
	}
}

func (pu *playerUsecase) AddNewPlayer(ctx context.Context, input *AddPlayerReq) (*AddPlayerRes, error) {

	teamMember, err := pu.teamMemberRepo.AddTeamMember(ctx, &entity.TeamMember{
		ID:       uuid.New(),
		TeamID:   input.TeamID,
		FullName: input.FullName,
		Role:     entity.PLAYER,
	})

	if err != nil {
		return nil, err
	}

	playerRes, err := pu.playerRepository.CreatePlayer(ctx, &entity.Player{
		ID:           uuid.New(),
		FullName:     input.FullName,
		TeamMemberID: teamMember.ID,
		DateOfBirth:  input.DOB,
		JerseyNumber: input.JerseyNumber,
		Position:     input.Postion,
		Height:       input.Height,
		Weight:       input.Weight,
		Status:       entity.PlayerStatusPendingActionvation,
	})

	if err != nil {
		return nil, err
	}

	teamcode, err := pu.teamRepository.GetTeamCode(ctx, input.TeamID)
	if err != nil {
		return nil, err
	}  

	

	pu.teamInviteRepo.CreateInvitation(ctx,&entity.TeamInvite{
		ID: uuid.New(), 
		TeamID: input.TeamID, 
		TeamCode: teamcode, 
		Code: ,
	})

	return nil, nil
}

func (pu *playerUsecase) RemovePlayer(ctx context.Context, input *RemovePlayerReq) (*RemovePlayerRes, error) {

}

func (pu *playerUsecase) SuspendPlayer(ctx context.Context, input *SuspendPlayerReq) (*SuspendPlayerRes, error) {

}

func (pu *playerUsecase) GetPlayers(ctx context.Context, input *GetTeamPlayersReq) (*GetTeamPlayersRes, error) {
	return nil, nil
}

func (pu *playerUsecase) GetPlayer(ctx context.Context, input *GetPlayerReq) (*GetPlayerRes, error) {
	return nil, nil
}
