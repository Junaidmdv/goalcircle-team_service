package playerinvitehandler

import (
	teampb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
)

func toSendPlayerInviteReq(input *teampb.SendPlayerInviteReq) *SendPlayerInvitationReq {
	return &SendPlayerInvitationReq{
		PlayerID:         input.PlayerId,
		UserID:           input.SenderUserId,
		TeamID:           input.TeamId,
		Message:          input.Message,
		OfferedJerseyNum: input.OfferedJerseyNum,
	}
}

func mapTeamInviteStatus(status teampb.PlayerInvitationStatus) entity.PlayerInvitationStatus {

	switch status {
	case teampb.PlayerInvitationStatus_PLAYER_INVITATION_STATUS_ACCEPTED:
		return entity.InvitationAccepted
	case teampb.PlayerInvitationStatus_PLAYER_INVITATION_STATUS_CANCELLED:
		return entity.InvitationCancelled
	case teampb.PlayerInvitationStatus_PLAYER_INVITATION_STATUS_EXPIRED:
		return entity.InvitationExpired
	case teampb.PlayerInvitationStatus_PLAYER_INVITATION_STATUS_PENDING:
		return entity.InvitationPending
	case teampb.PlayerInvitationStatus_PLAYER_INVITATION_STATUS_UNSPECIFIED:
		return entity.InvitationInvalid
	}

	return entity.InvitationInvalid
}

func toListTeamInviteReq(input *teampb.TeamInviteListReq) *ListTeamInviteReq {
	status := mapTeamInviteStatus(input.Status)

	return &ListTeamInviteReq{
		TeamID: input.TeamId,
		Page:   int(input.Page),
		Limit:  int(input.Limit),
		Status: status,
		Search: *input.Search,
	}
}

func toAcceptTeamInviteReq(input *teampb.AcceptPlayerInviteReq) *AcceptPlayerInviteReq {
	return &AcceptPlayerInviteReq{
		InviteID: input.InviteId,
		UserID:   input.UserId,
	}
}

func toRejectTeamInvite(input *teampb.RejectPlayerInviteReq) *RejectPlayerInviteReq {
	return &RejectPlayerInviteReq{
		InviteID: input.InviteId,
		UserID:   input.UserId,
		Reason:   input.Reason,
	}
}

func toCancelPlayerInvite(input *teampb.CancelTeamInviteReq) *CancelPlayerInviteReq {
	return &CancelPlayerInviteReq{
		InviteID: input.InviteId,
		UserID:   input.UserId,
		TeamID:   input.TeamId,
	}
}

func toPlayerInvitesReq(input *teampb.PlayerInvitesReq) *PlayerInviteReq {
	return &PlayerInviteReq{
		UserID: input.UserId,
	}
}
