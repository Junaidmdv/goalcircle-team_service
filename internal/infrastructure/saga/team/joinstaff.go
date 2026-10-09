package teamsaga

import (
	"context"
	"errors"

	usrclient "github.com/Junaidmdv/goalcircle-protos/user/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/saga"
	staffuc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/staff"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
)

type JoinStaffSaga interface {
	RegisterStaff(
		ctx context.Context,
		state *JoinStaffSagaState,
	) (*JoinStaffSagaState, error)
}

type joinStaffSaga struct {
	staffUc    staffuc.StaffUsecase
	userclient usrclient.AuthServiceClient
	logger     logger.Logger
}

type JoinStaffSagaState struct {
	JoinStaffReq   *staffuc.JoinStaffReq
	JoinStaffRes   *staffuc.JoinStaffRes
	AddUserRoleReq *AddUserRoleReq
	AddUserRoleRes *AddUserRoleRes
}

func NewJoinStaffSaga(
	staffUc staffuc.StaffUsecase,
	client usrclient.AuthServiceClient,
	logger logger.Logger,
) JoinStaffSaga {
	return &joinStaffSaga{
		staffUc:    staffUc,
		userclient: client,
		logger:     logger,
	}
}

func (tm *joinStaffSaga) RegisterStaff(ctx context.Context, state *JoinStaffSagaState) (*JoinStaffSagaState, error) {

	steps := []saga.SagaStep{

		{
			Name: "register staff",

			Action: func(ctx context.Context, sagaState interface{}) error {

				r := sagaState.(*JoinStaffSagaState)

				res, err := tm.staffUc.JoinStaff(ctx, state.JoinStaffReq)

				if err != nil {
					return err
				}

				r.JoinStaffRes = res

				return nil
			},

			Compensate: func(
				ctx context.Context,
				sagaState interface{},
			) error {

				r := sagaState.(*JoinStaffSagaState)

				if r.JoinStaffRes == nil {
					return nil
				}

				return tm.staffUc.CompensateJoinStaff(
					ctx,
					staffuc.ComponsateJoinStaffReq{
						StaffID:      r.JoinStaffRes.ID,
						TeamMemberID: r.JoinStaffRes.TeamMemberID,
						InviteID:     r.JoinStaffRes.InviteID,
					},
				)
			},
		},

		{
			Name: "update staff role",
			Action: func(ctx context.Context, sagaState interface{}) error {
				r := sagaState.(*JoinStaffSagaState)
				res, err := tm.userclient.AddUserRole(ctx, &usrclient.AddUserRoleReq{
					UserId:       r.AddUserRoleReq.UserID,
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

				r, ok := sagaState.(*JoinStaffSagaState)
				if !ok {
					return apperror.NewInternalError(apperror.InternalErrorMsg, errors.New("invalid saga data"))
				}

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

	orchestrator := saga.NewOrchestrator(
		steps,
		tm.logger,
	)

	if err := orchestrator.Execute(ctx, state); err != nil {
		return nil, err
	}

	return state, nil
}
