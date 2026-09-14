package playeruc

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/player"
	staffrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staff"
	teamRepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teamplayer"
	code "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/invitation"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/imageutil"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
)

type PlayerUsecase interface {
	AddNewPlayer(context.Context, *AddPlayerReq) (*AddPlayerRes, error)
	UpdatePlayerDetails(context.Context, *UpdatePlayerReq) (*UpdatePlayersRes, error)
	ListPlayer(ctx context.Context, input *ListPlayerReq) (*ListPlayerRes, error)
	GetPlayer(context.Context, *GetPlayerReq) (*GetPlayerRes, error)
	UpdateImage(context.Context, *UpdatePlayerImageReq) (*UpdatePlayerImageRes, error)
	GetPlayerPresignedUrl(context.Context, *GetPlayerPresignedUrlReq) (*GetPlayerPresignedUrlRes, error)
	RemovePlayerImage(context.Context, *RemovePlayerImageReq) (*RemovePlayerImageRes, error)
}

type playerUsecase struct {
	playerRepository    player.PlayerRepository
	teamMemberRepo      teammember.TeamMemberRepository
	teamRepository      teamRepo.TeamRepository
	teamPlayerRepo      teamplayer.TeamPlayerRepository
	staffRepo           staffrepo.StaffRepository
	code                code.CodeGenerater
	logger              logger.Logger
	ObjectStorage       storage.ObjectStorage
	ObjectStorageConfig *config.ObjectStorageConfig
}

func NewPlayerUsecase(player player.PlayerRepository,
	tmRepo teammember.TeamMemberRepository,
	teamRepo teamRepo.TeamRepository,
	staffRepo staffrepo.StaffRepository,
	teamPlayerRepo teamplayer.TeamPlayerRepository,
	logger logger.Logger,
	objectStorage storage.ObjectStorage,
	objectStorageconfig *config.ObjectStorageConfig,
	code code.CodeGenerater) PlayerUsecase {

	return &playerUsecase{
		playerRepository:    player,
		teamMemberRepo:      tmRepo,
		teamRepository:      teamRepo,
		teamPlayerRepo:      teamPlayerRepo,
		staffRepo:           staffRepo,
		code:                code,
		ObjectStorage:       objectStorage,
		ObjectStorageConfig: objectStorageconfig,
		logger:              logger,
	}
}

func (pu *playerUsecase) AddNewPlayer(ctx context.Context, input *AddPlayerReq) (*AddPlayerRes, error) {

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid player id")
	}

	exist, err := pu.playerRepository.IsPlayerExist(ctx, userID)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, apperror.NewConflictError("user already registered as player")
	}

	playerRes, err := pu.playerRepository.CreatePlayer(ctx, &entity.Player{
		ID:            uuid.New(),
		UserID:        userID,
		FullName:      input.FullName,
		DateOfBirth:   input.DOB,
		Nationality:   input.Nationality,
		Position:      input.Position,
		PreferredFoot: input.PreferredFoot,
	})

	if err != nil {
		return nil, err
	}

	webpBytes, err := imageutil.ConvertImageIntoWebpbFormate(input.ImageBytes)
	if err != nil {
		return nil, err
	}

	objectName := fmt.Sprintf("/players/%s/logo.webp", playerRes.ID)

	key, err := pu.ObjectStorage.Upload(ctx, pu.ObjectStorageConfig.Bucket, objectName, bytes.NewReader(webpBytes), int64(len(webpBytes)), "image/webp")

	if err != nil {
		return nil, err
	}

	if err := pu.playerRepository.UpdateImageKey(ctx, playerRes.ID, key); err != nil {
		return nil, err
	}

	presignedUrl, err := pu.ObjectStorage.GetPresignedURL(ctx, pu.ObjectStorageConfig.Bucket, objectName, pu.ObjectStorageConfig.PresignedURLExpiry)
	if err != nil {
		return nil, err
	}

	return &AddPlayerRes{
		ID:            playerRes.ID,
		FullName:      playerRes.FullName,
		DOB:           playerRes.DateOfBirth,
		Nationality:   playerRes.Nationality,
		PreferredFoot: playerRes.PreferredFoot,
		Position:      playerRes.Position,
		Status:        playerRes.Status,
		PresignedUrl:  presignedUrl,
	}, nil

}

func (pu *playerUsecase) UpdatePlayerDetails(ctx context.Context, input *UpdatePlayerReq) (*UpdatePlayersRes, error) {

	playerID, err := uuid.Parse(input.PlayerID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid player id")
	}

	// userID, err := uuid.Parse(input.UserID)
	// if err != nil {
	// 	return nil, apperror.NewBadRequestError("invalid user id")
	// }

	if err := pu.playerRepository.UpdatePlayerDetails(ctx, playerID, &entity.Player{
		FullName:      *input.FullName,
		DateOfBirth:   *input.DOB,
		Position:      *input.Position,
		PreferredFoot: *input.PreferredFoot,
	}); err != nil {

		return nil, err
	}

	player, err := pu.playerRepository.GetPlayer(ctx, playerID)
	if err != nil {
		return nil, err
	}

	return &UpdatePlayersRes{
		PlayerID:      playerID,
		FullName:      player.FullName,
		Status:        player.Status,
		DOB:           player.DateOfBirth,
		Position:      player.Position,
		Nationality:   player.Nationality,
		PreferredFoot: player.PreferredFoot,
	}, nil
}

