package team

import (
	"bytes"
	"context"
	"io"
	"time"

	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	teamsaga "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/saga/team"
	team_uc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/team"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/imageutil"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/validater"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TeamHandler struct {
	teamUsecase team_uc.TeamUsecase
	pb.UnimplementedTeamServiceServer
	teamSaga  teamsaga.TeamSagaMaker
	timeOut   time.Duration
	validater *validater.Validater
}

func NewTeamHandler(tu team_uc.TeamUsecase, teamSaga teamsaga.TeamSagaMaker, timeout time.Duration, validater *validater.Validater) *TeamHandler {
	return &TeamHandler{
		teamUsecase: tu,
		teamSaga:    teamSaga,
		timeOut:     timeout,
		validater:   validater,
	}
}

func (th *TeamHandler) CreateTeam(ctx context.Context, req *pb.CreateTeamReq) (*pb.CreateTeamRes, error) {

	context, cancel := context.WithTimeout(ctx, th.timeOut)
	defer cancel()

	request := ToCreateTeam(req)

	if validationErrs := th.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := validater.ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := th.teamSaga.CreateTeamSaga(context, &teamsaga.TeamSagaState{
		UserID:       req.Owner.UserId,
		Role:         entity.TEAM,
		RefreshToken: req.RefreshToken,
		TeamName:     req.Name,
		City:         req.City,
		Description:  req.Description,
		FullName:     req.Owner.FullName,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.CreateTeamRes{
		TeamDetails: &pb.TeamDetailsRes{
			Id:        res.TeamRes.ID.String(),
			Name:      res.TeamRes.Name,
			ShortName: res.TeamRes.ShortName,
			City:      res.TeamRes.City,
			TeamCode:  res.TeamRes.Code,
		},
		User: &pb.UserRes{
			SessionId:          res.AddUserRes.SessionID,
			UserId:             res.AddUserRes.UserID,
			Email:              res.AddUserRes.Email,
			AccessToken:        res.AddUserRes.AccessToken,
			AccessTokenExpiry:  timestamppb.New(res.AddUserRes.AccessTokenExpiry),
			RefreshToken:       res.AddUserRes.RefreshToken,
			RefreshTokenExpiry: timestamppb.New(res.AddUserRes.RefreshTokenExpiry),
		},
		TeamMember: &pb.TeamMemberDetails{
			TeamMemberId: res.TeamMemberRes.TeamMemberID.String(),
			FullName:     res.TeamMemberRes.FullName,
			Role:         string(res.TeamMemberRes.Role),
		},
	}, nil
}

func (th *TeamHandler) UpdateTeam(ctx context.Context, req *pb.UpdateTeamReq) (*pb.UpdateTeamRes, error) {

	context, cancel := context.WithTimeout(ctx, th.timeOut)
	defer cancel()
	request := ToUpdateTeam(req)

	if validationErrs := th.validater.Validation(request); validationErrs != nil {
		stWithDetails, err := validater.ValidationError(validationErrs)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := th.teamUsecase.UpdateTeamDetails(context, &team_uc.UpdateTeamDetailsReq{
		TeamID:      req.TeamId,
		UserID:      req.UserId,
		Name:        req.Name,
		City:        req.City,
		Description: req.Description,
		ShortName:   req.ShortName,
		Email:       req.Email,
		PhoneNum:    req.PhoneNumber,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.UpdateTeamRes{
		TeamId:      res.TeamID,
		Name:        res.Name,
		City:        res.City,
		Description: res.Description,
		ShortName:   res.ShortName,
	}, nil
}

func (th *TeamHandler) SetCaption(ctx context.Context, req *pb.SetCaptainReq) (*pb.SetCaptainRes, error) {
	context, cancel := context.WithTimeout(ctx, th.timeOut)
	defer cancel()

	res, err := th.teamUsecase.ChangeCaptain(context, &team_uc.ChangeCaptainReq{
		UserID:   req.UserId,
		TeamID:   req.TeamId,
		PlayerID: req.PlayerId,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.SetCaptainRes{
		PlayerId: res.PlayerID.String(),
	}, nil

}

func (th *TeamHandler) SetViceCaptain(ctx context.Context, req *pb.SetViceCaptainReq) (*pb.SetViceCaptainRes, error) {
	context, cancel := context.WithTimeout(ctx, th.timeOut)
	defer cancel()

	res, err := th.teamUsecase.ChangeViceCaptain(context, &team_uc.ChangeViceCaptainReq{
		UserID:   req.UserId,
		TeamID:   req.TeamId,
		PlayerID: req.PlayerId,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.SetViceCaptainRes{
		PlayerId: res.PlayerID.String(),
	}, nil

}

func (th *TeamHandler) ListTeams(ctx context.Context, req *pb.ListTeamReq) (*pb.ListTeamRes, error) {
	context, cancel := context.WithTimeout(ctx, th.timeOut)
	defer cancel()

	status := protoTeamStatusToEntity[req.TeamStatus]

	res, err := th.teamUsecase.ListTeams(context, &team_uc.ListTeamsReq{
		Page:   req.Limit,
		Limit:  req.Limit,
		Status: status,
		Search: req.Search,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	var teams []*pb.TeamDetails

	for _, tms := range res.Teams {

		tmstatus := entityTeamStatusToProto[tms.TeamStatus]

		teams = append(teams, &pb.TeamDetails{
			TeamId:     tms.TeamID.String(),
			Name:       tms.Name,
			City:       tms.City,
			LogoUrl:    tms.LogoUrl,
			TeamCode:   tms.TeamCode,
			TeamStatus: tmstatus,
		})
	}

	return &pb.ListTeamRes{
		TeamDetails: teams,
		Pagination: &pb.PaginateTeam{
			TotalPage: res.Pagination.TotalPage,
			Page:      res.Pagination.Page,
			Limit:     res.Pagination.Limit,
			TotalItem: int64(res.Pagination.TotalItem),
		},
	}, nil

}

func (th *TeamHandler) RegisterTeamMember(ctx context.Context, req *pb.RegisterTeamMemberReq) (*pb.RegisterTeamMemberRes, error) {
	return &pb.RegisterTeamMemberRes{}, nil
}

func (th *TeamHandler) AddLogo(stream grpc.ClientStreamingServer[pb.AddLogoReq, pb.AddLogoRes]) error {

	ctx := stream.Context()

	var (
		meta      *pb.TeamLogoMetaData
		buffer    bytes.Buffer
		imagesize int64
	)

	ctx, cancel := context.WithTimeout(ctx, th.timeOut)
	defer cancel()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch data := req.Data.(type) {
		case *pb.AddLogoReq_MetaData:
			if meta != nil {
				return status.Error(codes.InvalidArgument, "logo meta data already send")
			}

			meta = data.MetaData

		case *pb.AddLogoReq_Chunk:
			if meta == nil {
				return status.Error(codes.InvalidArgument, "logo meta data is missing")
			}
			buffer.Write(data.Chunk)
			size := len(data.Chunk)
			imagesize += int64(size)

			if imagesize > 1<<entity.PictureMaxLen {
				return status.Error(codes.InvalidArgument, "logo size exceeds maximum size")
			}
		}

	}

	contentype, err := imageutil.ValidateImage(buffer.Bytes(), imageutil.TeamLogo)
	if err != nil {
		return apperror.GRPCStatus(err)
	}

	res, err := th.teamUsecase.UploadLogo(ctx, &team_uc.UploadLogoReq{
		TeamID:      meta.TeamId,
		LogoData:    buffer.Bytes(),
		ContentType: contentype,
		Size:        imagesize,
	})

	if err != nil {
		return apperror.GRPCStatus(err)
	}

	return stream.SendAndClose(&pb.AddLogoRes{
		TeamId:      res.TeamID,
		ResignedUrl: res.PresignedUrl,
	})

}


func(th *TeamHandler)GetPresignedURL(ctx context.Context,req *pb.GetPlayerPresignedUrlReq)(*pb.GetPlayerPresignedUrlRes,error){
	return nil,nil
}