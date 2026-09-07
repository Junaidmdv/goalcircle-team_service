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
	teamrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/datetime"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/imageutil"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
)

type StaffUsecase interface {
	AddStaff(context.Context, *AddStaffReq) (*AddStaffRes, error)
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
	teamMemberRepo teammember.TeamMemberRepository
	staffRepo      staffrepo.StaffRepository
	teamRepo       teamrepo.TeamRepository
	date           datetime.DateCalculator
	objstore       storage.ObjectStorage
	objcfg         *config.ObjectStorageConfig
	logger         logger.Logger
}

func NewTeamStaffUsecase(tm teammember.TeamMemberRepository,
	sr staffrepo.StaffRepository,
	tmrepo teamrepo.TeamRepository,
	date datetime.DateCalculator,
	objst storage.ObjectStorage,
	cnfg *config.ObjectStorageConfig,
	logger logger.Logger) StaffUsecase {
	return &staffUsecase{
		teamMemberRepo: tm,
		staffRepo:      sr,
		teamRepo:       tmrepo,
		date:           date,
		objstore:       objst,
		objcfg:         cnfg,
		logger:         logger,
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

	teamMember, err := su.teamMemberRepo.AddTeamMember(ctx, &entity.TeamMember{
		ID:     uuid.New(),
		TeamID: teamID,
		Status: req.TeamMemberStatus,
	})

	if err != nil {
		return nil, err
	}

	count, err := su.staffRepo.CountByDesination(ctx, teamID, req.Designation)
	if err != nil {
		return nil, err
	}

	if err := StaffMaxLimit(req.Designation, count); err != nil {
		return nil, err
	}

	staff, err := su.staffRepo.AddStaff(ctx, &entity.Staff{
		ID:           uuid.New(),
		TeamMemberID: teamMember.ID,
		FullName:     req.FullName,
		Role:         req.Role,
		Designation:  req.Designation,
		PhoneNum:     *req.PhoneNum,
	})
	if err != nil {
		return nil, err
	}

	webbytes, err := imageutil.ConvertImageIntoWebpbFormate(req.ImageBytes)
	if err != nil {
		return nil, err
	}

	objectName := fmt.Sprintf("staff/%s/profile.webp", staff.ID)

	key, err := su.objstore.Upload(ctx, su.objcfg.Bucket, objectName, bytes.NewReader(webbytes), int64(len(webbytes)), "image/webp")

	if err != nil {
		return nil, err
	}

	if err := su.staffRepo.UpdateImageKey(ctx, staff.ID, key); err != nil {
		return nil, err
	}

	presignedUrl, err := su.objstore.GetPresignedURL(ctx, su.objcfg.Bucket, key, su.objcfg.PresignedURLExpiry)
	if err != nil {
		return nil, err
	}

	return &AddStaffRes{
		StaffID:      staff.ID,
		TeamMemberID: teamMember.ID,
		FullName:     staff.FullName,
		Role:         staff.Role,
		Designation:  req.Designation,
		PresignedUrl: presignedUrl,
	}, nil
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
		Status:       string(staff.TeamMember.Status),
		JoinedAt:     staff.TeamMember.JoinedAt.Format("2006-01-02 15:04:05"),
		ReleasedAt:   staff.TeamMember.ReleasedAt.Format("2006-01-02 15:04:05"),
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

	for i, s := range res {
		staff = append(staff, StaffDetails{
			StaffID:      s.ID,
			TeamMemberID: s.TeamMemberID,
			FullName:     s.FullName,
			Role:         string(s.Role),
			Designation:  string(s.Designation),
			Status:       string(s.TeamMember.Status),
		})

		presignedUrl, err := su.objstore.GetPresignedURL(ctx, su.objcfg.Bucket, *s.ImageKey, su.objcfg.PresignedURLExpiry)
		if err == nil {
			staff[i].PresignedUrl = presignedUrl
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
		Status:       string(staff.TeamMember.Status),
		JoinedAt:     staff.TeamMember.JoinedAt.Format("2006-01-02 15:04:05"),
		ReleasedAt:   staff.TeamMember.ReleasedAt.Format("2006-01-02 15:04:05"),
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
