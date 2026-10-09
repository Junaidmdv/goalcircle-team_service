package staffuc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/permission"
	staffrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staff"
	staff_invite_repo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staffinvite"
	teamrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
	teamplayerrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teamplayer"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/invitation"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/datetime"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/imageutil"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
)

type StaffUsecase interface {
	AddStaff(context.Context, *AddStaffReq) (*AddStaffRes, error)
	JoinStaff(ctx context.Context, req *JoinStaffReq) (*JoinStaffRes, error)
	CompensateJoinStaff(ctx context.Context, req ComponsateJoinStaffReq) error
	GetStaff(context.Context, *GetStaffReq) (*GetStaffRes, error)
	UploadStaffImage(context.Context, *UploadStaffImageReq) (*UploadStaffImageRes, error)
	GetStaffImagePresignedUrl(context.Context, *StaffImagePresignedUrlReq) (*StaffImagePresignedUrlRes, error)
	RemoveStaffImage(context.Context, *RemoveStaffImageReq) (*RemoveStaffImageRes, error)
	UpdateStaff(context.Context, *UpdateStaffReq) (*UpdateStaffRes, error)
	ListTeamStaff(context.Context, *ListTeamStaffReq) (*ListTeamStaffRes, error)
	ReleaseStaff(context.Context, *ReleaseStaffReq) (*ReleaseStaffRes, error)
	TransferOwnership(context.Context, *TransferOwnershipReq) (*TransferOwnershipRes, error)
	GetStaffProfile(ctx context.Context, req *GetStaffProfileReq) (*GetStaffProfileRes, error)
}

type staffUsecase struct {
	teamMemberRepo  teammember.TeamMemberRepository
	staffInviteRepo staff_invite_repo.StaffInviteRepository
	staffRepo       staffrepo.StaffRepository
	teamRepo        teamrepo.TeamRepository
	teamPlayerRepo  teamplayerrepo.TeamPlayerRepository
	date            datetime.DateCalculator
	objstore        storage.ObjectStorage
	objcfg          *config.ObjectStorageConfig
	logger          logger.Logger
	codeGenerater   invitation.CodeGenerater
}

func NewTeamStaffUsecase(tm teammember.TeamMemberRepository,
	sr staffrepo.StaffRepository,
	tmrepo teamrepo.TeamRepository,
	staffInviteRepo staff_invite_repo.StaffInviteRepository,
	teamplayerRepo teamplayerrepo.TeamPlayerRepository,
	date datetime.DateCalculator,
	objst storage.ObjectStorage,
	cnfg *config.ObjectStorageConfig,
	code invitation.CodeGenerater,
	logger logger.Logger) StaffUsecase {
	return &staffUsecase{
		teamMemberRepo:  tm,
		staffRepo:       sr,
		staffInviteRepo: staffInviteRepo,
		teamRepo:        tmrepo,
		teamPlayerRepo:  teamplayerRepo,
		date:            date,
		codeGenerater:   code,
		objstore:        objst,
		objcfg:          cnfg,
		logger:          logger,
	}
}

func (su *staffUsecase) AddStaff(ctx context.Context, req *AddStaffReq) (*AddStaffRes, error) {

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		su.logger.Error("invalid team id", "error", errors.New("invalid user id from token"))
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	teamOwnerRole, err := su.staffRepo.GetStaffDesignation(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}
	permite := permission.HasPermission(teamOwnerRole, permission.PermissionAddStaff)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user not allowed to create staff")
	}

	if !IsValidStaffRoleDesignation(req.Role, req.Designation) {
		return nil, apperror.NewBadRequestError(
			"invalid designation for staff role",
		)
	}

	count, err := su.staffRepo.CountByDesination(ctx, teamID, req.Designation)
	if err != nil {
		return nil, err
	}

	if err := StaffMaxLimit(req.Designation, count); err != nil {
		return nil, err
	}

	code, err := su.codeGenerater.GenerateCode(req.TeamID)
	if err != nil {
		su.logger.Error("failed generate code", "error", err)
		return nil, apperror.NewInternalError("failed to generate code", err)
	}

	staff, err := su.staffInviteRepo.AddStaffDetails(ctx, &entity.StaffInvite{
		ID:          uuid.New(),
		TeamID:      teamID,
		FullName:    req.FullName,
		Role:        req.Role,
		Designation: req.Designation,
		PhoneNume:   req.PhoneNum,
		InviteCode:  code,
		ExpiresAt:   time.Now().Add(time.Hour * 24 * entity.InvitationTime),
	})
	if err != nil {
		return nil, err
	}

	return &AddStaffRes{
		Code:        code,
		FullName:    req.FullName,
		Role:        req.Role,
		Designation: req.Designation,
		Expiredat:   staff.ExpiresAt,
	}, nil
}

