package teamowner

import (
	"context"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/google/uuid"
)

type TeamOwnerUsecase interface {
	AddTeamOwner(context.Context, *AddTeamOwnerReq) (*AddTeamOwnerRes, error)
	DeleteTeamOwner(context.Context, *uuid.UUID) error
}

type teamMemberUsecase struct {
	teamMemberRepo teammember.TeamMemberRepository
}

func NewTeamOwnerUsecase(tm teammember.TeamMemberRepository) TeamOwnerUsecase {
	return &teamMemberUsecase{
		teamMemberRepo: tm,
	}
}

func (tm *teamMemberUsecase) AddTeamOwner(ctx context.Context, req *AddTeamOwnerReq) (*AddTeamOwnerRes, error) {
	res, err := tm.teamMemberRepo.AddTeamMember(ctx, &entity.TeamMember{
		ID:       uuid.New(),
		TeamID:   req.TeamID,
		FullName: req.FullName,
		Role:     entity.OWNER,
	})

	if err != nil {
		return nil, err
	}

	return &AddTeamOwnerRes{
		TeamMemberID: res.ID,
		TeamID:       res.TeamID,
		UserID:       res.UserID,
		FullName:     req.FullName,
		Role:         res.Role,
	}, nil
}

func (tm *teamMemberUsecase) DeleteTeamOwner(ctx context.Context, UserId *uuid.UUID) error {
	return tm.teamMemberRepo.RemoveTeamMember(ctx, UserId)
}

func (tm *teamMemberUsecase) AddPlayer(ctx context.Context, req *AddPlayerReq) (*AddPlayerRes, error) {

	tm.teamMemberRepo.AddTeamMember(ctx, &entity.TeamMember{})

	return nil, nil
}

// func (tm *teamMemberUsecase) AddStaff(ctx context.Context, req *AddStaffReq) (*AddStaffRes, error) {
// 	return nil, nil
// }

// func (tm *teamMemberUsecase) AddManager(ctx context.Context, req *AddManagerReq) (*AddManagerRes, error) {
// 	return nil, nil
// }
