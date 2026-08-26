package teammemberuc

import (
	"context"
	"errors"
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	staffrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staff"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teaminvite"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/datetime"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
)

type TeamMemberUsecase interface {
	AddOwner(ctx context.Context, req *AddOwnerReq) (*AddOwnerRes, error)
	RegisterTeamMember(context.Context, *RegisterTeamMemberReq) (*RegisterTeamMemberRes, error)
	CompensateRegisterTeamMember(context.Context, *CompensateRegisterTeamMemberReq) error
	DeleteTeamOwner(ctx context.Context, teamMemberID *uuid.UUID) error
}

type teamMemberUsecase struct {
	teamMemberRepo teammember.TeamMemberRepository
	teamInviteRepo teaminvite.TeamInviteRepository
	StaffRepo      staffrepo.StaffRepository
	logger         logger.Logger
	dateCalc       datetime.DateCalculator
}

func NewTeamMemberUsecase(tmr teammember.TeamMemberRepository, ti teaminvite.TeamInviteRepository, staffRepo staffrepo.StaffRepository, logger logger.Logger) TeamMemberUsecase {
	return &teamMemberUsecase{
		teamMemberRepo: tmr,
		teamInviteRepo: ti,
		logger:         logger,
	}
}

func (tm *teamMemberUsecase) RegisterTeamMember(ctx context.Context, req *RegisterTeamMemberReq) (*RegisterTeamMemberRes, error) {

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		tm.logger.Error("token error", "error", errors.New("invalid user id from token"))
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	hasMembership, err := tm.teamMemberRepo.HasUnreleasedMembership(ctx, userID)
	if err != nil {
		return nil, err
	}

	if hasMembership {
		return nil, apperror.NewBadRequestError("user already has membershipt with another team")
	}

	invite, err := tm.teamInviteRepo.GetInvitationByCode(ctx, req.Code)
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

	if err := tm.teamMemberRepo.UpdateUserID(
		ctx,
		invite.TeamMemberID,
		userID,
	); err != nil {
		return nil, err
	}

	if err := tm.teamInviteRepo.UpdateTeamInvite(
		ctx,
		invite.ID,
		&teaminvite.UpdateTeamInviteReq{
			UserID: req.UserID,
			IsUsed: true,
		},
	); err != nil {
		return nil, err
	}

	member, err := tm.teamMemberRepo.GetTeamMemberByID(
		ctx,
		invite.TeamMemberID,
	)
	if err != nil {
		return nil, err
	}

	return &RegisterTeamMemberRes{
		InvitationID: invite.ID,
		TeamMemberID: member.ID,
		TeamID:       member.TeamID,
		UserID:       userID,
		Role:         string(member.Role),
		Status:       string(member.Status),
		JoinedAt:     member.JoinedAt,
	}, nil
}

func (tm *teamMemberUsecase) CompensateRegisterTeamMember(ctx context.Context, req *CompensateRegisterTeamMemberReq) error {
	if err := tm.teamMemberRepo.UpdateUserID(ctx, req.TeamMemberID, req.UserID); err != nil {
		return err
	}
	if err := tm.teamInviteRepo.UpdateTeamInvite(ctx, req.InvitationID, &teaminvite.UpdateTeamInviteReq{
		UserID: "",
		IsUsed: false,
	}); err != nil {
		return err

	}
	return nil
}

func (tm *teamMemberUsecase) AddOwner(ctx context.Context, req *AddOwnerReq) (*AddOwnerRes, error) {
	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	teamMember, err := tm.teamMemberRepo.AddTeamMember(ctx, &entity.TeamMember{
		ID:     uuid.New(),
		TeamID: teamID,
		Role:   entity.TeamMemberRoleStaff,
		Status: req.TeamMemberStatus,
	})

	if err != nil {
		return nil, err
	}

	staff, err := tm.StaffRepo.AddStaff(ctx, &entity.Staff{
		ID:           uuid.New(),
		TeamMemberID: teamMember.ID,
		FullName:     req.FullName,
		DateOfBirth:  req.DOB,
		Role:         req.Role,
		Designation:  req.Designation,
	})
	if err != nil {
		return nil, err
	}

	age := tm.dateCalc.CalculateAge(req.DOB)

	return &AddOwnerRes{
		StaffID:      staff.ID,
		TeamMemberID: teamMember.ID,
		FullName:     staff.FullName,
		Age:          age,
		Role:         staff.Role,
		Designation:  req.Designation,
	}, nil
}

func (tm *teamMemberUsecase) DeleteTeamOwner(ctx context.Context, teamMemberID *uuid.UUID) error {
	return tm.teamMemberRepo.RemoveTeamMember(ctx, teamMemberID)
}

func (tm *teamMemberUsecase) LeaveTeam(ctx context.Context, req *LeaveTeamReq) (*LeaveTeamRes, error) {

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewInvalidArgumentError("invalid team id")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid user id")
	}

	teamMember, err := tm.teamMemberRepo.GetTeamMember(ctx, teamID, userID)

	if err != nil {
		return nil, err
	}

	if teamMember.Status == entity.TeamMemberStatusRelease && teamMember.ReleasedAt.IsZero() {
		return nil, apperror.NewBadRequestError("team member can leave the team only after being released")
	}

	return nil, nil

}