func (su *staffUsecase) JoinStaff(ctx context.Context, req *JoinStaffReq) (*JoinStaffRes, error) {

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		su.logger.Error("invalid team id", "error", errors.New("invalid user id from token"))
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	exist, err := su.teamPlayerRepo.IsTeamPlayerExist(ctx, userID)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, apperror.NewBadRequestError("user is an active player")
	}

	staffExist, err := su.staffRepo.IsStaffExistByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if staffExist {
		return nil, apperror.NewBadRequestError(
			"user is already registered as staff",
		)
	}

	invite, err := su.staffInviteRepo.GetStaffInvitationByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	if time.Now().After(invite.ExpiresAt) {
		return nil, apperror.NewFailedPreCondition(
			"invitation code expired",
		)
	}

	if invite.IsUsed {
		return nil, apperror.NewFailedPreCondition(
			"invitation code already used",
		)
	}

	staffDetail, err := su.staffInviteRepo.GetStaffDetails(ctx, invite.ID)
	if err != nil {
		return nil, err
	}

	if err := su.staffInviteRepo.MarkUsedStaffInvitation(ctx, staffDetail.ID); err != nil {
		return nil, err
	}

	teamMember, err := su.teamMemberRepo.AddTeamMember(ctx, &entity.TeamMember{
		ID:     uuid.New(),
		TeamID: invite.TeamID,
		UserID: userID,
		Role:   entity.TeamMemberStaff,
		Status: entity.TeamMemberStatusActive,
	})

	staff, err := su.staffRepo.AddStaff(ctx, &entity.Staff{
		ID:           uuid.New(),
		TeamMemberID: teamMember.ID,
		FullName:     staffDetail.FullName,
		Role:         staffDetail.Role,
		Designation:  staffDetail.Designation,
		PhoneNum:     *staffDetail.PhoneNume,
		Status:       entity.StaffStatusActive,
		JoinedAt:     time.Now(),
	})

	return &JoinStaffRes{
		ID:           staff.ID,
		TeamMemberID: teamMember.ID,
		UserID:       teamMember.UserID,
		FullName:     staff.FullName,
		Role:         string(staff.Role),
		Designatin:   string(staff.Designation),
		PhoneNum:     staff.PhoneNum,
		CreatedAT:    staff.CreatedAt,
	}, nil

}

func (su *staffUsecase) CompensateJoinStaff(ctx context.Context, req ComponsateJoinStaffReq) error {


	if err := su.staffInviteRepo.MarkUsedStaffInvitation(ctx, req.InviteID); err != nil {
		return err
	}

	if err := su.teamMemberRepo.DeleteTeamMemberByID(ctx, req.TeamMemberID); err != nil {
		return err
	}

	if err := su.staffRepo.DeleteStaff(ctx, req.StaffID); err != nil {
		return err
	}

	return nil
}

