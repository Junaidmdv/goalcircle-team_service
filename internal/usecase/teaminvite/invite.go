package teaminvite

import (
	"context"
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/permission"
	staffrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staff"
	tminviterepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teaminvite"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	code "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/invitation"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/google/uuid"
)

type TeamInviteUsecase interface {
	CreateInvitation(context.Context, *TeamInviteReq) (*TeamInviteRes, error)
}

type teamInviteUsecase struct {
	teamInviteRepo tminviterepo.TeamInviteRepository
	CodeGenerater  code.CodeGenerater
	teamMemberRepo teammember.TeamMemberRepository
	staffrepo      staffrepo.StaffRepository
}

func NewTeamInviteUsecase(tm tminviterepo.TeamInviteRepository, code code.CodeGenerater, teamMemberRepo teammember.TeamMemberRepository, staffrepo staffrepo.StaffRepository) TeamInviteUsecase {
	return &teamInviteUsecase{
		teamInviteRepo: tm,
		CodeGenerater:  code,
		teamMemberRepo: teamMemberRepo,
		staffrepo:      staffrepo,
	}
}

func (ti *teamInviteUsecase) CreateInvitation(ctx context.Context, req *TeamInviteReq) (*TeamInviteRes, error) {

	//left to add team id  add it later

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid user id")
	}

	teamOwnerRole, err := ti.staffrepo.GetStaffDesignation(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}

	permite := permission.HasPermission(teamOwnerRole, permission.PermissionAddStaff)
	if !permite {
		return nil, apperror.NewUnAuthenticatedError("user not allowed to create staff")
	}

	teamMemberID, err := uuid.Parse(req.TeamMemberID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid user id")
	}

	teamMember, err := ti.teamMemberRepo.GetTeamMemberByID(ctx, teamMemberID)
	if err != nil {
		return nil, err
	}

	invite, exist, err := ti.teamInviteRepo.GetInvitation(ctx, &teamMember.ID)
	if err != nil {
		return nil, err
	}

	if !exist || invite.IsUsed || time.Now().After(invite.ExpiresAt) {
		code, err := ti.CodeGenerater.GenerateCode("TEAM")
		if err != nil {
			return nil, err
		}

		invite, err = ti.teamInviteRepo.CreateInvitation(ctx, &entity.TeamInvite{
			ID:           uuid.New(),
			TeamMemberID: teamMember.ID,
			Code:         code,
			IsUsed:       false,
			ExpiresAt:    time.Now().Add(entity.CodeExpiryDays * 24 * time.Hour),
			CreatedAt:    time.Now(),
		})
		if err != nil {
			return nil, err
		}
	}
	return &TeamInviteRes{
		TeamMemberId: teamMember.ID,
		Code:         invite.Code,
		ExpiresAt:    invite.ExpiresAt,
		CreatedAt:    invite.CreatedAt,
	}, nil
}
