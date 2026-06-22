package teamsaga

import (
	"context"

	usrclient "github.com/Junaidmdv/goalcircle-protos/user/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/saga"
	"github.com/Junaidmdv/goalcircle-team_service/internal/usecase/team"
	"github.com/Junaidmdv/goalcircle-team_service/internal/usecase/teamowner"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
)

type TeamSagaMaker interface {
	CreateTeamSaga(context.Context, *TeamSagaState) (*TeamSagaState, error)
}

type teamSaga struct {
	teamUsecase      team.TeamUsecase
	teamOwnerUsecase teamowner.TeamOwnerUsecase
	userclient       usrclient.AuthServiceClient
	logger           logger.Logger
}

func NewTeamSagaMaker(teamuc team.TeamUsecase, teamMemberUsecase teamowner.TeamOwnerUsecase, usrclient usrclient.AuthServiceClient, logger logger.Logger) TeamSagaMaker {
	return &teamSaga{
		teamUsecase:      teamuc,
		teamOwnerUsecase: teamMemberUsecase,
		userclient:       usrclient,
		logger:           logger,
	}
}

type TeamSagaState struct {
	UserID        string
	Role          entity.UserRole
	RefreshToken  string
	TeamName      string
	City          string
	Description   string
	FullName      string
	AddUserRes    *AddUserRoleRes
	TeamRes       *CreateTeamRes
	TeamMemberRes *TeamMemberRes
}

func (ts *teamSaga) CreateTeamSaga(ctx context.Context, req *TeamSagaState) (*TeamSagaState, error) {

	steps := []saga.SagaStep{
		{
			Name: "update role",
			Action: func(ctx context.Context, sagaState interface{}) error {
				r := sagaState.(*TeamSagaState)
				res, err := ts.userclient.AddUserRole(ctx, &usrclient.AddUserRoleReq{
					UserId:       r.UserID,
					Role:         string(r.Role),
					RefreshToken: r.RefreshToken,
				})

				if err != nil {
					return err
				}
				r.AddUserRes = &AddUserRoleRes{
					SessionID:          res.SessionId,
					UserID:             res.UserId,
					Email:              res.Email,
					AccessToken:        res.AccessToken,
					AccessTokenExpiry:  res.AccessTokenExpiry.AsTime(),
					RefreshToken:       res.RefreshToken,
					RefreshTokenExpiry: res.AccessTokenExpiry.AsTime(),
				}

				return nil
			},
			Compensate: func(ctx context.Context, sagaState interface{}) error {

				r := sagaState.(*TeamSagaState)

				_, err := ts.userclient.RemoveUserRole(ctx, &usrclient.RemoveUserRoleReq{
					UserId:       r.UserID,
					RefreshToken: r.RefreshToken,
				})

				if err != nil {
					return err
				}

				return nil
			},
		},

		{
			Name: "Create Team",
			Action: func(ctx context.Context, sagaState interface{}) error {

				req := sagaState.(*TeamSagaState)

				res, err := ts.teamUsecase.CreateTeam(ctx, &team.CreateTeamReq{
					Name:        req.TeamName,
					City:        req.City,
					Description: req.Description,
				})

				if err != nil {
					return err
				}

				req.TeamRes = &CreateTeamRes{
					ID:        res.TeamID,
					Name:      res.Name,
					ShortName: res.ShortName,
					City:      res.City,
					Code:      res.Code,
				}

				return nil
			},
			Compensate: func(ctx context.Context, req interface{}) error {
				request := req.(*TeamSagaState)

				err := ts.teamUsecase.DeleteTeam(ctx, request.TeamRes.ID)
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name: "add team owner",
			Action: func(ctx context.Context, sagaState interface{}) error {

				req, _ := sagaState.(*TeamSagaState)

				// res, err := ts.TeamOwnerUsecase.AddTeamMember(
				// })

				res, err := ts.teamOwnerUsecase.AddTeamOwner(ctx, &teamowner.AddTeamOwnerReq{
					TeamID:   req.TeamRes.ID,
					UserId:   req.UserID,
					FullName: req.FullName,
					Role:     entity.OWNER,
				},
				)

				if err != nil {
					return err
				}

				req.TeamMemberRes = &TeamMemberRes{
					TeamMemberID: res.TeamMemberID,
					TeamID:       req.TeamRes.ID,
					UserID:       req.UserID,
					FullName:     req.FullName,
					Role:         req.Role,
				}

				return nil
			},
			Compensate: func(ctx context.Context, sagaState interface{}) error {
				req, _ := sagaState.(*TeamSagaState)

				if err := ts.teamOwnerUsecase.DeleteTeamOwner(ctx, &req.TeamMemberRes.TeamMemberID); err != nil {
					return err
				}
				return nil
			},
		},
	}

	teamOrchestrator := saga.NewOrchestrator(steps, ts.logger)
	if err := teamOrchestrator.Execute(ctx, req); err != nil {
		return nil, err
	}

	return req, nil

}
