package entity

import (
	"time"

	"github.com/google/uuid"
)

type StaffDesignation string

const (
	StaffDesignationCoach           StaffDesignation = "COACH"
	StaffDesignationAssistantCoach  StaffDesignation = "ASSISTANT_COACH"
	StaffDesignationManager         StaffDesignation = "MANAGER"
	StaffDesignationPhysiotherapist StaffDesignation = "PHYSIOTHERAPIST"
	StaffDesignationDoctor          StaffDesignation = "DOCTOR"
	StaffDesignationAnalyst         StaffDesignation = "ANALYST"
	StaffDesignationMediaManager    StaffDesignation = "MEDIA_MANAGER"
	StaffDesignationTrainer         StaffDesignation = "TRAINER"
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
