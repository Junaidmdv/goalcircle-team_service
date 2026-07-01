package teammemberuc

import (
	"context"
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teaminvite"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
)

type TeamMemberUsecase interface {
	AddTeamOwner(context.Context, *AddTeamOwnerReq) (*AddTeamOwnerRes, error)
	DeleteTeamOwner(context.Context, *uuid.UUID) error
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

func (tm *teamMemberUsecase) AddTeamOwner(ctx context.Context, req *AddTeamOwnerReq) (*AddTeamOwnerRes, error) {
	res, err := tm.teamMemberRepo.AddTeamMember(ctx, &entity.TeamMember{
		ID:       uuid.New(),
		TeamID:   req.TeamID,
		FullName: req.FullName,
		Role:     entity.TeamMemberRoleOwner,
	})

	if err != nil {
		return nil, err
	}

	return &AddTeamOwnerRes{
		TeamMemberID: res.ID,
		TeamID:       res.TeamID,
		UserID:       res.UserID,
		FullName:     res.FullName,
		Role:         res.Role,
	}, nil
}

func (tm *teamMemberUsecase) RegisterTeamMember(ctx context.Context, req *RegisterTeamMemberReq) (*RegisterTeamMemberRes, error) {

	invite, err := tm.teamInviteRepo.GetInvitationByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	if time.Now().After(invite.ExpiresAt) {
		return nil, apperror.NewFailedPreCondition("invitation code expired")
	}

	if invite.IsUsed {
		return nil, apperror.NewFailedPreCondition("invitation code already used")
	}

	if err := tm.teamMemberRepo.UpdateUserID(ctx, &invite.TeamMemberID, req.UserID); err != nil {
		return nil, err
	}
	if err := tm.teamInviteRepo.UpdateTeamInvite(ctx, invite.ID, &teaminvite.UpdateTeamInviteReq{
		UserID: req.UserID,
		IsUsed: true,
	}); err != nil {
		return nil, err

	}

	return &RegisterTeamMemberRes{
		InvitationID: invite.ID,
	}, nil

}

func (tm *teamMemberUsecase) CompensateRegisterTeamMember(ctx context.Context, req *CompensateRegisterTeamMemberReq) error {
	if err := tm.teamMemberRepo.UpdateUserID(ctx, &req.InvitationID, ""); err != nil {
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
