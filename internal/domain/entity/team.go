package entity

import "time"

type Team struct {
	ID          string
	UUID        string `gorm:"primaryKey"`
	Name        string
	City        string
	Logo        string
	ManagerID   string
	TeamStatus  string
	PlayerCount string 
	CreatedAt   time.Time   
}



