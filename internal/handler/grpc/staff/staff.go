package staff

import (
	"context"
	"time"

	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	teamv1 "github.com/Junaidmdv/goalcircle-protos/team/v1"
	staffuc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/staff"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/validater"
)

type StaffHandler struct {
	teamv1.UnimplementedStaffServiceServer
	StaffUsecase staffuc.StaffUsecase
	timeOut      time.Duration
	Validater    validater.Validater
}

func NewStaffHandler(stafuc staffuc.StaffUsecase, timeout time.Duration) *StaffHandler {
	return &StaffHandler{
		StaffUsecase: stafuc,
		timeOut:      timeout,
	}

}

func (sh *StaffHandler) AddStaff(ctx context.Context, req *pb.AddStaffReq) (*pb.AddStaffRes, error) {
	context, cancel := context.WithTimeout(ctx, sh.timeOut)
	defer cancel()

	role := MapTeamMemberRole(req.Role)
	designation := MapStaffDesignation(req.Designation)

	res, err := sh.StaffUsecase.AddStaff(context, &staffuc.AddStaffReq{
		UserID:      req.UserId,
		TeamID:      req.TeamId,
		FullName:    req.FullName,
		Designation: designation,
		DOB:         req.Dob.AsTime(),
		Role:        role,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.AddStaffRes{
		StaffId:      res.StaffID.String(),
		TeamMemberId: res.TeamMemberID.String(),
		FullName:     res.FullName,
		Age:          res.Age,
		Designatin:   string(res.Designation),
	}, nil
}



func (tm *StaffHandler) ChangeStaffStatus(ctx context.Context){

}

func (tm *StaffHandler) ListTeamStaff(ctx context.Context){

}


// func(tm *StaffHandler)
