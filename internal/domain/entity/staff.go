package entity

import (
	"time"

	"github.com/google/uuid"
)

type StaffDesignation string

const (
	StaffDesignationHeadCoach        StaffDesignation = "HEAD_COACH"
	StaffDesignationAssistantCoach   StaffDesignation = "ASSISTANT_COACH"
	StaffDesignationGoalkeeperCoach  StaffDesignation = "GOALKEEPER_COACH"
	StaffDesignationFitnessCoach     StaffDesignation = "FITNESS_COACH"
	StaffDesignationPhysio           StaffDesignation = "PHYSIO"
	StaffDesignationDoctor           StaffDesignation = "DOCTOR"
	StaffDesignationAnalyst          StaffDesignation = "ANALYST"
	StaffDesignationTeamCoordinator  StaffDesignation = "TEAM_COORDINATOR"
	StaffDesignationEquipmentManager StaffDesignation = "EQUIPMENT_MANAGER"
	StaffDesignationOther            StaffDesignation = "OTHER"
)


type StaffStatus string

const (
	StaffStatusPending   StaffStatus = "PENDING"
	StaffStatusActive    StaffStatus = "ACTIVE"
	StaffStatusSuspended StaffStatus = "SUSPENDED"
	StaffStatusArchived  StaffStatus = "ARCHIVED"
)

type Staff struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	TeamMemberID uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	FullName     string
	DateOfBirth  time.Time
	Designation  StaffDesignation
	Status       StaffStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}


