package teamsaga

import (
	"context"

	usrclient "github.com/Junaidmdv/goalcircle-protos/user/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/saga"
	playeruc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/player"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
)

type PlayerSaga interface {
	AddPlayer(ctx context.Context, req *PlayerSagaState) (*PlayerSagaState, error)
}

type AddPlayerSaga struct {
	playerUc   playeruc.PlayerUsecase
	userclient usrclient.AuthServiceClient
	logger     logger.Logger
}

type PlayerSagaState struct {
	AddPlayerReq   *playeruc.AddPlayerReq
	AddPlayerRes   *playeruc.AddPlayerRes
	AddUserRoleReq *AddUserRoleReq
	AddUserRoleRes *AddUserRoleRes
}

func NewPlayerSaga(playeruc playeruc.PlayerUsecase, client usrclient.AuthServiceClient, logger logger.Logger) PlayerSaga {
	return &AddPlayerSaga{
		playerUc:   playeruc,
		userclient: client,
		logger:     logger,
	}
}

func (tm *AddPlayerSaga) AddPlayer(ctx context.Context, req *PlayerSagaState) (*PlayerSagaState, error) {
	steps := []saga.SagaStep{
		{
			Name: "add new player",
			Action: func(ctx context.Context, sagaState interface{}) error {

				r := sagaState.(*PlayerSagaState)

				res, err := tm.playerUc.AddNewPlayer(ctx, &playeruc.AddPlayerReq{
					UserID:        r.AddPlayerReq.UserID,
					FullName:      r.AddPlayerReq.FullName,
					DOB:           r.AddPlayerReq.DOB,
					Nationality:   r.AddPlayerReq.Nationality,
					PreferredFoot: r.AddPlayerReq.PreferredFoot,
					Position:      r.AddPlayerReq.Position,
					ImageBytes:    r.AddPlayerReq.ImageBytes,
				})

				if err != nil {
					return err
				}

				r.AddPlayerRes = res

				return nil
			},
			Compensate: func(ctx context.Context, sagaState interface{}) error {

				return nil
			},
		},
		{
			Name: "update player role",
			Action: func(ctx context.Context, sagaState interface{}) error {
				r := sagaState.(*PlayerSagaState)
				res, err := tm.userclient.AddUserRole(ctx, &usrclient.AddUserRoleReq{
					UserId:       r.AddPlayerReq.UserID,
					Role:         r.AddUserRoleReq.Role,
					RefreshToken: r.AddUserRoleReq.RefreshToken,
				})

				if err != nil {
					return err
				}
				r.AddUserRoleRes = &AddUserRoleRes{
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

				r := sagaState.(*PlayerSagaState)

				_, err := tm.userclient.RemoveUserRole(ctx, &usrclient.RemoveUserRoleReq{
					UserId:       r.AddUserRoleReq.UserID,
					RefreshToken: r.AddUserRoleReq.RefreshToken,
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
