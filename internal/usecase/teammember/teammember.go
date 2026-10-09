package teammemberuc

import (
	"context"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	staffrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staff"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/datetime"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/google/uuid"
)

type TeamMemberUsecase interface {
	AddOwner(ctx context.Context, req *AddOwnerReq) (*AddOwnerRes, error)
	DeleteTeamOwner(ctx context.Context, teamMemberID uuid.UUID, staffID uuid.UUID) error
}

type teamMemberUsecase struct {
	teamMemberRepo teammember.TeamMemberRepository
	StaffRepo      staffrepo.StaffRepository
	logger         logger.Logger
	dateCalc       datetime.DateCalculator
}

func NewTeamMemberUsecase(tmr teammember.TeamMemberRepository, staffRepo staffrepo.StaffRepository, logger logger.Logger, datecal datetime.DateCalculator) TeamMemberUsecase {
	return &teamMemberUsecase{
		teamMemberRepo: tmr,
		logger:         logger,
		StaffRepo:      staffRepo,
		dateCalc:       datecal,
	}
}

func (tm *teamMemberUsecase) AddOwner(ctx context.Context, req *AddOwnerReq) (*AddOwnerRes, error) {
	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid user id")
	}

	teamMember, err := tm.teamMemberRepo.AddTeamMember(ctx, &entity.TeamMember{
		ID:     uuid.New(),
		TeamID: teamID,
		Role:   entity.TeamMemberStaff,
		UserID: userID,
	})

	if err != nil {
		return nil, err
	}

	staff, err := tm.StaffRepo.AddStaff(ctx, &entity.Staff{
		ID:           uuid.New(),
		TeamMemberID: teamMember.ID,
		FullName:     req.FullName,
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

func (tm *teamMemberUsecase) DeleteTeamOwner(ctx context.Context, teamMemberID uuid.UUID, staffID uuid.UUID) error {
	err := tm.teamMemberRepo.RemoveTeamMember(ctx, teamMemberID)
	if err != nil {
		return err
	}

	err = tm.StaffRepo.DeleteStaff(ctx, staffID)
	if err != nil {
		return err
	}
	return err

}

