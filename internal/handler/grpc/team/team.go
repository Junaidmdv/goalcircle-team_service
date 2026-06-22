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
	tu team_uc.TeamUsecase
	pb.UnimplementedTeamServiceServer
	teamSaga teamsaga.TeamSagaMaker
	timeOut  time.Duration
}

func NewTeamHandler(tu team_uc.TeamUsecase, teamSaga teamsaga.TeamSagaMaker, timeout time.Duration) *TeamHandler {
	return &TeamHandler{
		tu:       tu,
		teamSaga: teamSaga,
		timeOut:  timeout,
	}
}

func (th *TeamHandler) CreateTeam(ctx context.Context, req *pb.CreateTeamReq) (*pb.CreateTeamRes, error) {

	context, cancel := context.WithTimeout(ctx, th.timeOut)
	defer cancel()

	res, err := th.teamSaga.CreateTeamSaga(context, &teamsaga.TeamSagaState{
		UserID:       req.Owner.UserId,
		Role:         entity.OWNER,
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
	return nil, nil
}
