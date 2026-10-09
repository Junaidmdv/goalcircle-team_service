package playerinvitehandler

import (
	"context"
	"time"

	teampb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	playerinviteuc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/invite"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/validater"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PlayerInviteHandler struct {
	validater      *validater.Validater
	playerinviteUc playerinviteuc.PlayerInviteUsecase
	TimeOut        time.Duration
	teampb.UnimplementedPlayerInviteServiceServer
}

func NewPlayerInviteHandler(validater *validater.Validater,
	playerinviteruc playerinviteuc.PlayerInviteUsecase,
	timeOut time.Duration,

) *PlayerInviteHandler {
	return &PlayerInviteHandler{
		validater:      validater,
		playerinviteUc: playerinviteruc,
		TimeOut:        timeOut,
	}
}

func (h *PlayerInviteHandler) SendPlayerInvitation(ctx context.Context, input *teampb.SendPlayerInviteReq) (*teampb.SendPlayerInviteRes, error) {

	c, cancel := context.WithTimeout(ctx, h.TimeOut)
	defer cancel()

	data := toSendPlayerInviteReq(input)

	if validationError := h.validater.Validation(data); validationError != nil {
		stWithDetails, err := validater.ValidationError(validationError)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := h.playerinviteUc.SendPlayerRequest(c, &playerinviteuc.SendPlayerInvitationReq{
		SendUserID:       data.UserID,
		TeamID:           data.TeamID,
		PlayerID:         data.PlayerID,
		OfferedJerseyNum: data.OfferedJerseyNum,
		Message:          data.Message,
	})
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &teampb.SendPlayerInviteRes{
		InviteId:   res.RequestID,
		Status:     res.Status,
		CreatedAt:  timestamppb.New(res.CreatedAt),
		ExpirestAt: timestamppb.New(res.ExpiresAt),
	}, nil
}

func (h PlayerInviteHandler) TeamInviteList(ctx context.Context, input *teampb.TeamInviteListReq) (*teampb.TeamInviteListRes, error) {

	c, cancel := context.WithTimeout(ctx, h.TimeOut)
	defer cancel()
	data := toListTeamInviteReq(input)

	if data.Status == entity.InvitationInvalid {
		err := apperror.NewBadRequestError("invalid player invite status")
		return nil, apperror.GRPCStatus(err)
	}

	if validationError := h.validater.Validation(data); validationError != nil {
		stWithDetails, err := validater.ValidationError(validationError)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := h.playerinviteUc.TeamInviteList(c, &playerinviteuc.TeamInviteListReq{
		TeamID: data.TeamID,
		Page:   data.Page,
		Limit:  data.Limit,
		Status: data.Status,
		Search: &data.Search,
	})
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	var inviteRes []*teampb.TeamInviteListRes_InvitationDetails

	for _, invite := range res.Data {
		inviteRes = append(inviteRes, &teampb.TeamInviteListRes_InvitationDetails{
			InviteID:  invite.RequestID,
			InvitedBy: invite.InvitedBy,
			Player: &teampb.TeamInviteListRes_PlayerDetail{
				PlayerId:      invite.Player.ID,
				UserId:        invite.Player.UserID,
				Position:      string(invite.Player.Position),
				PreferredFoot: string(invite.Player.Position),
				Code:          invite.Player.Code,
			},
			JerseyNumber:  invite.JerseyNumber,
			Status:        string(invite.Status),
			InviteMessage: invite.Message,
			RejectReason:  invite.RejectReason,
			CreatedAt:     timestamppb.New(invite.CreatedAt),
			ExpiresAt:     timestamppb.New(invite.ExpiresAt),
		})
	}

	return &teampb.TeamInviteListRes{
		InvitationDetails: inviteRes,
		Paginate: &teampb.TeamInviteListRes_Paginate{
			TotalPage: res.Paginate.TotalPage,
			Page:      res.Paginate.Page,
			Limit:     res.Paginate.Limit,
			Total:     res.Paginate.Total,
		},
	}, nil
}

func (h *PlayerInviteHandler) AcceptPlayerInvite(ctx context.Context, input *teampb.AcceptPlayerInviteReq) (*teampb.AcceptPlayerInviteRes, error) {

	c, cancel := context.WithTimeout(ctx, h.TimeOut)
	defer cancel()
	data := toAcceptTeamInviteReq(input)

	if validationError := h.validater.Validation(data); validationError != nil {
		stWithDetails, err := validater.ValidationError(validationError)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := h.playerinviteUc.AcceptPlayerInvite(c, &playerinviteuc.AcceptPlayerInviteReq{
		InviteID: data.InviteID,
		UserID:   data.UserID,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &teampb.AcceptPlayerInviteRes{
		InviteId:    res.InviteID.String(),
		Status:      res.RequestStatus,
		RespondedAt: timestamppb.New(res.RespondedAt),
		TeamId:      res.TeamPlayerID.String(),
	}, nil
}

func (h *PlayerInviteHandler) RejectPlayerInvite(ctx context.Context, input *teampb.RejectPlayerInviteReq) (*teampb.RejectPlayerInviteRes, error) {
	c, cancel := context.WithTimeout(ctx, h.TimeOut)
	defer cancel()
	data := toRejectTeamInvite(input)

	if validationError := h.validater.Validation(data); validationError != nil {
		stWithDetails, err := validater.ValidationError(validationError)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := h.playerinviteUc.RejectPlayerInvite(c, &playerinviteuc.RejectPlayerInviteReq{
		InviteID: data.InviteID,
		UserID:   data.UserID,
		Reason:   data.Reason,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &teampb.RejectPlayerInviteRes{
		InvitedId:   res.InviteID.String(),
		RespondedAt: timestamppb.New(res.RespondedAt),
	}, nil
}

func (h *PlayerInviteHandler) CancelPlayerInvite(ctx context.Context, input *teampb.CancelTeamInviteReq) (*teampb.CancelTeamInviteRes, error) {
	c, cancel := context.WithTimeout(ctx, h.TimeOut)
	defer cancel()
	data := toCancelPlayerInvite(input)

	if validationError := h.validater.Validation(data); validationError != nil {
		stWithDetails, err := validater.ValidationError(validationError)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := h.playerinviteUc.PlayerInviteCancel(c, &playerinviteuc.PlayerInviteCancelReq{
		InviteID: data.InviteID,
		UserID:   data.UserID,
		TeamID:   data.TeamID,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &teampb.CancelTeamInviteRes{
		InviteId:    res.InviteID.String(),
		RespondedAt: timestamppb.New(res.RespondedAt),
	}, nil
}

func (h *PlayerInviteHandler) PlayerInvites(ctx context.Context, input *teampb.PlayerInvitesReq) (*teampb.PlayerInvitesRes, error) {
	c, cancel := context.WithTimeout(ctx, h.TimeOut)
	defer cancel()
	data := toPlayerInvitesReq(input)

	if validationError := h.validater.Validation(data); validationError != nil {
		stWithDetails, err := validater.ValidationError(validationError)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := h.playerinviteUc.PlayerInvites(c, &playerinviteuc.PlayerInvitesReq{
		UserID: data.UserID,
	})
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	var invites []*teampb.PlayerInvitesRes_InviteDetails

	for _, invite := range res.InvitationDetails {
		invites = append(invites, &teampb.PlayerInvitesRes_InviteDetails{
			InviteId:             invite.InviteID.String(),
			InvitedByUser:        invite.InviteID.String(),
			TeamId:               invite.TeamID.String(),
			OfferredJerseyNumber: *invite.OfferedJerseyNum,
			Message:              *invite.Message,
			CreatedAt:            timestamppb.New(invite.CreatedAt),
		})
	}

	return &teampb.PlayerInvitesRes{
		UserId:        res.UserID.String(),
		InviteDetails: invites,
	}, nil
}
