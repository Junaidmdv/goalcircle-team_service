package teaminvite

import (
	"context"
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	tminviterepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teaminvite"
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
}

func NewTeamInviteUsecase(tm tminviterepo.TeamInviteRepository) TeamInviteUsecase {
	return &teamInviteUsecase{
		teamInviteRepo: tm,
	}
}

func (ti *teamInviteUsecase) CreateInvitation(ctx context.Context, req *TeamInviteReq) (*TeamInviteRes, error) {

	teamMemberID, err := uuid.Parse(req.TeamMemberID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team member id")
	}

    

	invite, exist, err := ti.teamInviteRepo.GetInvitation(ctx, &teamMemberID)
	if err != nil {
		return nil, err
	}

	if !exist || invite.IsUsed || time.Now().After(invite.ExpiresAt) {
		code, err := ti.CodeGenerater.GenerateCode(string(invite.Role))
		if err != nil {
			return nil, err
		}

		invite, err = ti.teamInviteRepo.CreateInvitation(ctx, &entity.TeamInvite{
			ID:           uuid.New(),
			TeamMemberID: teamMemberID,
			Code:         code,
			Role:         req.Role,
			IsUsed:       false,
			ExpiresAt:    time.Now().Add(entity.CodeExpiryDays * 24 * time.Hour),
			CreatedAt:    time.Now(),
		})
		if err != nil {
			return nil, err
		}
	}
	return &TeamInviteRes{
		TeamMemberId: invite.TeamMemberID,
		Code:         invite.Code,
		ExpiresAt:    invite.ExpiresAt,
		CreatedAt:    invite.CreatedAt,
	}, nil
}
