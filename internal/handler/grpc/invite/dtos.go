package playerinvitehandler

import "github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"

type SendPlayerInvitationReq struct {
	PlayerID         string  `json:"player_id" validate:"required"`
	TeamID           string  `json:"team_id" validate:"required"`
	UserID           string  `json:"user_id" validate:"required"`
	OfferedJerseyNum *int32  `json:"offered_jersey_number" validate:"omitempty"`
	Message          *string `json:"message" validate:"omitempty,max=200,min=50" `
}

type ListTeamInviteReq struct {
	TeamID string                        `json:"team_id" validate:"required"`
	Page   int                           `json:"page" validate:"omitempty,min=1"`
	Limit  int                           `json:"limit" validate:"omitempty,min=1"`
	Status entity.PlayerInvitationStatus `json:"status" validate:"omitempty"`
	Search string                        `json:"search" validate:"omitempty"`
}

type AcceptPlayerInviteReq struct {
	InviteID     string `json:"invite_id" validate:"required"`
	UserID       string `json:"user_id" validate:"required"`
	JerseyNumber int32  `json:"jersey_number" validate:"required,min=1,max=99"`
}

type RejectPlayerInviteReq struct {
	InviteID string `json:"invite_id" validate:"required"`
	UserID   string `json:"user_id" validate:"required"`
	Reason   string `json:"reject_reason" validate:"omitempty,min=50,max=200"`
}

type CancelPlayerInviteReq struct {
	InviteID string `json:"invite_id" validate:"required"`
	UserID   string `json:"user_id" validate:"required"`
	TeamID   string `json:"team_id" validate:"required"`
}

type PlayerInviteReq struct {
	UserID string `json:"user_id" validate:"required"`
}
