package entity

import "time"

type TeamMember struct {
	ID        string
	UUID      string `gorm:"primaryKey"`
	TeamID    string
	UserID    string
	FirstName string
	LastName  string
	Role      string
	CreatedAt time.Time
}
