package player

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/permission"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/player"
	teamRepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	code "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/invitation"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/imageutil"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
)

type PlayerUsecase interface {
	AddNewPlayer(context.Context, *AddPlayerReq) (*AddPlayerRes, error)
	UpdatePlayerStatus(context.Context, *UpdatPlayerStatusReq) (*UpdatePlayerStatusRes, error)
	ListTeamPlayers(context.Context, *ListTeamPlayersReq) ([]PlayerRes, *PaginateDetails, error)
	GetPlayer(context.Context, *GetPlayerReq) (*GetPlayerRes, error)
	UpdateImage(context.Context, *UpdatePlayerImageReq) (*UpdatePlayerImageRes, error)
	GetPlayerPresignedUrl(context.Context, *GetPlayerPresignedUrlReq) (*GetPlayerPresignedUrlRes, error)
	ReleasePlayer(context.Context, *ReleasePlayerReq) (*ReleasePlayerRes, error)
}

type playerUsecase struct {
	playerRepository    player.PlayerRepository
	teamMemberRepo      teammember.TeamMemberRepository
	teamRepository      teamRepo.TeamRepository
	code                code.CodeGenerater
	logger              logger.Logger
	ObjectStorage       storage.ObjectStorage
	ObjectStorageConfig *config.ObjectStorageConfig
}

func NewPlayerUsecase(player player.PlayerRepository,
	tmRepo teammember.TeamMemberRepository,
	teamRepo teamRepo.TeamRepository,
	logger logger.Logger,
	objectStorage storage.ObjectStorage,
	objectStorageconfig *config.ObjectStorageConfig,
	code code.CodeGenerater) PlayerUsecase {

	return &playerUsecase{
		playerRepository:    player,
		teamMemberRepo:      tmRepo,
		teamRepository:      teamRepo,
		code:                code,
		ObjectStorage:       objectStorage,
		ObjectStorageConfig: objectStorageconfig,
		logger:              logger,
	}
}

func (pu *playerUsecase) AddNewPlayer(ctx context.Context, input *AddPlayerReq) (*AddPlayerRes, error) {

	teamID, err := uuid.Parse(input.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	role, err := pu.teamMemberRepo.GetTeamMemeberRole(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}
	permit := permission.HasPermissionTeam(role, permission.PermissionAddPlayer)
	if !permit {
		return nil, apperror.NewFailedPreCondition("user are not authorized to add players")
	}

	jerseyNumTaken, err := pu.teamRepository.IsJerseyNumOccupied(ctx, teamID, input.JerseyNumber)
	if jerseyNumTaken {
		return nil, apperror.NewBadRequestError("jerey number already taken by the active player.")
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
		Status:       entity.PlayerStatusActive,
	})

	if err != nil {
		return nil, err
	}

	webpBytes, err := imageutil.ConvertImageIntoWebpbFormate(input.ImageBytes)
	if err != nil {
		return nil, err
	}

	objectName := fmt.Sprintf("/players/%s/logo.webp", playerRes.ID)

	key, err := pu.ObjectStorage.Upload(ctx, pu.ObjectStorageConfig.Bucket, objectName, bytes.NewReader(webpBytes), int64(len(webpBytes)), input.ContentType)

	if err != nil {
		return nil, err
	}

	if err := pu.playerRepository.UpdateImageKey(ctx, playerRes.ID, key); err != nil {
		return nil, err
	}

	presignedUrl, err := pu.ObjectStorage.GetPresignedURL(ctx, pu.ObjectStorageConfig.Bucket, objectName, pu.ObjectStorageConfig.PresignedURLExpiry)

	return &AddPlayerRes{
		ID:           playerRes.ID,
		TeamMemberID: playerRes.TeamMemberID,
		FullName:     playerRes.FullName,
		JerseyNumber: playerRes.JerseyNumber,
		Position:     playerRes.Position,
		Status:       playerRes.Status,
		PresignedUrl: presignedUrl,
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
		PlayerID: input.PlayerID,
		Status:   input.Status,
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

func (pu *playerUsecase) ReleasePlayer(ctx context.Context, input *ReleasePlayerReq) (*ReleasePlayerRes, error) {

	teamID, err := uuid.Parse(input.TeamID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid team id")
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid user id")
	}

	if err := pu.teamMemberRepo.DeleteTeamMember(ctx, teamID, userID); err != nil {
		return nil, err
	}

	playerID, err := uuid.Parse(input.PlayerID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid player id")
	}
	if err := pu.playerRepository.ReleasePlayer(ctx, playerID); err != nil {
		return nil, err
	}

	return &ReleasePlayerRes{
		Success: true,
	}, nil
}

func (pu *playerUsecase) UpdateImage(ctx context.Context, input *UpdatePlayerImageReq) (*UpdatePlayerImageRes, error) {

	teamID, err := uuid.Parse(input.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid user id")
	}

	role, err := pu.teamMemberRepo.GetTeamMemeberRole(ctx, teamID, userID)

	permite := permission.HasPermissionTeam(role, permission.PermissionAddTeamPlayerImage)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user is not allowed to update player image")
	}

	playerID, err := uuid.Parse(input.PlayerID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid team id")
	}

	player, err := pu.playerRepository.GetPlayer(ctx, &playerID)
	if err != nil {
		return nil, err
	}

	webpBytes, err := imageutil.ConvertImageIntoWebpbFormate(input.ImageData)
	if err != nil {
		return nil, err
	}

	Reader := bytes.NewReader(webpBytes)

	key, err := pu.ObjectStorage.Upload(ctx, pu.ObjectStorageConfig.Bucket, player.ImageKey, Reader, int64(len(webpBytes)), input.ContentType)

	if err != nil {
		return nil, err
	}

	if err = pu.playerRepository.UpdateImageKey(ctx, playerID, key); err != nil {
		return nil, err
	}

	presignedUrl, err := pu.ObjectStorage.GetPresignedURL(ctx, pu.ObjectStorageConfig.Bucket, player.ImageKey, pu.ObjectStorageConfig.PresignedURLExpiry)
	if err != nil {
		return nil, err
	}

	return &UpdatePlayerImageRes{
		TeamID:       input.TeamID,
		PlayerID:     input.PlayerID,
		PresignedUrl: presignedUrl,
	}, nil

}

func (pu *playerUsecase) GetPlayerPresignedUrl(ctx context.Context, input *GetPlayerPresignedUrlReq) (*GetPlayerPresignedUrlRes, error) {
	teamID, err := uuid.Parse(input.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid user id")
	}

	playerID, err := uuid.Parse(input.PlayerID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid player id")
	}

	exist, err := pu.teamMemberRepo.IsTeamMemberExist(ctx, teamID, userID)

	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, apperror.NewNotFoundError("user is not a member of the team ")
	}
	player, err := pu.playerRepository.GetPlayer(ctx, &playerID)
	if err != nil {
		return nil, err
	}

	presignedUrl, err := pu.ObjectStorage.GetPresignedURL(ctx, pu.ObjectStorageConfig.Bucket, player.ImageKey, pu.ObjectStorageConfig.PresignedURLExpiry)
	if err != nil {
		return nil, err
	}

	return &GetPlayerPresignedUrlRes{
		TeamID:       input.TeamID,
		PlayerId:     input.PlayerID,
		PresignedUrl: presignedUrl,
	}, nil
}
