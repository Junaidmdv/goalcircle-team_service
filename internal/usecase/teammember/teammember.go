package teammemberuc

import (
	"context"
	"errors"
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teaminvite"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
)

type TeamMemberUsecase interface {
	RegisterTeamMember(context.Context, *RegisterTeamMemberReq) (*RegisterTeamMemberRes, error)
	CompensateRegisterTeamMember(context.Context, *CompensateRegisterTeamMemberReq) error
}

type teamMemberUsecase struct {
	teamMemberRepo teammember.TeamMemberRepository
	teamInviteRepo teaminvite.TeamInviteRepository
	logger         logger.Logger
}

func NewTeamMemberUsecase(tmr teammember.TeamMemberRepository, ti teaminvite.TeamInviteRepository, logger logger.Logger) TeamMemberUsecase {
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

func (tm *teamMemberUsecase) DeleteTeamOwner(ctx context.Context, teamMemberID *uuid.UUID) error {
	return tm.teamMemberRepo.RemoveTeamMember(ctx, teamMemberID)
}
