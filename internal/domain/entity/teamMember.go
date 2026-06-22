package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UNSPECIFIED UserRole = "unspecified"
	ORGAINISER  UserRole = "organiser"
	MANAGER     UserRole = "manager"
	ADMIN       UserRole = "admin"
	OWNER       UserRole = "owner"
	COACH       UserRole = "coach"
	STAFF       UserRole = "staff"
	PLAYER      UserRole = "player"
)

type TeamMember struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	TeamID    uuid.UUID
	UserID    string
	FullName  string
	Role      UserRole
	CreatedAt time.Time
}
