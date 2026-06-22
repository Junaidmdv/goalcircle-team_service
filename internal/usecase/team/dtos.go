package team

import "github.com/google/uuid"

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
}
