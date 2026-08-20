package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UNSPECIFIED UserRole = "unspecified"
	ORGAINISER  UserRole = "organiser"
	TEAM        UserRole = "team"
)

type TeamMemberRole string

const (
	TeamMemberRoleStaff  TeamMemberRole = "STAFF"
	TeamMemberRolePlayer TeamMemberRole = "PLAYER"
)

type TeamMemberStatus string

const (
	TeamMemberStatusActive   TeamMemberStatus = "ACTIVE"
	TeamMemberStatusInactive TeamMemberStatus = "INACTIVE"
	TeamMemberStatusRelease  TeamMemberStatus = "RELEASED"
)

type TeamMember struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	TeamID     uuid.UUID
	UserID     string
	Role       TeamMemberRole
	Status     TeamMemberStatus
	JoinedAt   time.Time
	ReleasedAt time.Time
	ResignedAt time.Time
	Player     *Player     `gorm:"foreignKey:TeamMemberID"`
	Staff      *Staff      `gorm:"foreignKey:TeamMemberID"`
	TeamInvite *TeamInvite `gorm:"foreignKey:TeamMemberID"`
}