func (su *staffUsecase) GetStaff(ctx context.Context, req *GetStaffReq) (*GetStaffRes, error) {
	staffID, err := uuid.Parse(req.StaffID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid staff id")
	}

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	staff, err := su.staffRepo.GetStaff(ctx, teamID, staffID)
	if err != nil {
		return nil, err
	}

	return &GetStaffRes{
		ID:           staff.ID.String(),
		TeamMemberID: staff.TeamMemberID.String(),
		FullName:     staff.FullName,
		Role:         string(staff.Role),
		Designation:  string(staff.Designation),
		PhoneNum:     staff.PhoneNum,
		Status:       string(staff.Status),
		JoinedAt:     staff.JoinedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (su *staffUsecase) UploadStaffImage(ctx context.Context, req *UploadStaffImageReq) (*UploadStaffImageRes, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		su.logger.Error("invalid team id", "error", errors.New("invalid user id from token"))
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	staffID, err := uuid.Parse(req.StaffID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid staffid")
	}

	teamOwnerRole, err := su.staffRepo.GetStaffDesignation(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}

	permite := permission.HasPermission(teamOwnerRole, permission.PermissionAddStaff)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user not allowed to create staff")
	}

	staff, err := su.staffRepo.GetStaff(ctx, teamID, staffID)
	if err != nil {
		return nil, err
	}

	webbytes, err := imageutil.ConvertImageIntoWebpbFormate(req.ImageData)
	if err != nil {
		return nil, err
	}

	if staff.ImageKey == nil || *staff.ImageKey == "" {
		*staff.ImageKey = fmt.Sprintf("staff/%s/profile.webp", staff.ID)
		if err := su.staffRepo.UpdateImageKey(ctx, staffID, *staff.ImageKey); err != nil {
			return nil, err
		}
	}

	key, err := su.objstore.Upload(ctx, su.objcfg.Bucket, *staff.ImageKey, bytes.NewReader(webbytes), int64(len(webbytes)), "image/webp")

	if err != nil {
		return nil, err
	}

	presignedUrl, err := su.objstore.GetPresignedURL(ctx, su.objcfg.Bucket, key, su.objcfg.PresignedURLExpiry)
	if err != nil {
		return nil, err
	}

	return &UploadStaffImageRes{
		StaffID:      staff.ID.String(),
		PresignedUrl: presignedUrl,
	}, nil
}

func (tm *staffUsecase) GetStaffImagePresignedUrl(ctx context.Context, req *StaffImagePresignedUrlReq) (*StaffImagePresignedUrlRes, error) {

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	staffID, err := uuid.Parse(req.StaffID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid staff id")
	}

	staff, err := tm.staffRepo.GetStaff(ctx, teamID, staffID)
	if err != nil {
		return nil, err
	}

	if *staff.ImageKey == "" {
		return nil, apperror.NewNotFoundError("image is not added or removed")
	}

	presignedurl, err := tm.objstore.GetPresignedURL(ctx, tm.objcfg.Bucket, *staff.ImageKey, tm.objcfg.PresignedURLExpiry)
	if err != nil {
		return nil, err
	}

	return &StaffImagePresignedUrlRes{
		StaffID:      req.StaffID,
		PresignedUrl: presignedurl,
	}, nil
}

func (su *staffUsecase) RemoveStaffImage(ctx context.Context, req *RemoveStaffImageReq) (*RemoveStaffImageRes, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		su.logger.Error("invalid data", "error", errors.New("invalid user id from token"))
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	teamOwnerRole, err := su.staffRepo.GetStaffDesignation(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}

	permite := permission.HasPermission(teamOwnerRole, permission.PermissionAddStaff)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user not allowed to create staff")
	}

	staffID, err := uuid.Parse(req.StaffID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid staffid")
	}

	staff, err := su.staffRepo.GetStaff(ctx, teamID, staffID)
	if err != nil {
		return nil, err
	}

	err = su.objstore.Delete(ctx, su.objcfg.Bucket, *staff.ImageKey)
	if err != nil {
		return nil, err
	}

	if err := su.staffRepo.UpdateStaffImageURL(ctx, teamID, staffID); err != nil {
		return nil, err
	}

	return &RemoveStaffImageRes{
		Success: true,
	}, nil
}

