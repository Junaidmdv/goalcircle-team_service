package team

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/permission"
	playerrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/player"
	staffrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staff"
	team_repo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/transaction"
	code "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/invitation"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/datetime"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/imageutil"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
)

type TeamUsecase interface {
	CreateTeam(context.Context, *CreateTeamReq) (*CreateTeamRes, error)
	DeleteTeam(ctx context.Context, teamID uuid.UUID, staffID uuid.UUID) error
	UpdateTeamDetails(context.Context, *UpdateTeamDetailsReq) (*UpdateTeamDetailsRes, error)
	ListTeams(context.Context, *ListTeamsReq) (*ListTeamsRes, error)
	ChangeCaptain(context.Context, *ChangeCaptainReq) (*ChangeCaptainRes, error)
	ChangeViceCaptain(context.Context, *ChangeViceCaptainReq) (*ChangeViceCaptainRes, error)
	UploadLogo(context.Context, *UploadLogoReq) (*UploadLogoRes, error)
	GetPresignedURL(context.Context, *GetPresignedUrlReq) (*GetPresignedUrlRes, error)
	DeleteLogo(ctx context.Context, req *RemoveLogoReq) (*RemoveLogoRes, error)
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
	datecal           datetime.DateCalculator
	transactions      transaction.UnitOfWork
	staffRepo         staffrepo.StaffRepository
}

func NewTeamUsecase(teamrepo team_repo.TeamRepository,
	logger logger.Logger,
	code code.CodeGenerater,
	tmrepo teammember.TeamMemberRepository,
	playerepo playerrepo.PlayerRepository,
	ob storage.ObjectStorage,
	obconfig *config.ObjectStorageConfig,
	dateCalc datetime.DateCalculator,
	trans transaction.UnitOfWork,
	staffrepo staffrepo.StaffRepository,
) TeamUsecase {
	return &teamUsecase{
		teamRepo:          teamrepo,
		logger:            logger,
		code:              code,
		teamMemberRepo:    tmrepo,
		playerRepo:        playerepo,
		objectStore:       ob,
		objectStoreConfig: obconfig,
		datecal:           dateCalc,
		transactions:      trans,
		staffRepo:         staffrepo,
	}
}

func (tu *teamUsecase) CreateTeam(ctx context.Context, dt *CreateTeamReq) (*CreateTeamRes, error) {

	userId, err := uuid.Parse(dt.UserID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid user id")
	}

	exist, err := tu.teamMemberRepo.HasUnreleasedMembership(ctx, userId)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, apperror.NewConflictError("user already assosite a team")
	}

	code, err := tu.code.GenerateCode("TM-")
	if err != nil {
		tu.logger.Error("Failed to generate code", "error", err, "method", "teamUsercase")
	}

	formatedTeamName := FormatTeamName(dt.TeamName)

	var staff *entity.Staff
	var teamMember *entity.TeamMember
	var team *entity.Team

	err = tu.transactions.Execute(ctx, func(r transaction.Repositories) error {
		teamRes, err := r.Team().CreateTeam(ctx, &entity.Team{
			ID:          uuid.New(),
			Name:        formatedTeamName,
			ShortName:   tu.code.GenerateShortName(dt.TeamName, 2, true),
			City:        dt.City,
			Description: dt.Description,
			TeamCode:    code,
			TeamStatus:  entity.TeamStatusActive,
		})

		if err != nil {
			return nil
		}

		teamMemberRes, err := r.TeamMember().AddTeamMember(ctx, &entity.TeamMember{
			ID:       uuid.New(),
			TeamID:   teamRes.ID,
			UserID:   userId,
			Status:   dt.TeamMemberStatus,
			JoinedAt: time.Now(),
			Role:     entity.TeamMemberStaff,
		})
		if err != nil {
			return err
		}

		staffRes, err := r.Staff().AddStaff(ctx, &entity.Staff{
			ID:           uuid.New(),
			TeamMemberID: teamMemberRes.ID,
			FullName:     dt.OwnerName,
			Role:         dt.Role,
			Designation:  dt.Designation,
		})
		if err != nil {
			return err
		}

		team = teamRes
		teamMember = teamMemberRes
		staff = staffRes

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &CreateTeamRes{
		Team: &TeamDetails{
			TeamID:    team.ID,
			Name:      team.Name,
			ShortName: team.ShortName,
			City:      team.City,
			Code:      team.TeamCode,
		},
		Owner: &TeamMemberDetails{
			TeamMemberID: teamMember.ID,
			StaffID:      staff.ID,
			FullName:     staff.FullName,
			Role:         staff.Role,
			Designation:  staff.Designation,
		},
	}, nil
}

func (tu *teamUsecase) DeleteTeam(ctx context.Context, teamId uuid.UUID, staffID uuid.UUID) error {
	return tu.transactions.Execute(ctx, func(r transaction.Repositories) error {

		if err := r.Staff().DeleteStaff(ctx, staffID); err != nil {
			return err
		}

		if err := r.Team().DeleteTeam(ctx, teamId); err != nil {
			return err
		}

		return nil
	})
}

func (tu *teamUsecase) UpdateTeamDetails(ctx context.Context, req *UpdateTeamDetailsReq) (*UpdateTeamDetailsRes, error) {

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	teamOwnerRole, err := tu.staffRepo.GetStaffDesignation(ctx, teamID, userID)
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

		if req.ShortName == nil {
			shortName := tu.code.GenerateShortName(*req.Name, 2, true)
			req.ShortName = &shortName
		}
	}

	if err := tu.teamRepo.UpdateTeamDetails(ctx, teamMember.TeamID, &entity.Team{
		ID:          teamMember.TeamID,
		Name:        *req.Name,
		ShortName:   *req.ShortName,
		City:        *req.City,
		Description: *req.Description,
		PhoneNum:    *req.PhoneNum,
		Email:       *req.Email,
	}); err != nil {
		return nil, err
	}

	team, err := tu.teamRepo.GetTeamDetails(ctx, teamMember.TeamID)
	if err != nil {
		return nil, err
	}

	return &UpdateTeamDetailsRes{
		TeamID:      team.ID.String(),
		Name:        team.Name,
		City:        team.City,
		Description: team.Description,
		ShortName:   team.ShortName,
	}, nil
}

