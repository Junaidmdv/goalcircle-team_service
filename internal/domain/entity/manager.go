package entity

import "time"

type TeamManager struct {
	ID        string
	UUID      string `gorm:"primaryKey"`
	UserId    string
	Name      string
	Phone     string
	CreatedAt time.Time
}


