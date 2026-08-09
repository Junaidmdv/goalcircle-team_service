package teaminvite

import (
	"context"
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/permission"
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
}

func NewTeamInviteUsecase(tm tminviterepo.TeamInviteRepository, code code.CodeGenerater, teamMemberRepo teammember.TeamMemberRepository) TeamInviteUsecase {
	return &teamInviteUsecase{
		teamInviteRepo: tm,
		CodeGenerater:  code,
		teamMemberRepo: teamMemberRepo,
	}
}

func (ti *teamInviteUsecase) CreateInvitation(ctx context.Context, req *TeamInviteReq) (*TeamInviteRes, error) {

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid user id")
	}

	teamMemberID, err := uuid.Parse(req.TeamMemberID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid user id")
	}

	authorised, err := ti.teamMemberRepo.GetActiveTeamMemberByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	permite := permission.HasPermissionTeam(authorised.Role, permission.PermissionCreateInvite)

	if !permite {
		return nil, apperror.NewFailedPreCondition("user is not allowed to create invitation")
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
		code, err := ti.CodeGenerater.GenerateCode(string(teamMember.Role))
		if err != nil {
			return nil, err
		}
        
		invite, err = ti.teamInviteRepo.CreateInvitation(ctx, &entity.TeamInvite{
			ID:           uuid.New(),
			TeamMemberID: teamMember.ID,
			Code:         code,
			Role:         string(teamMember.Role),
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
