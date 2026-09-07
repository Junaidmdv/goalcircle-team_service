package staffuc

import (
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/google/uuid"
)

type AddStaffReq struct {
	UserID           string
	TeamID           string
	FullName         string
	PhoneNum         *string
	Role             entity.StaffRole
	Designation      entity.StaffDesignation
	TeamMemberStatus entity.TeamMemberStatus
	ImageBytes       []byte
}

type AddStaffRes struct {
	StaffID      uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	Age          int32
	Designation  entity.StaffDesignation
	Role         entity.StaffRole
	PresignedUrl string
}

type UploadStaffImageReq struct {
	UserID    string
	TeamID    string
	StaffID   string
	ImageData []byte
}

type UploadStaffImageRes struct {
	StaffID      string
	PresignedUrl string
}

type StaffImagePresignedUrlReq struct {
	TeamID  string
	StaffID string
}

type StaffImagePresignedUrlRes struct {
	StaffID      string
	PresignedUrl string
}

type RemoveStaffImageReq struct {
	UserID  string
	TeamID  string
	StaffID string
}

type RemoveStaffImageRes struct {
	Success bool
}

type ListTeamStaffReq struct {
	TeamID      string
	Page        int32
	Limit       int32
	Role        entity.StaffRole
	Designation entity.StaffDesignation
	Search      string
}

type StaffDetails struct {
	StaffID      uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	Role         string
	Designation  string
	PresignedUrl string
	Status       string
}

type PaginationDetails struct {
	TotalPage int32
	Page      int32
	Limit     int32
	Total     int64
}

type ListTeamStaffRes struct {
	Staff      []StaffDetails
	Pagination *PaginationDetails
}

type UpdateStaffReq struct {
	UserID      string
	TeamID      string
	StaffID     string
	FullName    string
	Role        entity.StaffRole
	Designation entity.StaffDesignation
	PhoneNum    string
}

type UpdateStaffRes struct {
	StaffID      uuid.UUID
	TeamMemberID uuid.UUID
	FullName     string
	DateOfBirth  time.Time
	Role         string
	Designation  string
	PhoneNum     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ReleaseStaffReq struct {
	TeamID  string
	UserID  string
	StaffID string
}

type ReleaseStaffRes struct {
	StaffID      string
	TeamMemberID string
	FullName     string
	Role         string
	Designation  string
	Status       string
	JoinedAt     string
	ReleasedAt   string
}

type TransferOwnershipReq struct {
}

type TransferOwnershipRes struct {
}

type GetStaffReq struct {
	StaffID string
	TeamID  string
}

type GetStaffRes struct {
	ID           string
	TeamMemberID string
	FullName     string
	Role         string
	Designation  string
	PhoneNum     string
	CreatedAt    time.Time
	Status       string
	JoinedAt     string
	ReleasedAt   string
}

type RegisterTeamMemberReq struct {
	UserID string
	Code   string
}

type RegisterTeamMemberRes struct {
}

type GetStaffProfileReq struct {
	UserID string
}

type StaffDetail struct {
	TeamMemberID string
	StaffID      string
	UserID       string
	TeamID       string
	FullName     string
	Role         entity.StaffRole
	Designation  entity.StaffDesignation
	PhoneNum     string
	CreatedAt    time.Time
}

type TeamDetail struct {
	TeamID    string
	Name      string
	ShortName string
	City      string
	TeamCode  string
}

type GetStaffProfileRes struct {
	Team  *TeamDetail
	Staff *StaffDetail
}
