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

type TeamMemberStatus string

const (
	TeamMemberStatusActive   TeamMemberStatus = "ACTIVE"
	TeamMemberStatusInactive TeamMemberStatus = "INACTIVE"
	TeamMemberStatusRelease  TeamMemberStatus = "RELEASED"
)

type TeamMemberRole string

const (
	TeamMemberStaff  TeamMemberRole = "STAFF"
	TeamMemberPlayer TeamMemberRole = "PLAYER"
)

type TeamMember struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	TeamID     uuid.UUID
	UserID     uuid.UUID
	Status     TeamMemberStatus
	JoinedAt   time.Time
	ReleasedAt time.Time
	LeftAt     time.Time
	Role       TeamMemberRole
	Player     *Player     `gorm:"foreignKey:TeamMemberID"`
	Staff      *Staff      `gorm:"foreignKey:TeamMemberID"`
	TeamInvite *TeamInvite `gorm:"foreignKey:TeamMemberID"`
	Team       Team        `gorm:"foreignKey:TeamID;references:ID;"`
}
