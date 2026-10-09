package playerinviteuc

import (
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type SendPlayerInvitationReq struct {
	SendUserID       string
	PlayerID         string
	TeamID           string
	OfferedJerseyNum *int32
	Message          *string
}

type SendPlayerInvitationRes struct {
	RequestID string
	Status    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type TeamInviteListReq struct {
	TeamID string
	Page   int
	Limit  int
	Status entity.PlayerInvitationStatus
	Search *string
}

type PaginationDetails struct {
	TotalPage int32
	Page      int32
	Limit     int32
	Total     int64
}

type PlayerDetails struct {
	ID            string
	UserID        string
	FullName      string
	ImageUrl      string
	Position      entity.PlayerPosition
	PreferredFoot string
	Code          string
}

type InvitationDetails struct {
	RequestID    string
	InvitedBy    string
	Player       PlayerDetails
	JerseyNumber int32
	Status       entity.PlayerInvitationStatus
	Message      string
	RejectReason string
	CreatedAt    time.Time
	ExpiresAt    time.Time
	RespondedAt  time.Time
}

type TeamInviteListRes struct {
	Data     []InvitationDetails
	Paginate *PaginationDetails
}

type AcceptPlayerInviteReq struct {
	InviteID string
	UserID   string
}

type AcceptPlayerInviteRes struct {
	InviteID      uuid.UUID
	RequestStatus string
	RespondedAt   time.Time
	TeamPlayerID  uuid.UUID
	TeamID        uuid.UUID
}

type RejectPlayerInviteReq struct {
	InviteID string
	UserID   string
	Reason   string
}

type RejectPlayerInviteRes struct {
	InviteID    uuid.UUID
	RespondedAt time.Time
}

type PlayerInviteCancelReq struct {
	InviteID string
	UserID   string
	TeamID   string
}

type PlayerInviteCancelRes struct {
	InviteID    uuid.UUID
	RespondedAt time.Time
}

type PlayerInvitesReq struct {
	UserID string
}

type PlayerInvitationDetails struct {
	InviteID         uuid.UUID
	InvitedByUser    uuid.UUID
	TeamID           uuid.UUID
	OfferedJerseyNum *int32
	Message          *string
	CreatedAt        time.Time
}

type PlayerInvitesRes struct {
	UserID            uuid.UUID
	InvitationDetails []PlayerInvitationDetails
}
