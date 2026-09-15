package teamplayeruc

import (
	"context"

	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/permission"
	staffrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staff"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teamplayer"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/google/uuid"
)

type TeamPlayerUsecase interface {
	ListTeamPlayers(context.Context, *ListTeamPlayersReq) ([]PlayerRes, *PaginateDetails, error)
	GetTeamPlayer(ctx context.Context, input *GetTeamPlayerReq) (*GetTeamPlayerRes, error)
	ReleasePlayer(context.Context, *ReleasePlayerReq) (*ReleasePlayerRes, error)  
}

type teamPlayerUsecase struct {
	teamPlayerRepo      teamplayer.TeamPlayerRepository
	objectStorage       storage.ObjectStorage
	objectStorageConfig *config.ObjectStorageConfig
	staffRepo           staffrepo.StaffRepository
	teamMemberRepo      teammember.TeamMemberRepository
}

func NewTeamPlayerUsecas() TeamPlayerUsecase {
	return &teamPlayerUsecase{}
}

func (pu *teamPlayerUsecase) ListTeamPlayers(ctx context.Context, input *ListTeamPlayersReq) ([]PlayerRes, *PaginateDetails, error) {

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

	players, total, err := pu.teamPlayerRepo.GetTeamPlayers(ctx, &teamplayer.ListPlayerDetailsReq{
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

	for i, player := range players {
		playerList = append(playerList, PlayerRes{
			ID:           player.ID,
			TeamMemberID: player.TeamMemberID,
			FullName:     player.Player.FullName,
			Position:     player.Player.Position,
			Status:       player.Player.Status,
		})

		presignedUrl, err := pu.objectStorage.GetPresignedURL(ctx, pu.objectStorageConfig.Bucket, *player.Player.ProfileImageKey, pu.objectStorageConfig.PresignedURLExpiry)
		if err == nil {
			playerList[i].ImageUrl = presignedUrl
		}

	}

	return playerList, &PaginateDetails{
		TotalPage: int32(totalPage),
		Page:      input.Page,
		Limit:     input.Limit,
		Total:     total,
	}, nil
}

func (pu *teamPlayerUsecase) GetTeamPlayer(ctx context.Context, input *GetTeamPlayerReq) (*GetTeamPlayerRes, error) {

	playerID, err := uuid.Parse(input.PlayerID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid player id")

	}
	teamID, err := uuid.Parse(input.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	player, err := pu.teamPlayerRepo.GetTeamPlayer(ctx, teamID, playerID)
	if err != nil {
		return nil, err
	}

	return &GetTeamPlayerRes{
		TeamID:      player.TeamMemberID.String(),
		PlayerID:    player.ID,
		FullName:    player.Player.FullName,
		DateOfBirth: player.Player.DateOfBirth,
		Position:    player.Player.Position,
		Status:      player.Player.Status,
		CreatedAt:   player.CreatedAt,
	}, nil
}

func (pu *teamPlayerUsecase) ReleasePlayer(ctx context.Context, input *ReleasePlayerReq) (*ReleasePlayerRes, error) {
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}
	teamID, err := uuid.Parse(input.TeamID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid team id")
	}

	teamOwnerRole, err := pu.staffRepo.GetStaffDesignation(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}

	permite := permission.HasPermission(teamOwnerRole, permission.PermissionReleasePlayer)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user not allowed to create staff")
	}

	playerID, err := uuid.Parse(input.PlayerID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid player id")
	}

	if err := pu.teamPlayerRepo.ReleaseTeamPlayer(ctx, teamID, playerID); err != nil {
		return nil, err
	}

	playerDetails, err := pu.teamPlayerRepo.GetTeamPlayer(ctx, teamID, playerID)
	if err != nil {
		return nil, err
	}

	if err := pu.teamMemberRepo.ReleaseMember(ctx, teamID, playerDetails.TeamMemberID); err != nil {
		return nil, err
	}

	return &ReleasePlayerRes{
		ID:           playerDetails.ID,
		TeamMemberID: playerDetails.TeamMemberID,
		FullName:     playerDetails.Player.FullName,
		Position:     playerDetails.Player.Position,
		Status:       playerDetails.Player.Status,
		ReleasedAt:   playerDetails.LeftAt,
		JoinedAt:     playerDetails.JoinedAt,
	}, nil
}
