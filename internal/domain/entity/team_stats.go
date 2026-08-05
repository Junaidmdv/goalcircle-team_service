package entity

import "github.com/google/uuid"

type TeamStats struct {
	ID            uuid.UUID `gorm:"primaryKey"`
	TeamID        uuid.UUID  
	MatchesPlayed int
	Wins          int
	Draws         int
	Losses        int
	GoalsScored   int
	GoalsConceded int  
}
