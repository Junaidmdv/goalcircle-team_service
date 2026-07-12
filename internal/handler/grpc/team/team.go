package team

import (
	"context"
	"time"

	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	teamsaga "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/saga/team"
	team_uc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/team"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TeamHandler struct {
	teamUsecase team_uc.TeamUsecase
	pb.UnimplementedTeamServiceServer
	teamSaga teamsaga.TeamSagaMaker
	timeOut  time.Duration
}

func NewTeamHandler(tu team_uc.TeamUsecase, teamSaga teamsaga.TeamSagaMaker, timeout time.Duration) *TeamHandler {
	return &TeamHandler{
		teamUsecase: tu,
		teamSaga:    teamSaga,
		timeOut:     timeout,
	}
}

func (th *TeamHandler) CreateTeam(ctx context.Context, req *pb.CreateTeamReq) (*pb.CreateTeamRes, error) {

	context, cancel := context.WithTimeout(ctx, th.timeOut)
	defer cancel()

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

func (th *TeamHandler) UpdateTeamDetails(ctx context.Context, req *pb.UpdateTeamReq) (*pb.UpdateTeamRes, error) {

	context, cancel := context.WithTimeout(ctx, th.timeOut)
	defer cancel()

	res, err := th.teamUsecase.UpdateTeamDetails(context, &team_uc.UpdateTeamDetailsReq{
		TeamID:       req.TeamId,
		TeamMemberID: req.TeamMemberId,
		Name:         req.Name,
		City:         req.City,
		Description:  req.Description,
		ShortName:    req.ShortName,
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

func (th *TeamHandler) UpdateTeamContactDetails(ctx context.Context, req *pb.UpdateTeamContactDetailsReq) (*pb.UpdateTeamContactDetailsRes, error) {
	context, cancel := context.WithTimeout(ctx, th.timeOut)
	defer cancel()

	res, err := th.teamUsecase.UpdateTeamContactDetails(context, &team_uc.UpdateTeamContactDetailsReq{
		TeamID:       req.TeamId,
		TeamMemberID: req.TeamMemberId,
		ContactEmail: &req.ContactEmail,
		ContactPhone: &req.ContactPhoneNum,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.UpdateTeamContactDetailsRes{
		TeamId:          res.TeamID.String(),
		ContactEmail:    res.ContactEmail,
		ContactPhoneNum: res.ContactPhone,
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

func (th TeamHandler) RegisterTeamMember(ctx context.Context, req *pb.RegisterTeamMemberReq) (*pb.RegisterTeamMemberRes, error) {
	return &pb.RegisterTeamMemberRes{}, nil
}