func (tu *teamUsecase) ChangeCaptain(ctx context.Context, req *ChangeCaptainReq) (*ChangeCaptainRes, error) {

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	teamId, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewFailedPreCondition("invalid team id")
	}

	teamOwnerRole, err := tu.staffRepo.GetStaffDesignation(ctx, teamId, userID)
	if err != nil {
		return nil, err
	}

	permite := permission.HasPermission(teamOwnerRole, permission.PermissionAddStaff)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user not allowed to create staff")
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
	teamId, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewFailedPreCondition("invalid team id")
	}
	teamOwnerRole, err := tu.staffRepo.GetStaffDesignation(ctx, teamId, userID)
	if err != nil {
		return nil, err
	}

	permite := permission.HasPermission(teamOwnerRole, permission.PermissionAddStaff)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user not allowed to create staff")
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
			TeamCode:   team.TeamCode,
			TeamStatus: team.TeamStatus,
		})

		if team.LogoKey != nil && *team.LogoKey != "" {
			logoURL, err := tu.objectStore.GetPresignedURL(
				ctx,
				tu.objectStoreConfig.Bucket,
				*team.LogoKey,
				tu.objectStoreConfig.PresignedURLExpiry,
			)
			if err != nil {
				tu.logger.Error("failed to generate team logo presigned URL", "error", err)
			} else {
				team.LogoKey = &logoURL
			}
		}

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
		return nil, apperror.NewBadRequestError("invalid team id")
	}
	teamId, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewFailedPreCondition("invalid team id")
	}

	teamOwnerRole, err := tu.staffRepo.GetStaffDesignation(ctx, teamId, userID)
	if err != nil {
		return nil, err
	}

	permite := permission.HasPermission(teamOwnerRole, permission.PermissionAddStaff)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user not allowed to create staff")
	}

	team, err := tu.teamRepo.GetTeamDetails(ctx, teamId)
	if err != nil {
		return nil, err
	}

	webpBytes, err := imageutil.ConvertImageIntoWebpbFormate(req.LogoData)
	if err != nil {
		tu.logger.Error("failed update photo", "error", err)
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

	if team.LogoKey == nil {
		return nil, apperror.NewBadRequestError("logo is not updated or removed")
	}

	presignedUrl, err := tu.objectStore.GetPresignedURL(ctx, tu.objectStoreConfig.Bucket, *team.LogoKey, tu.objectStoreConfig.PresignedURLExpiry)
	if err != nil {
		return nil, err
	}

	return &GetPresignedUrlRes{
		TeamID:       team.ID.String(),
		PresignedUrl: presignedUrl,
	}, nil
}

func (tu *teamUsecase) DeleteLogo(ctx context.Context, req *RemoveLogoReq) (*RemoveLogoRes, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}
	teamId, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewFailedPreCondition("invalid team id")
	}
	teamOwnerRole, err := tu.staffRepo.GetStaffDesignation(ctx, teamId, userID)
	if err != nil {
		return nil, err
	}

	permite := permission.HasPermission(teamOwnerRole, permission.PermissionAddStaff)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user not allowed to create staff")
	}

	team, err := tu.teamRepo.GetTeamDetails(ctx, teamId)
	if err != nil {
		return nil, err
	}

	if team.LogoKey == nil {
		return nil, apperror.NewBadRequestError("logo alreay deleted or not set")
	}

	if err := tu.objectStore.Delete(ctx, tu.objectStoreConfig.Bucket, *team.LogoKey); err != nil {
		return nil, err
	}

	if err := tu.teamRepo.ClearTeamLogo(ctx, teamId); err != nil {
		return nil, err
	}

	return &RemoveLogoRes{
		TeamId: team.ID.String(),
	}, nil

}

func (th *teamUsecase) LeaveTeam(ctx context.Context, input *LeaveTeamReq) (*LeaveTeamRes, error) {

	return nil, nil
}
