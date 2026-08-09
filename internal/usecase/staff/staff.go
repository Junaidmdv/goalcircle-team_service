package staffuc

import (
	"context"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	staffrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staff"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/datetime"
	"github.com/google/uuid"
)

type StaffUsecase interface {
	AddStaff(context.Context, *AddStaffReq) (*AddStaffRes, error)
	ChangeStaffStatus(context.Context, *ChangeStaffStatusReq) (*ChangeStaffStatusRes, error)
	ListTeamStaff(context.Context, *ListTeamStaffReq) (*ListTeamStaffRes, error)
}

type staffUsecase struct {
	teamMemberRepo teammember.TeamMemberRepository
	staffRepo      staffrepo.StaffRepository
	date           datetime.DateCalculator
}

func NewTeamStaffUsecase(tm teammember.TeamMemberRepository, sr staffrepo.StaffRepository,date datetime.DateCalculator) StaffUsecase {
	return &staffUsecase{
		teamMemberRepo: tm,
		staffRepo:      sr, 
		date: date,
	}
}

func (tm *staffUsecase) AddStaff(ctx context.Context, req *AddStaffReq) (*AddStaffRes, error) {

	teamId, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	teamMember, err := tm.teamMemberRepo.AddTeamMember(ctx, &entity.TeamMember{
		ID:       uuid.New(),
		TeamID:   teamId,
		FullName: req.FullName,
		Role:     req.Role,
	})

	if err != nil {
		return nil, err
	}

	staff, err := tm.staffRepo.AddStaff(ctx, &entity.Staff{
		ID:           uuid.New(),
		TeamMemberID: teamMember.ID,
		FullName:     req.FullName,
		DateOfBirth:  req.DOB,
		Designation:  req.Designation,
		Status:       entity.StaffStatusPending,
	})
	if err != nil {
		return nil, err
	}

	age := tm.date.CalculateAge(req.DOB)

	return &AddStaffRes{
		StaffID:      staff.ID,
		TeamMemberID: teamMember.ID,
		FullName:     staff.FullName,
		Age:          age,
		Designation:  req.Designation,
	}, nil
}

func (tm *staffUsecase) ChangeStaffStatus(ctx context.Context, req *ChangeStaffStatusReq) (*ChangeStaffStatusRes, error) {
	return nil, nil
}

func (tm *staffUsecase) ListTeamStaff(ctx context.Context, req *ListTeamStaffReq) (*ListTeamStaffRes, error) {
	return nil, nil
}
