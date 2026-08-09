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
}

type teamMemberSaga struct {
	teamUsecase  team.TeamUsecase
	teamMemberUc teammemberuc.TeamMemberUsecase
	userclient   usrclient.AuthServiceClient
	logger       logger.Logger
}

type TeamMemberSagaState struct {
	UserID       string
	Code         string
	Role         entity.UserRole
	RefreshToken string
	AddUserRes   *AddUserRoleRes
}

func NewTeamMemberSaga() TeamMemberSaga {
	return &teamMemberSaga{}
}

func (tm *teamMemberSaga) RegisterTeamMember(ctx context.Context, req *TeamMemberSagaState) (*TeamMemberSagaState, error) {
	steps := []saga.SagaStep{
		{
			Name: "team member",
			Action: func(ctx context.Context, sagaState interface{}) error {

				r := sagaState.(*TeamMemberSagaState)

				_,err:=tm.teamMemberUc.RegisterTeamMember(ctx, &teammemberuc.RegisterTeamMemberReq{
					UserID: r.UserID,
					Code:   r.Code,
				})

				if err != nil{
					return err 
				}


				return nil
			},
			Compensate: func(ctx context.Context, sagaState interface{}) error {
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
