package team

import (
	"time"

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
	PhoneNum    string
	Email       string
}

type UpdateTeamDetailsReq struct {
	UserID      string
	Name        *string
	City        *string
	Description *string
	ShortName   *string
	PhoneNum    *string
	Email       *string
}

type UpdateTeamDetailsRes struct {
	TeamID      string
	Name        string
	City        string
	Description string
	ShortName   string
	LogoKey     string
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
	TeamID string
}

type GetTeamRes struct {
	ID            string
	Name          string
	ShortName     string
	City          string
	LogoKey       string
	TeamCode      string
	Description   string
	Email         string
	PhoneNum      string
	TeamStatus    string
	PlayerCount   int32
	CaptainID     string
	ViceCaptainID string
	CreatedAt     time.Time
}

type UploadLogoReq struct {
	UserID      string
	LogoData    []byte
	ContentType string
	Size        int64
}

type UploadLogoRes struct {
	TeamID       string
	PresignedUrl string
}

type GetPresignedUrlReq struct {
	TeamID string
}

type GetPresignedUrlRes struct {
	TeamID       string
	PresignedUrl string
}
