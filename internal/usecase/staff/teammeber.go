package staff

import (
	"context"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/google/uuid"
)

type StaffUsecase interface {
	AddTeamOwner(context.Context, *AddTeamOwnerReq) (*AddTeamOwnerRes, error)
	DeleteStaff(context.Context, *uuid.UUID) error
}

type staffUsecase struct {
	teamMemberRepo teammember.TeamMemberRepository
}

func NewTeamStaffUsecase(tm teammember.TeamMemberRepository) StaffUsecase {
	return &staffUsecase{
		teamMemberRepo: tm,
	}
}

func (tm *staffUsecase) AddTeamOwner(ctx context.Context, req *AddTeamOwnerReq) (*AddTeamOwnerRes, error) {
	res, err := tm.teamMemberRepo.AddTeamMember(ctx, &entity.TeamMember{
		ID:       uuid.New(),
		TeamID:   req.TeamID,
		FullName: req.FullName,
		Role:     req.Role,
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

func (tm *staffUsecase) DeleteStaff(ctx context.Context, staffID *uuid.UUID) error {
	return tm.teamMemberRepo.RemoveTeamMember(ctx, staffID)
}

func (tm *staffUsecase) AddManager(ctx context.Context,req *AddManagerReq)(*AddManagerRes,error){
	return &AddManagerRes{},nil
}


func(tm *staffUsecase)AddCoach(ctx context.Context)