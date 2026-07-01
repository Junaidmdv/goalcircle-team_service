package team

import "github.com/google/uuid"

type UpdateTeamRes struct {
	TeamID uuid.UUID  
	TeamName string 
}