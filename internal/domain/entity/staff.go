package entity

import (
	"time"

	"github.com/google/uuid"
)

type StaffDesignation string

const (
	StaffDesignationManager         StaffDesignation = "MANAGER"
	StaffDesignationCoach           StaffDesignation = "COACH"
	StaffDesignationAssistantCoach  StaffDesignation = "ASSISTANT_COACH"
	StaffDesignationPhysiotherapist StaffDesignation = "PHYSIOTHERAPIST"
	StaffDesignationOther           StaffDesignation = "OTHER"
)

type Staff struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	TeamMemberID uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	FullName     string
	Designation  StaffDesignation
	CreatedAt    time.Time
	UpdatedAt    time.Time
}


