package teamsaga

import (
	"context"
	"time"

	usrclient "github.com/Junaidmdv/goalcircle-protos/user/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/saga"
	staffuc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/staff"
	"github.com/Junaidmdv/goalcircle-team_service/internal/usecase/team"
	teammemberuc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
)

type TeamSagaMaker interface {
	CreateTeamSaga(context.Context, *TeamSagaState) (*TeamSagaState, error)
}

type teamSaga struct {
	teamUsecase  team.TeamUsecase
	staffUsecase staffuc.StaffUsecase
	userclient   usrclient.AuthServiceClient
	logger       logger.Logger
}

func NewTeamSagaMaker(teamuc team.TeamUsecase, tmuc teammemberuc.TeamMemberUsecase, usrclient usrclient.AuthServiceClient, stuc staffuc.StaffUsecase, logger logger.Logger) TeamSagaMaker {
	return &teamSaga{
		teamUsecase:  teamuc,
		staffUsecase: stuc,
		userclient:   usrclient,
		logger:       logger,
	}
}

type TeamSagaState struct {
	UserID       string
	Role         entity.UserRole
	RefreshToken string
	TeamName     string
	City         string
	Description  string
	FullName     string
	PhoneNum     string
	Email        string
	DOB          time.Time
	AddUserRes   *AddUserRoleRes
	TeamRes      *CreateTeamRes
	StaffRes     *staffuc.AddStaffRes
	// TeamMemberRes *TeamMemberRes
}

func (ts *teamSaga) CreateTeamSaga(ctx context.Context, req *TeamSagaState) (*TeamSagaState, error) {

	steps := []saga.SagaStep{

		{
			Name: "Create Team",
			Action: func(ctx context.Context, sagaState interface{}) error {

				req := sagaState.(*TeamSagaState)

				res, err := ts.teamUsecase.CreateTeam(ctx, &team.CreateTeamReq{
					UserID:      req.UserID,
					Name:        req.TeamName,
					City:        req.City,
					Description: req.Description,
					PhoneNum:    req.PhoneNum,
					Email:       req.Email,
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

				res, err := ts.staffUsecase.AddStaff(ctx, &staffuc.AddStaffReq{
					UserID:           req.UserID,
					TeamID:           req.TeamRes.ID.String(),
					FullName:         req.FullName,
					Role:             entity.StaffRoleManagement,
					Designation:      entity.StaffDesignationOwner,
					DOB:              req.DOB,
					TeamMemberStatus: entity.TeamMemberStatusActive,
				})

				if err != nil {
					return err
				}

				req.StaffRes = &staffuc.AddStaffRes{
					StaffID:      res.StaffID,
					TeamMemberID: res.TeamMemberID,
					FullName:     res.FullName,
					Age:          res.Age,
					Designation:  res.Designation,
					Role:         res.Role,
					PresignedUrl: res.PresignedUrl,
				}

				return nil
			},
			Compensate: func(ctx context.Context, sagaState interface{}) error {
				// req, _ := sagaState.(*TeamSagaState)

				// if err := ts.teamMemberUc.DeleteTeamOwner(ctx, &req.TeamMemberRes.TeamMemberID); err != nil {
				// 	return err
				// }

				return nil
			},
		},

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
	}

	teamOrchestrator := saga.NewOrchestrator(steps, ts.logger)
	if err := teamOrchestrator.Execute(ctx, req); err != nil {
		return nil, err
	}

	return req, nil

}
