package player

import (
	"context"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/player"
	teamRepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	code "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/invitation"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
)

type PlayerUsecase interface {
	AddNewPlayer(context.Context, *AddPlayerReq) (*AddPlayerRes, error)
	UpdatePlayerStatus(context.Context, *UpdatPlayerStatusReq) (*UpdatePlayerStatusRes, error)
	ListTeamPlayers(context.Context, *ListTeamPlayersReq) ([]PlayerRes, *PaginateDetails, error)
	GetPlayer(context.Context, *GetPlayerReq) (*GetPlayerRes, error)
}

type playerUsecase struct {
	playerRepository player.PlayerRepository
	teamMemberRepo   teammember.TeamMemberRepository
	teamRepository   teamRepo.TeamRepository
	code             code.CodeGenerater
	logger           logger.Logger
}

func NewPlayerUsecase(player player.PlayerRepository,
	tmRepo teammember.TeamMemberRepository,
	teamRepo teamRepo.TeamRepository,

	code code.CodeGenerater) PlayerUsecase {

	return &playerUsecase{
		playerRepository: player,
		teamMemberRepo:   tmRepo,
		teamRepository:   teamRepo,
		code:             code,
	}
}

func (pu *playerUsecase) AddNewPlayer(ctx context.Context, input *AddPlayerReq) (*AddPlayerRes, error) {

	teamID, err := uuid.Parse(input.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid player id")

	}

	teamMember, err := pu.teamMemberRepo.AddTeamMember(ctx, &entity.TeamMember{
		ID:       uuid.New(),
		TeamID:   teamID,
		FullName: input.FullName,
		Role:     entity.TeamMemberRolePlayer,
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

	return &AddPlayerRes{
		ID:           playerRes.ID,
		TeamMemberID: playerRes.TeamMemberID,
		FullName:     playerRes.FullName,
		JerseyNumber: playerRes.JerseyNumber,
		Position:     playerRes.Position,
		Status:       playerRes.Status,
	}, nil

}

func (pu *playerUsecase) UpdatePlayerStatus(ctx context.Context, input *UpdatPlayerStatusReq) (*UpdatePlayerStatusRes, error) {

	playerID, err := uuid.Parse(input.PlayerID)

	if err != nil {
		return nil, apperror.NewBadRequestError("invalid player id")
	}

	if err := pu.playerRepository.UpdatePlayerStatus(ctx, &playerID, &input.Status); err != nil {
		return nil, err
	}
	return &UpdatePlayerStatusRes{
		Success: true,
	}, nil
}

func (pu *playerUsecase) ListTeamPlayers(ctx context.Context, input *ListTeamPlayersReq) ([]PlayerRes, *PaginateDetails, error) {

	teamId, err := uuid.Parse(input.TeamID)
	if err != nil {
		return nil, nil, apperror.NewBadRequestError("invalid player id")

	}
	if input.Page <= 0 {
		input.Page = 1
	}

	switch {
	case input.Limit > entity.MaxDefaultPaginateLimit:
		input.Limit = entity.MaxDefaultPaginateLimit
	case input.Limit <= 0:
		input.Limit = entity.MinDefaultPagination
	}

	players, total, err := pu.playerRepository.GetTeamPlayers(ctx, &player.ListUserReq{
		TeamID:       teamId,
		Page:         input.Page,
		Limit:        input.Limit,
		PlayerStatus: input.PlayerStatus,
		Position:     input.Position,
		Search:       input.Search,
	})

	if err != nil {
		return nil, nil, err
	}
	totalPage := total / int64(input.Limit)
	if total%int64(input.Limit) > 0 {
		totalPage += 1
	}

	var playerList []PlayerRes

	for _, player := range players {
		playerList = append(playerList, PlayerRes{
			ID:           player.ID,
			TeamMemberID: player.TeamMemberID,
			FullName:     player.FullName,
			JerseyNumber: player.JerseyNumber,
			Position:     player.Position,
			Status:       player.Status,
		})
	}

	return playerList, &PaginateDetails{
		TotalPage: int32(totalPage),
		Page:      input.Page,
		Limit:     input.Limit,
		Total:     total,
	}, nil
}

func (pu *playerUsecase) GetPlayer(ctx context.Context, input *GetPlayerReq) (*GetPlayerRes, error) {

	playerID, err := uuid.Parse(input.PlayerID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid player id")

	}

	player, err := pu.playerRepository.GetPlayer(ctx, &playerID)
	if err != nil {
		return nil, err
	}

	return &GetPlayerRes{
		PlayerID:     player.ID,
		TeamMemberID: player.TeamMemberID,
		FullName:     player.FullName,
		DateOfBirth:  player.DateOfBirth,
		JerseyNumber: player.JerseyNumber,
		Position:     player.Position,
		Height:       player.Height,
		Weight:       player.Weight,
		Status:       player.Status,
		CreatedAt:    player.CreatedAt,
	}, nil
}