func (su *staffUsecase) UpdateStaff(ctx context.Context, req *UpdateStaffReq) (*UpdateStaffRes, error) {

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		su.logger.Error("invalid user id from token", "error", err)
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	teamOwnerRole, err := su.staffRepo.GetStaffDesignation(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}

	permite := permission.HasPermission(teamOwnerRole, permission.PermissionAddStaff)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user not allowed to create staff")
	}

	staffID, err := uuid.Parse(req.StaffID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid staff id")
	}

	staff, err := su.staffRepo.GetStaff(ctx, teamID, staffID)
	if err != nil {
		return nil, err
	}

	existingRole := staff.Role
	existingDesignation := staff.Designation

	newRole := existingRole
	newDesignation := existingDesignation

	if req.Role != "" {
		newRole = req.Role
	}

	if req.Designation != "" {
		newDesignation = req.Designation
	}

	if newRole != existingRole || newDesignation != existingDesignation {
		if !IsValidStaffRoleDesignation(newRole, newDesignation) {
			return nil, apperror.NewBadRequestError("invalid designation for staff role")
		}
	}

	res, err := su.staffRepo.UpdateStaff(ctx, staffID, &entity.Staff{
		ID:          staffID,
		FullName:    req.FullName,
		Role:        req.Role,
		Designation: req.Designation,
		PhoneNum:    req.PhoneNum,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &UpdateStaffRes{
		StaffID:      res.ID,
		TeamMemberID: res.TeamMemberID,
		FullName:     res.FullName,
		Role:         string(res.Role),
		Designation:  string(res.Designation),
		PhoneNum:     res.PhoneNum,
	}, nil

}

func (su *staffUsecase) ListTeamStaff(ctx context.Context, req *ListTeamStaffReq) (*ListTeamStaffRes, error) {
	teamId, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	switch {
	case req.Limit > entity.MaxDefaultPaginateLimit:
		req.Limit = entity.MaxDefaultPaginateLimit
	case req.Limit <= 0:
		req.Limit = entity.MinDefaultPagination
	}

	res, total, err := su.staffRepo.ListStaff(ctx, teamId, &staffrepo.ListStaffDetails{
		Page:        req.Page,
		Limit:       req.Limit,
		Role:        req.Role,
		Designation: req.Designation,
		Search:      req.Search,
	})

	if err != nil {
		return nil, err
	}

	var staff []StaffDetails

	totalPage := total / req.Limit
	if total%req.Limit > 0 {
		totalPage += 1
	}

	for _, s := range res {
		staff = append(staff, StaffDetails{
			StaffID:      s.ID,
			TeamMemberID: s.TeamMemberID,
			FullName:     s.FullName,
			Role:         string(s.Role),
			Designation:  string(s.Designation),
			Status:       string(s.Status),
		})
		if s.ImageKey != nil && *s.ImageKey != "" {
			presignedUrl, err := su.objstore.GetPresignedURL(ctx, su.objcfg.Bucket, *s.ImageKey, su.objcfg.PresignedURLExpiry)
			if err != nil {
				su.logger.Error("object storage error", "error", err)
			} else {
				s.ImageKey = &presignedUrl
			}
		}
	}

	return &ListTeamStaffRes{
		Staff: staff,
		Pagination: &PaginationDetails{
			TotalPage: totalPage,
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     int64(total),
		},
	}, nil
}

func (su *staffUsecase) ReleaseStaff(ctx context.Context, req *ReleaseStaffReq) (*ReleaseStaffRes, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		su.logger.Error("invalid user id from token", "error", err)
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}
	teamOwnerRole, err := su.staffRepo.GetStaffDesignation(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}

	permite := permission.HasPermission(teamOwnerRole, permission.PermissionAddStaff)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user not allowed to create staff")
	}

	staffID, err := uuid.Parse(req.StaffID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid staff id")
	}

	if err := su.staffRepo.ReleaseStaff(ctx, teamID, staffID); err != nil {
		return nil, err
	}

	staff, err := su.staffRepo.GetStaff(ctx, teamID, staffID)
	if err != nil {
		return nil, err
	}

	return &ReleaseStaffRes{
		StaffID:      staff.ID.String(),
		TeamMemberID: staff.TeamMemberID.String(),
		FullName:     staff.FullName,
		Role:         string(staff.Role),
		Designation:  string(staff.Designation),
		Status:       string(staff.Status),
		JoinedAt:     staff.JoinedAt.Format("2006-01-02 15:04:05"),
	}, nil

}

func (tm *staffUsecase) TransferOwnership(ctx context.Context, req *TransferOwnershipReq) (*TransferOwnershipRes, error) {
	return nil, nil
}

func (tm *staffUsecase) GetStaffProfile(ctx context.Context, req *GetStaffProfileReq) (*GetStaffProfileRes, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid user id")
	}

	teamMember, err := tm.teamMemberRepo.GetActiveTeamMemberByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if teamMember.Role != entity.TeamMemberStaff {
		return nil, apperror.NewBadRequestError("user is not staff")
	}

	team, err := tm.teamRepo.GetTeamDetails(ctx, teamMember.TeamID)
	if err != nil {
		return nil, err
	}

	staff, err := tm.staffRepo.GetStaffByTeamMemberID(ctx, teamMember.ID)
	if err != nil {
		return nil, err
	}

	return &GetStaffProfileRes{
		Team: &TeamDetail{
			TeamID:    team.ID.String(),
			Name:      team.Name,
			ShortName: team.ShortName,
			City:      team.City,
			TeamCode:  team.TeamCode,
		},
		Staff: &StaffDetail{
			TeamMemberID: teamMember.ID.String(),
			StaffID:      staff.ID.String(),
			UserID:       teamMember.UserID.String(),
			TeamID:       teamMember.TeamID.String(),
			FullName:     staff.FullName,
			Role:         staff.Role,
			Designation:  staff.Designation,
			PhoneNum:     staff.PhoneNum,
			CreatedAt:    staff.CreatedAt,
		},
	}, nil
}
