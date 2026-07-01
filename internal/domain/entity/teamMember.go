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
	TeamMemberRoleOwner   TeamMemberRole = "OWNER"
	TeamMemberRoleManager TeamMemberRole = "MANAGER"
	TeamMemberRoleStaff   TeamMemberRole = "STAFF"
	TeamMemberRolePlayer  TeamMemberRole = "PLAYER"
)

type TeamMember struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	TeamID    uuid.UUID
	UserID    string
	FullName  string
	Role      TeamMemberRole
	CreatedAt time.Time
}
