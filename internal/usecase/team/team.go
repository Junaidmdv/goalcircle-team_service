package team

import (
	"bytes"
	"context"
	"image"
	"net/http"

	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/permission"
	playerrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/player"
	team_repo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	code "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/invitation"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
)

type TeamUsecase interface {
	CreateTeam(context.Context, *CreateTeamReq) (*CreateTeamRes, error)
	DeleteTeam(context.Context, uuid.UUID) error
	UpdateTeamDetails(context.Context, *UpdateTeamDetailsReq) (*UpdateTeamDetailsRes, error)
	ListTeams(context.Context, *ListTeamsReq) (*ListTeamsRes, error)
	ChangeCaptain(context.Context, *ChangeCaptainReq) (*ChangeCaptainRes, error)
	ChangeViceCaptain(context.Context, *ChangeViceCaptainReq) (*ChangeViceCaptainRes, error)
	UploadLogo(context.Context, *UploadLogoReq) (*UploadLogoRes, error)
}

type teamUsecase struct {
	teamRepo          team_repo.TeamRepository
	teamMemberRepo    teammember.TeamMemberRepository
	logger            logger.Logger
	code              code.CodeGenerater
	playerRepo        playerrepo.PlayerRepository
	objectStore       storage.ObjectStorage
	objectStoreConfig *config.ObjectStorageConfig
}

func NewTeamUsecase(teamrepo team_repo.TeamRepository, logger logger.Logger, code code.CodeGenerater, tmrepo teammember.TeamMemberRepository, playerepo playerrepo.PlayerRepository, ob storage.ObjectStorage, obconfig *config.ObjectStorageConfig) TeamUsecase {
	return &teamUsecase{
		teamRepo:          teamrepo,
		logger:            logger,
		code:              code,
		teamMemberRepo:    tmrepo,
		playerRepo:        playerepo,
		objectStore:       ob,
		objectStoreConfig: obconfig,
	}
}

