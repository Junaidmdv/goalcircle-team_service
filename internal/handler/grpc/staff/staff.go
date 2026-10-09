package staff

import (
	"bytes"
	"context"
	"io"
	"time"

	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	teamv1 "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	teamsaga "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/saga/team"
	staffuc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/staff"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/imageutil"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/validater"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StaffHandler struct {
	StaffUsecase staffuc.StaffUsecase
	timeOut      time.Duration
	validater    *validater.Validater
	joinsaga     teamsaga.JoinStaffSaga
	teamv1.UnimplementedStaffServiceServer
}

func NewStaffHandler(stafuc staffuc.StaffUsecase, timeout time.Duration, validater *validater.Validater, jsaga teamsaga.JoinStaffSaga) *StaffHandler {
	return &StaffHandler{
		StaffUsecase: stafuc,
		timeOut:      timeout,
		validater:    validater,
		joinsaga:     jsaga,
	}

}

func (sh *StaffHandler) AddStaff(ctx context.Context, req *pb.AddStaffReq) (*pb.AddStaffRes, error) {

	request := ToCreateStaff(req)

	if validationErrs := sh.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := validater.ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	staff, err := sh.StaffUsecase.AddStaff(ctx, &staffuc.AddStaffReq{
		TeamID:      request.TeamID,
		UserID:      request.UserID,
		FullName:    request.FullName,
		Role:        request.Role,
		Designation: request.Designation,
		PhoneNum:    request.PhoneNum,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.AddStaffRes{
		Code:        staff.Code,
		FullName:    staff.FullName,
		Role:        string(staff.Role),
		Designation: string(staff.Designation),
	}, nil

}

func (sh *StaffHandler) JoinStaff(ctx context.Context, req *pb.JoinStaffReq) (*pb.JoinStaffRes, error) {
	res, err := sh.joinsaga.RegisterStaff(ctx, &teamsaga.JoinStaffSagaState{
		JoinStaffReq: &staffuc.JoinStaffReq{
			UserID: req.UserId,
			Code:   req.Code,
		},
		AddUserRoleReq: &teamsaga.AddUserRoleReq{
			UserID:       req.UserId,
			Role:         string(entity.TEAM),
			RefreshToken: req.RefreshToken,
		},
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.JoinStaffRes{
		Staff: &teamv1.JoinStaffRes_StaffDetails{
			StaffId:      res.JoinStaffRes.ID.String(),
			TeamMemberId: res.JoinStaffRes.TeamMemberID.String(),
			UserId:       res.JoinStaffRes.UserID.String(),
			FullName:     res.JoinStaffRes.FullName,
			Role:         res.JoinStaffRes.Role,
			Designation:  res.JoinStaffRes.Designatin,
			CreatedAt:    timestamppb.New(res.JoinStaffRes.CreatedAT),
		},
		User: &teamv1.JoinStaffRes_UserResponse{
			SessionId:          res.AddUserRoleRes.SessionID,
			UserId:             res.AddUserRoleRes.UserID,
			Email:              res.AddUserRoleRes.Email,
			AccessToken:        res.AddUserRoleRes.AccessToken,
			AccessTokenExpiry:  timestamppb.New(res.AddUserRoleRes.AccessTokenExpiry),
			RefreshToken:       res.AddUserRoleRes.RefreshToken,
			RefreshTokenExpiry: timestamppb.New(res.AddUserRoleRes.RefreshTokenExpiry),
		},
	}, nil
}

func (tm *StaffHandler) ListTeamStaff(ctx context.Context, input *pb.ListTeamStaffReq) (*pb.ListTeamStaffRes, error) {
	ctx, cancel := context.WithTimeout(ctx, tm.timeOut)
	defer cancel()

	role := MapStaffRole(input.Role)
	desig := MapStaffDesignation(input.Designation)

	res, err := tm.StaffUsecase.ListTeamStaff(ctx, &staffuc.ListTeamStaffReq{
		TeamID:      input.TeamId,
		Page:        input.Page,
		Limit:       input.Limit,
		Role:        role,
		Designation: desig,
		Search:      input.Search,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	var staffList []*pb.StaffDetailResponse

	for _, v := range res.Staff {
		staffList = append(staffList, &pb.StaffDetailResponse{
			StaffId:      v.StaffID.String(),
			TeamMemberId: v.TeamMemberID.String(),
			FullName:     v.FullName,
			Role:         v.Role,
			Designation:  v.Designation,
			PresignedUrl: v.PresignedUrl,
			Status:       v.Status,
		})
	}

	return &teamv1.ListTeamStaffRes{
		StaffDetail: staffList,
		Pagination: &teamv1.StaffPaginationDetails{
			TotalPage: res.Pagination.TotalPage,
			Page:      res.Pagination.Page,
			Limit:     res.Pagination.Limit,
			Total:     res.Pagination.Total,
		},
	}, nil
}

func (sh *StaffHandler) UpdateStaff(ctx context.Context, input *pb.UpdateStaffReq) (*pb.UpdateStaffRes, error) {
	ctx, cancel := context.WithTimeout(ctx, sh.timeOut)
	defer cancel()

	request := ToUpdateStaff(input)

	if validationErrs := sh.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := validater.ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := sh.StaffUsecase.UpdateStaff(ctx, &staffuc.UpdateStaffReq{
		UserID:      request.UserID,
		TeamID:      request.TeamID,
		StaffID:     request.StaffID,
		FullName:    *request.FullName,
		Role:        *request.Role,
		Designation: *request.Designation,
		PhoneNum:    *request.PhoneNum,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.UpdateStaffRes{
		StaffId:      res.StaffID.String(),
		TeamMemberId: res.TeamMemberID.String(),
		FullName:     res.FullName,
		DateOfBirth:  timestamppb.New(res.DateOfBirth),
		Role:         res.Role,
		Designation:  res.Designation,
		PhoneNum:     res.PhoneNum,
		CreatedAt:    timestamppb.New(res.CreatedAt),
		UpdatedAt:    timestamppb.New(res.UpdatedAt),
	}, nil
}

func (sh *StaffHandler) GetStaff(ctx context.Context, input *pb.GetStaffReq) (*pb.GetStaffRes, error) {

	res, err := sh.StaffUsecase.GetStaff(ctx, &staffuc.GetStaffReq{
		StaffID: input.StaffId,
		TeamID:  input.TeamId,
	})
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.GetStaffRes{
		Id:           res.ID,
		TeamMemberId: res.TeamMemberID,
		FullName:     res.FullName,
		Role:         res.Role,
		Designation:  res.Designation,
		PhoneNum:     res.PhoneNum,
		CreatedAt:    timestamppb.New(res.CreatedAt),
		Status:       res.Status,
		JoinedAt:     res.JoinedAt,
	}, nil
}

func (sh *StaffHandler) UpdateStaffImage(stream grpc.ClientStreamingServer[pb.UpdateStaffImageReq, pb.UpdateStaffImageRes]) error {
	c := stream.Context()
	ctx, cancel := context.WithTimeout(c, sh.timeOut)
	defer cancel()

	var (
		meta   *pb.StaffImageMeta
		buffer bytes.Buffer
	)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return status.Errorf(codes.Internal, "failed to receive stream: %v", err)
		}

		switch data := req.Data.(type) {
		case *pb.UpdateStaffImageReq_Meta:
			if meta != nil {
				return status.Error(codes.InvalidArgument, "staff details already sent")
			}
			meta = data.Meta

		case *pb.UpdateStaffImageReq_Chunks:
			if meta == nil {
				return status.Error(codes.InvalidArgument, "player details is missing")
			}
			buffer.Write(data.Chunks)
		}
	}

	if err := imageutil.ValidateImage(buffer.Bytes(), imageutil.StaffImage); err != nil {
		return apperror.GRPCStatus(err)
	}

	res, err := sh.StaffUsecase.UploadStaffImage(ctx, &staffuc.UploadStaffImageReq{
		UserID:    meta.UserId,
		TeamID:    meta.TeamId,
		StaffID:   meta.StaffId,
		ImageData: buffer.Bytes(),
	})

	if err != nil {
		return apperror.GRPCStatus(err)
	}

	return stream.SendAndClose(&pb.UpdateStaffImageRes{
		StaffId:      res.StaffID,
		PresignedUrl: res.PresignedUrl,
	})
}

func (sh *StaffHandler) GetStaffImageUrl(ctx context.Context, input *pb.GetStaffImageUrlReq) (*pb.GetStaffImageUrlRes, error) {
	ctx, cancel := context.WithTimeout(ctx, sh.timeOut)
	defer cancel()

	res, err := sh.StaffUsecase.GetStaffImagePresignedUrl(ctx, &staffuc.StaffImagePresignedUrlReq{
		StaffID: input.StaffId,
		TeamID:  input.TeamId,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.GetStaffImageUrlRes{
		StaffId:      res.StaffID,
		PresignedUrl: res.PresignedUrl,
	}, nil
}

func (sh *StaffHandler) RemoveStaffImage(ctx context.Context, input *pb.RemoveStaffImageReq) (*pb.RemoveStaffImageRes, error) {
	ctx, cancel := context.WithTimeout(ctx, sh.timeOut)
	defer cancel()

	res, err := sh.StaffUsecase.RemoveStaffImage(ctx, &staffuc.RemoveStaffImageReq{
		UserID:  input.UserId,
		TeamID:  input.TeamId,
		StaffID: input.StaffId,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}
	return &pb.RemoveStaffImageRes{
		Success: res.Success,
	}, nil
}

func (sh *StaffHandler) ReleaseStaff(ctx context.Context, input *pb.ReleaseStaffReq) (*pb.ReleaseStaffRes, error) {
	ctx, cancel := context.WithTimeout(ctx, sh.timeOut)
	defer cancel()

	res, err := sh.StaffUsecase.ReleaseStaff(ctx, &staffuc.ReleaseStaffReq{
		TeamID:  input.TeamId,
		UserID:  input.UserId,
		StaffID: input.StaffId,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.ReleaseStaffRes{
		StaffId:      res.StaffID,
		TeamMemberId: res.TeamMemberID,
		FullName:     res.FullName,
		Role:         res.Role,
		Designation:  res.Designation,
		Status:       res.Status,
		JoinedAt:     res.JoinedAt,
		ReleasedAt:   res.ReleasedAt,
	}, nil
}

func (sh *StaffHandler) GetStaffProfile(ctx context.Context, input *pb.GetStaffProfileReq) (*pb.GetStaffProfileRes, error) {
	ctx, cancel := context.WithTimeout(ctx, sh.timeOut)
	defer cancel()

	res, err := sh.StaffUsecase.GetStaffProfile(ctx, &staffuc.GetStaffProfileReq{
		UserID: input.UserId,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.GetStaffProfileRes{
		Team: &pb.GetStaffProfileRes_TeamDetails{
			TeamId:    res.Team.TeamID,
			Name:      res.Team.Name,
			ShortName: res.Team.ShortName,
			City:      res.Team.City,
			TeamCode:  res.Team.TeamCode,
		},
		Staff: &pb.GetStaffProfileRes_StaffDetails{
			StaffId:      res.Staff.StaffID,
			TeamMemberId: res.Staff.TeamMemberID,
			UserId:       res.Staff.TeamMemberID,
			FullName:     res.Staff.FullName,
			Role:         string(res.Staff.Role),
			Designation:  string(res.Staff.Designation),
			PhoneNum:     res.Staff.PhoneNum,
			CreatedAt:    timestamppb.New(res.Staff.CreatedAt),
		},
	}, nil
}
