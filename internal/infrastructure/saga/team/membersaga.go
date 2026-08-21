package teamsaga

import (
	"context"

	usrclient "github.com/Junaidmdv/goalcircle-protos/user/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/saga"
	"github.com/Junaidmdv/goalcircle-team_service/internal/usecase/team"
	teammemberuc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
)

type TeamMemberSaga interface {
	RegisterTeamMember(ctx context.Context, req *TeamMemberSagaState) (*TeamMemberSagaState, error)
}

type teamMemberSaga struct {
	teamUsecase  team.TeamUsecase
	userclient   usrclient.AuthServiceClient
	logger       logger.Logger
	teamMemberUc teammemberuc.TeamMemberUsecase
}

type TeamMemberSagaState struct {
	UserID        string
	Code          string
	Role          entity.UserRole
	RefreshToken  string
	AddUserRes    *AddUserRoleRes
	TeamMemberRes *teammemberuc.RegisterTeamMemberRes
}

func NewTeamMemberSaga(tm team.TeamUsecase, teamMemberUc teammemberuc.TeamMemberUsecase, client usrclient.AuthServiceClient, logger logger.Logger) TeamMemberSaga {
	return &teamMemberSaga{
		teamUsecase:  tm,
		teamMemberUc: teamMemberUc,
		userclient:   client,
		logger:       logger,
	}
}

func (tm *teamMemberSaga) RegisterTeamMember(ctx context.Context, req *TeamMemberSagaState) (*TeamMemberSagaState, error) {
	steps := []saga.SagaStep{
		{
			Name: "register team member",
			Action: func(ctx context.Context, sagaState interface{}) error {

				r := sagaState.(*TeamMemberSagaState)

				res, err := tm.teamMemberUc.RegisterTeamMember(ctx, &teammemberuc.RegisterTeamMemberReq{
					UserID: r.UserID,
					Code:   r.Code,
				})

				r.TeamMemberRes = &teammemberuc.RegisterTeamMemberRes{
					InvitationID: res.InvitationID,
					TeamMemberID: res.TeamMemberID,
					TeamID:       res.TeamID,
					UserID:       res.UserID,
					Role:         res.Role,
					Status:       res.Status,
					JoinedAt:     res.JoinedAt,
				}

				if err != nil {
					return err
				}

				return nil
			},
			Compensate: func(ctx context.Context, sagaState interface{}) error {

				r := sagaState.(*TeamMemberSagaState)

				err := tm.teamMemberUc.CompensateRegisterTeamMember(ctx, &teammemberuc.CompensateRegisterTeamMemberReq{
					TeamMemberID: r.TeamMemberRes.TeamMemberID,
					InvitationID: r.TeamMemberRes.InvitationID, 
					UserID: r.TeamMemberRes.UserID,
				})
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name: "update role",
			Action: func(ctx context.Context, sagaState interface{}) error {
				r := sagaState.(*TeamMemberSagaState)
				res, err := tm.userclient.AddUserRole(ctx, &usrclient.AddUserRoleReq{
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

				r := sagaState.(*TeamMemberSagaState)

				_, err := tm.userclient.RemoveUserRole(ctx, &usrclient.RemoveUserRoleReq{
					UserId:       r.UserID,
					RefreshToken: r.RefreshToken,
				})

				if err != nil {
					return err
				}

				return nil
			},
		},
	}
	teamOrchestrator := saga.NewOrchestrator(steps, tm.logger)
	if err := teamOrchestrator.Execute(ctx, req); err != nil {
		return nil, err
	}

	return req, nil

}
