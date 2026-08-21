package team

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/permission"
	playerrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/player"
	team_repo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	code "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/invitation"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/imageutil"
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
	GetPresignedURL(context.Context, *GetPresignedUrlReq) (*GetPresignedUrlRes, error)
	GetTeam(context.Context, *GetTeamReq) (*GetTeamRes, error)
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

func NewTeamUsecase(teamrepo team_repo.TeamRepository,
	logger logger.Logger,
	code code.CodeGenerater,
	tmrepo teammember.TeamMemberRepository,
	playerepo playerrepo.PlayerRepository,
	ob storage.ObjectStorage,
	obconfig *config.ObjectStorageConfig) TeamUsecase {
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

	userId, err := uuid.Parse(dt.UserID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid user id")
	}
	tu.teamMemberRepo.GetActiveTeamMemberByUserID(ctx, userId)

	code, err := tu.code.GenerateCode("TM")
	if err != nil {
		tu.logger.Error("Failed to generate code", "error", err, "method", "teamUsercase")
	}

	formatedTeamName := FormatTeamName(dt.Name)

	res, err := tu.teamRepo.CreateTeam(ctx, &entity.Team{
		ID:          uuid.New(),
		Name:        formatedTeamName,
		ShortName:   tu.code.GenerateShortName(dt.Name, 2, true),
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

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	teamOwnerRole, err := tu.teamMemberRepo.GetStaffDesignation(ctx, userID)
	if err != nil {
		return nil, err
	}

	permite := permission.HasPermission(teamOwnerRole, permission.PermissionAddStaff)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user not allowed to create staff")
	}

	teamMember, err := tu.teamMemberRepo.GetActiveTeamMemberByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		*req.Name = FormatTeamName(*req.Name)
	}

	if req.Name != nil && req.ShortName == nil {
		*req.ShortName = tu.code.GenerateShortName(*req.Name, 2, true)
	}

	team, err := tu.teamRepo.UpdateTeamDetails(ctx, teamMember.TeamID, &team_repo.UpdateTeamReq{
		TeamID:      teamMember.TeamID,
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
		LogoKey:     team.LogoKey,
	}, nil
}

func (tu *teamUsecase) ChangeCaptain(ctx context.Context, req *ChangeCaptainReq) (*ChangeCaptainRes, error) {

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	teamOwnerRole, err := tu.teamMemberRepo.GetStaffDesignation(ctx, userID)
	if err != nil {
		return nil, err
	}

	permite := permission.HasPermission(teamOwnerRole, permission.PermissionAddStaff)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user not allowed to create staff")
	}

	teamId, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewFailedPreCondition("invalid team id")
	}

	playerId, err := uuid.Parse(req.PlayerID)
	if err != nil {
		return nil, apperror.NewFailedPreCondition("invalid player id")
	}

	player, err := tu.playerRepo.GetPlayer(ctx, teamId, playerId)
	if err != nil {
		return nil, err
	}

	if player.TeamMember.Status != entity.TeamMemberStatusActive {
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

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}
	teamOwnerRole, err := tu.teamMemberRepo.GetStaffDesignation(ctx, userID)
	if err != nil {
		return nil, err
	}

	permite := permission.HasPermission(teamOwnerRole, permission.PermissionAddStaff)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user not allowed to create staff")
	}

	teamId, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewFailedPreCondition("invalid team id")
	}

	playerId, err := uuid.Parse(req.PlayerID)
	if err != nil {
		return nil, apperror.NewFailedPreCondition("invalid player id")
	}

	player, err := tu.playerRepo.GetPlayer(ctx, teamId, playerId)
	if err != nil {
		return nil, err
	}

	if player.TeamMember.Status != entity.TeamMemberStatusActive {
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
	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	team, err := tu.teamRepo.GetTeamDetails(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return &GetTeamRes{
		ID:            team.ID.String(),
		Name:          team.Name,
		ShortName:     team.ShortName,
		City:          team.City,
		LogoKey:       team.LogoKey,
		TeamCode:      team.TeamCode,
		Description:   team.Description,
		Email:         team.Email,
		PhoneNum:      team.PhoneNum,
		TeamStatus:    string(team.TeamStatus),
		PlayerCount:   team.PlayerCount,
		CaptainID:     team.CaptainID.String(),
		ViceCaptainID: team.ViceCaptainID.String(),
		CreatedAt:     team.CreatedAt,
	}, nil
}

func (tu *teamUsecase) UploadLogo(ctx context.Context, req *UploadLogoReq) (*UploadLogoRes, error) {

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid team id")
	}

	teamMember, err := tu.teamMemberRepo.GetActiveTeamMemberByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	team, err := tu.teamRepo.GetTeamDetails(ctx, teamMember.TeamID)
	if err != nil {
		return nil, err
	}

	webpBytes, err := imageutil.ConvertImageIntoWebpbFormate(req.LogoData)
	if err != nil {
		return nil, err
	}

	objectName := fmt.Sprintf("/team/%s/logo.webp", team.ID)

	key, err := tu.objectStore.Upload(ctx, tu.objectStoreConfig.Bucket, objectName, bytes.NewReader(webpBytes), int64(len(webpBytes)), "image/webp")

	if err != nil {
		return nil, err
	}

	if err = tu.teamRepo.UpdateLogoKey(ctx, team.ID, key); err != nil {
		return nil, err
	}

	presignedUrl, err := tu.objectStore.GetPresignedURL(ctx, tu.objectStoreConfig.Bucket, key, tu.objectStoreConfig.PresignedURLExpiry)

	if err != nil {
		return nil, err
	}

	return &UploadLogoRes{
		TeamID:       team.ID.String(),
		PresignedUrl: presignedUrl,
	}, nil
}

func (tu *teamUsecase) GetPresignedURL(ctx context.Context, input *GetPresignedUrlReq) (*GetPresignedUrlRes, error) {

	teamID, err := uuid.Parse(input.TeamID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid team id.")
	}

	team, err := tu.teamRepo.GetTeamDetails(ctx, teamID)
	if err != nil {
		return nil, err
	}

	presignedUrl, err := tu.objectStore.GetPresignedURL(ctx, tu.objectStoreConfig.Bucket, team.LogoKey, tu.objectStoreConfig.PresignedURLExpiry)
	if err != nil {
		return nil, err
	}

	return &GetPresignedUrlRes{
		TeamID:       team.ID.String(),
		PresignedUrl: presignedUrl,
	}, nil
}
