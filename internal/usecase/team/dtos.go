package team

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type CreateTeamRes struct {
	TeamID    uuid.UUID `gorm:"primaryKey"`
	Name      string
	ShortName string
	City      string
	Code      string
}

type CreateTeamReq struct {
	Name        string
	City        string
	Description string
}

type UpdateTeamDetailsReq struct {
	TeamID       string
	TeamMemberID string
	Name         *string
	City         *string
	Description  *string
	ShortName    *string
	PhoneNum     *string
	Email        *string
}

type UpdateTeamDetailsRes struct {
	TeamID      string
	Name        string
	City        string
	Description string
	ShortName   string
}

type UpdateTeamContactDetailsReq struct {
	TeamID       string
	TeamMemberID string
	ContactEmail *string
	ContactPhone *string
}

type UpdateTeamContactDetailsRes struct {
	TeamID       uuid.UUID
	ContactPhone string
	ContactEmail string
}

type ChangeCaptainReq struct {
	UserID   string
	TeamID   string
	PlayerID string
}

type ChangeCaptainRes struct {
	PlayerID uuid.UUID
}

type ChangeViceCaptainReq struct {
	UserID   string
	TeamID   string
	PlayerID string
}
type ChangeViceCaptainRes struct {
	PlayerID uuid.UUID
}

type ListTeamsReq struct {
	Page   int32
	Limit  int32
	Status entity.TeamStatus
	City   string
	Search string
}

type TeamData struct {
	TeamID     uuid.UUID
	Name       string
	City       string
	LogoUrl    string
	TeamCode   string
	TeamStatus entity.TeamStatus
}

type PaginateDetails struct {
	TotalPage int32
	Page      int32
	Limit     int32
	TotalItem int
}

type ListTeamsRes struct {
	Teams      []TeamData
	Pagination *PaginateDetails
}

type GetTeamReq struct {
}

type GetTeamRes struct {
}

type UploadLogoReq struct {
}

type UploadLogoRes struct {
}
