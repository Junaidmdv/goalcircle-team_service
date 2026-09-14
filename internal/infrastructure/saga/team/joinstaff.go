package teamsaga

import (
	"context"

	usrclient "github.com/Junaidmdv/goalcircle-protos/user/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/saga"
	staffuc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/staff"
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
	UserID       string
	Code         string
	Role         entity.UserRole
	RefreshToken string

	StaffRes   *staffuc.JoinStaffRes
	AddUserRes *AddUserRoleRes
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

func (tm *joinStaffSaga) RegisterStaff(
	ctx context.Context,
	state *JoinStaffSagaState,
) (*JoinStaffSagaState, error) {

	steps := []saga.SagaStep{


		{
			Name: "register staff",

			Action: func(
				ctx context.Context,
				sagaState interface{},
			) error {

				r := sagaState.(*JoinStaffSagaState)

				res, err := tm.staffUc.JoinStaff(
					ctx,
					&staffuc.JoinStaffReq{
						UserID: r.UserID,
						Code:   r.Code,
					},
				)

				if err != nil {
					return err
				}

				r.StaffRes = res

				return nil
			},

			Compensate: func(
				ctx context.Context,
				sagaState interface{},
			) error {

				r := sagaState.(*JoinStaffSagaState)

				if r.StaffRes == nil {
					return nil
				}

				return tm.staffUc.CompensateJoinStaff(
					ctx,
					staffuc.ComponsateJoinStaffReq{
						StaffID:      r.StaffRes.ID,
						TeamMemberID: r.StaffRes.TeamMemberID,
						InviteID:     r.StaffRes.InviteID,
					},
				)
			},
		},



		{
			Name: "update user role",

			Action: func(
				ctx context.Context,
				sagaState interface{},
			) error {

				r := sagaState.(*JoinStaffSagaState)

				res, err := tm.userclient.AddUserRole(
					ctx,
					&usrclient.AddUserRoleReq{
						UserId:       r.UserID,
						Role:         string(r.Role),
						RefreshToken: r.RefreshToken,
					},
				)

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
					RefreshTokenExpiry: res.RefreshTokenExpiry.AsTime(),
				}

				return nil
			},

			Compensate: func(
				ctx context.Context,
				sagaState interface{},
			) error {

				r := sagaState.(*JoinStaffSagaState)

				_, err := tm.userclient.RemoveUserRole(
					ctx,
					&usrclient.RemoveUserRoleReq{
						UserId:       r.UserID,
						RefreshToken: r.RefreshToken,
					},
				)

				return err
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
