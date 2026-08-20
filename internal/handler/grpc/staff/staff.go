package staff

import (
	"bytes"
	"context"
	"io"
	"time"

	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	teamv1 "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
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
	teamv1.UnimplementedStaffServiceServer
	StaffUsecase staffuc.StaffUsecase
	timeOut      time.Duration
	validater    validater.Validater
}

func NewStaffHandler(stafuc staffuc.StaffUsecase, timeout time.Duration) *StaffHandler {
	return &StaffHandler{
		StaffUsecase: stafuc,
		timeOut:      timeout,
	}

}

func (sh *StaffHandler) AddStaff(ctx context.Context, stream grpc.ClientStreamingServer[pb.AddStaffReq, pb.AddStaffRes]) error {

	c := stream.Context()

	ctx, cancel := context.WithTimeout(c, sh.timeOut)
	defer cancel()

	var (
		staffDetails *pb.StaffDetails
		buffer       bytes.Buffer
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
		case *pb.AddStaffReq_Meta:
			if staffDetails != nil {
				return status.Error(codes.InvalidArgument, "staff details already sent")
			}
			staffDetails = data.Meta

		case *pb.AddStaffReq_Chunks:
			if staffDetails == nil {
				return status.Error(codes.InvalidArgument, "player details is missing")
			}
			buffer.Write(data.Chunks)
		}
	}

	if err := imageutil.ValidateImage(buffer.Bytes(), imageutil.StaffImage); err != nil {
		return apperror.GRPCStatus(err)
	}

	request := ToCreateStaff(staffDetails)

	if validationErrs := sh.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := validater.ValidationError(validationErrs)
		if err != nil {
			return status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return stWithDetails.Err()
	}

	staff, err := sh.StaffUsecase.AddStaff(ctx, &staffuc.AddStaffReq{
		TeamID:           request.TeamID,
		UserID:           request.UserID,
		FullName:         request.FullName,
		Role:             request.Role,
		Designation:      request.Designation,
		DOB:              request.DOB,
		TeamMemberStatus: entity.TeamMemberStatusInactive,
		ImageBytes:       buffer.Bytes(),
	})

	if err != nil {
		return apperror.GRPCStatus(err)
	}

	return stream.SendAndClose(&pb.AddStaffRes{
		StaffId:      staff.StaffID.String(),
		TeamMemberId: staff.TeamMemberID.String(),
		FullName:     staff.FullName,
		Age:          staff.Age,
		Role:         string(staff.Role),
		Designatin:   string(staff.Designation),
	})

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
			ReleasedAt:   timestamppb.New(v.ReleasedAt),
			JoinedAt:     timestamppb.New(v.JoinedAt),
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
		DateOfBirth: *request.DateOfBirth,
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
		Age:          res.Age,
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
		DateOfBirth:  res.DateOfBirth,
		Role:         res.Role,
		Designation:  res.Designation,
		PhoneNum:     res.PhoneNum,
		CreatedAt:    timestamppb.New(res.CreatedAt),
		Status:       res.Status,
		JoinedAt:     res.JoinedAt,
		ReleasedAt:   &res.ReleasedAt,
	}, nil
}

func (sh *StaffHandler) UpdateStaffImage(ctx context.Context, stream grpc.ClientStreamingServer[pb.UpdateStaffImageReq, pb.UpdateStaffImageRes]) error {
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