func (tu *teamUsecase) CreateTeam(ctx context.Context, dt *CreateTeamReq) (*CreateTeamRes, error) {

	code, err := tu.code.GenerateCode("TM")
	if err != nil {
		tu.logger.Error("Failed to generate code", "error", err, "method", "teamUsercase")
	}

	formatedTeamName := FormatTeamName(dt.Name)

	res, err := tu.teamRepo.CreateTeam(ctx, &entity.Team{
		ID:          uuid.New(),
		Name:        formatedTeamName,
		ShortName:   tu.code.GenerateShortName(dt.Name),
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

func (tu *teamUsecase) UpdateTeamDetails(ctx context.Context, req *UpdateTeamDetailsReq) (*UpdateTeamDetailsRes, error) {

	teamId, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewFailedPreCondition("invalid team id")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewFailedPreCondition("invalid team member id")
	}

	role, err := tu.teamMemberRepo.GetTeamMemeberRole(ctx, teamId, userID)
	if err != nil {
		return nil, err
	}

	permitted := permission.HasPermissionTeam(role, permission.PermissionUpdateTeamDetails)
	if !permitted {
		return nil, apperror.NewPermissionDenied("user not allowed to update team details")
	}

	if req.Name != nil {
		*req.Name = FormatTeamName(*req.Name)
	}

	if req.Name != nil && req.ShortName == nil {
		*req.ShortName = tu.code.GenerateShortName(*req.Name)
	}

	team, err := tu.teamRepo.UpdateTeamDetails(ctx, teamId, &team_repo.UpdateTeamReq{
		TeamID:      teamId,
		Name:        req.Name,
		ShortName:   req.ShortName,
		City:        req.City,
		Description: req.Description,
		PhoneNum:    req.PhoneNum,
		Email:       req.Email,
	})

	return &UpdateTeamDetailsRes{
		TeamID:      team.ID.String(),
		Name:        team.Name,
		City:        team.City,
		Description: team.Description,
		ShortName:   team.ShortName,
	}, nil
}

func (tu *teamUsecase) ChangeCaptain(ctx context.Context, req *ChangeCaptainReq) (*ChangeCaptainRes, error) {

	teamId, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewFailedPreCondition("invalid team id")
	}

	playerId, err := uuid.Parse(req.PlayerID)
	if err != nil {
		return nil, apperror.NewFailedPreCondition("invalid player id")
	}
	desig, err := tu.teamMemberRepo.GetStaffDesignation(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	permitted := permission.HasPermissionSquad(desig, permission.PermissionApointCaptain)
	if !permitted {
		return nil, apperror.NewPermissionDenied("user not allowed to set captain")
	}

	status, err := tu.playerRepo.GetPlayerStatus(ctx, teamId, playerId)
	if err != nil {
		return nil, err
	}

	if status != entity.PlayerStatusActive {
		return nil, apperror.NewFailedPreCondition("only active players can be selected as captain")
	}

	if err := tu.teamRepo.SetCaptain(ctx, teamId, playerId); err != nil {
		return nil, err
	}

	return &ChangeCaptainRes{
		PlayerID: playerId,
	}, nil
}

func (tu *teamUsecase) ChangeViceCaptain(ctx context.Context, req *ChangeViceCaptainReq) (*ChangeViceCaptainRes, error) {
	teamId, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewFailedPreCondition("invalid team id")
	}

	playerId, err := uuid.Parse(req.PlayerID)
	if err != nil {
		return nil, apperror.NewFailedPreCondition("invalid player id")
	}
	desig, err := tu.teamMemberRepo.GetStaffDesignation(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	permitted := permission.HasPermissionSquad(desig, permission.PermissionApointCaptain)
	if !permitted {
		return nil, apperror.NewPermissionDenied("user not allowed to set captain")
	}

	status, err := tu.playerRepo.GetPlayerStatus(ctx, teamId, playerId)
	if err != nil {
		return nil, err
	}

	if status != entity.PlayerStatusActive {
		return nil, apperror.NewFailedPreCondition("only active players can be selected as captain")
	}

	if err := tu.teamRepo.SetViceCaptain(ctx, teamId, playerId); err != nil {
		return nil, err
	}

	return &ChangeViceCaptainRes{
		PlayerID: playerId,
	}, nil
}

func (tu *teamUsecase) ListTeams(ctx context.Context, input *ListTeamsReq) (*ListTeamsRes, error) {

	switch {
	case input.Limit > entity.MaxDefaultPaginateLimit:
		input.Limit = entity.MaxDefaultPaginateLimit
	case input.Limit <= 0:
		input.Limit = entity.MinDefaultPagination
	}

	teams, total, err := tu.teamRepo.ListTeams(ctx, &team_repo.ListTeamsReq{
		Page:       input.Page,
		Limit:      input.Limit,
		City:       input.City,
		TeamStatus: input.Status,
		Search:     input.Search,
	})

	if err != nil {
		return nil, err
	}

	var teamList []TeamData

	for _, team := range teams {
		teamList = append(teamList, TeamData{
			TeamID:     team.ID,
			Name:       team.Name,
			City:       team.City,
			LogoUrl:    team.LogoKey,
			TeamCode:   team.TeamCode,
			TeamStatus: team.TeamStatus,
		})
	}
	totalPage := total / input.Limit
	if total%input.Limit > 0 {
		totalPage += 1
	}

	pagination := &PaginateDetails{
		TotalPage: totalPage,
		Page:      input.Page,
		Limit:     input.Limit,
		TotalItem: int(total),
	}

	return &ListTeamsRes{
		Pagination: pagination,
		Teams:      teamList,
	}, nil
}

func (tu *teamUsecase) GetTeam(ctx context.Context, req *GetTeamReq) (*GetTeamRes, error) {

	return &GetTeamRes{}, nil
}

func (tu *teamUsecase) UploadLogo(ctx context.Context, req *UploadLogoReq) (*UploadLogoRes, error) {

	teamid, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid team id")
	}

	exist, err := tu.teamRepo.IsTeamExist(ctx, teamid)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, apperror.NewNotFoundError("team is not found in this id")
	}

	contentType := http.DetectContentType(req.LogoData)

	if err = ImageAllowedFormate(contentType); err != nil {
		return nil, err
	}

	logoReader := bytes.NewReader(req.LogoData)

	config, _, err := image.DecodeConfig(logoReader)
	if err != nil {
		tu.logger.Error("failed decode image using image.DecodeConfig", "error", err)
		return nil, apperror.NewInternalError("failed decode image", err)
	}

	if err := ValidateImageDiamension(config.Height, config.Width); err != nil {
		return nil, err
	}

	objectName := CreateObjectName(req.TeamID)

	key, err := tu.objectStore.Upload(ctx, req.TeamID, tu.objectStoreConfig.Bucket, objectName, logoReader, req.Size, contentType)

	if err != nil {

		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	if err = tu.teamRepo.UpdateLogoKey(ctx, teamid, key); err != nil {
		return nil, err
	}

	presignedUrl, err := tu.objectStore.GetPresignedURL(ctx, tu.objectStoreConfig.Bucket, key, tu.objectStoreConfig.PresignedURLExpiry)

	if err != nil {
		return nil, err
	}

	return &UploadLogoRes{
		TeamID: req.TeamID, 
		PresignedUrl: presignedUrl,
	}, nil
}