func (pu *playerUsecase) ListPlayer(ctx context.Context, input *ListPlayerReq) (*ListPlayerRes, error) {

	if input.Page <= 0 {
		input.Page = 1
	}

	switch {
	case input.Limit > entity.MaxDefaultPaginateLimit:
		input.Limit = entity.MaxDefaultPaginateLimit
	case input.Limit <= 0:
		input.Limit = entity.MinDefaultPagination
	}

	players, total, err := pu.playerRepository.GetPlayers(ctx, &player.ListUserReq{
		Page:         input.Page,
		Limit:        input.Limit,
		PlayerStatus: input.PlayerStatus,
		Position:     input.Position,
		Search:       input.Search,
	})

	if err != nil {
		return nil, err
	}

	var playerList []ListPlayerDetails

	for i, p := range players {
		playerList = append(playerList, ListPlayerDetails{
			ID:           p.ID,
			FullName:     p.FullName,
			Position:     p.Position,
			PreferedFoot: p.PreferredFoot,
			Status:       p.Status,
		})
		presignedUrl, err := pu.ObjectStorage.GetPresignedURL(ctx, pu.ObjectStorageConfig.Bucket, *p.ProfileImageKey, pu.ObjectStorageConfig.PresignedURLExpiry)
		if err == nil {
			playerList[i].ImageUrl = presignedUrl
		}

	}
	totalPage := total / int64(input.Limit)
	if total%int64(input.Limit) > 0 {
		totalPage += 1
	}

	return &ListPlayerRes{
		Players: playerList,
		Pagination: PaginateDetails{
			TotalPage: int32(totalPage),
			Page:      input.Page,
			Limit:     input.Limit,
			Total:     total,
		},
	}, nil

}

func (pu *playerUsecase) GetPlayer(ctx context.Context, input *GetPlayerReq) (*GetPlayerRes, error) {

	playerID, err := uuid.Parse(input.PlayerID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid player id")

	}

	player, err := pu.playerRepository.GetPlayer(ctx, playerID)
	if err != nil {
		return nil, err
	}

	return &GetPlayerRes{
		PlayerID:    player.ID,
		FullName:    player.FullName,
		DateOfBirth: player.DateOfBirth,
		Position:    player.Position,
		Status:      player.Status,
		CreatedAt:   player.CreatedAt,
	}, nil
}

func (pu *playerUsecase) UpdateImage(ctx context.Context, input *UpdatePlayerImageReq) (*UpdatePlayerImageRes, error) {

	playerID, err := uuid.Parse(input.PlayerID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid team id")
	}

	player, err := pu.playerRepository.GetPlayer(ctx, playerID)
	if err != nil {
		return nil, err
	}

	webpBytes, err := imageutil.ConvertImageIntoWebpbFormate(input.ImageData)
	if err != nil {
		return nil, err
	}

	Reader := bytes.NewReader(webpBytes)
	if player.ProfileImageKey == nil {
		objectName := fmt.Sprintf("/player/%s/logo.webp", player.ID)
		player.ProfileImageKey = &objectName
	}

	key, err := pu.ObjectStorage.Upload(ctx, pu.ObjectStorageConfig.Bucket, *player.ProfileImageKey, Reader, int64(len(webpBytes)), "image/web")

	if err != nil {
		return nil, err
	}

	if err = pu.playerRepository.UpdateImageKey(ctx, playerID, key); err != nil {
		return nil, err
	}

	presignedUrl, err := pu.ObjectStorage.GetPresignedURL(ctx, pu.ObjectStorageConfig.Bucket, *player.ProfileImageKey, pu.ObjectStorageConfig.PresignedURLExpiry)
	if err != nil {
		return nil, err
	}

	return &UpdatePlayerImageRes{
		PlayerID:     input.PlayerID,
		PresignedUrl: presignedUrl,
	}, nil

}

func (pu *playerUsecase) GetPlayerPresignedUrl(ctx context.Context, input *GetPlayerPresignedUrlReq) (*GetPlayerPresignedUrlRes, error) {

	playerID, err := uuid.Parse(input.PlayerID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid player id")
	}

	player, err := pu.playerRepository.GetPlayer(ctx, playerID)
	if err != nil {
		return nil, err
	}

	if player.ProfileImageKey == nil {
		return nil, apperror.NewNotFoundError("image not uploaded or removed")
	}

	presignedUrl, err := pu.ObjectStorage.GetPresignedURL(ctx, pu.ObjectStorageConfig.Bucket, *player.ProfileImageKey, pu.ObjectStorageConfig.PresignedURLExpiry)
	if err != nil {
		return nil, err
	}

	return &GetPlayerPresignedUrlRes{
		TeamID:       input.TeamID,
		PlayerId:     input.PlayerID,
		PresignedUrl: presignedUrl,
	}, nil
}

func (pu *playerUsecase) RemovePlayerImage(ctx context.Context, req *RemovePlayerImageReq) (*RemovePlayerImageRes, error) {

	// userID, err := uuid.Parse(req.UserID)
	// if err != nil {
	// 	pu.logger.Error("invalid team id", "error", errors.New("invalid user id from token"))
	// 	return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	// }

	playerID, err := uuid.Parse(req.PlayerID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid player id")
	}

	key, err := pu.playerRepository.GetPlayerImageKey(ctx, playerID)
	if err != nil {
		return nil, err
	}
	err = pu.ObjectStorage.Delete(ctx, pu.ObjectStorageConfig.Bucket, key)
	if err != nil {
		return nil, err
	}

	return &RemovePlayerImageRes{
		Success: true,
	}, nil
}
