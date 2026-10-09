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
	PLAYER      UserRole = "player"
)

type TeamMemberStatus string

const (
	TeamMemberStatusActive TeamMemberStatus = "active"
	TeamMemberStatusLeft   TeamMemberStatus = "left"
)

type TeamMemberRole string

const (
	TeamMemberStaff  TeamMemberRole = "STAFF"
	TeamMemberPlayer TeamMemberRole = "PLAYER"
)

type TeamMember struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	TeamID     uuid.UUID `gorm:"type:uuid;not null"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	Role       TeamMemberRole
	Status     TeamMemberStatus
	Team       Team        `gorm:"foreignKey:TeamID;references:ID"`
	TeamPlayer *TeamPlayer `gorm:"foreignKey:TeamMemberID;references:ID"`
	Staff      *Staff      `gorm:"foreignKey:TeamMemberID;references:ID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
