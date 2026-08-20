package entity

import (
	"time"

	"github.com/google/uuid"
)

type StaffDesignation string

const (
	StaffDesignationOwner   StaffDesignation = "OWNER"
	StaffDesignationManager StaffDesignation = "MANAGER"

	StaffDesignationHeadCoach       StaffDesignation = "HEAD_COACH"
	StaffDesignationAssistantCoach  StaffDesignation = "ASSISTANT_COACH"
	StaffDesignationGoalkeeperCoach StaffDesignation = "GOALKEEPER_COACH"
	StaffDesignationFitnessCoach    StaffDesignation = "FITNESS_COACH"

	StaffDesignationPhysio StaffDesignation = "PHYSIO"
	StaffDesignationDoctor StaffDesignation = "DOCTOR"

	StaffDesignationAnalyst     StaffDesignation = "ANALYST"
	StaffDesignationOther       StaffDesignation = "OTHER"
	StaffDesignationUnspecified StaffDesignation = "UNSPECIFIED"
)

type StaffRole string

const (
	StaffRoleManagement  StaffRole = "MANAGEMENT"
	StaffRoleCoach       StaffRole = "COACH"
	StaffRoleMedical     StaffRole = "MEDICAL"
	StaffRoleOther       StaffRole = "OTHER"
	StaffRoleUnspecified StaffRole = "UNSPECIFIED"
)

type Staff struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	TeamMemberID uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	FullName     string
	DateOfBirth  time.Time
	Role         StaffRole
	Designation  StaffDesignation
	PhoneNum     string
	ImageKey     string
	CreatedAt    time.Time
	UpdatedAt    time.Time   
	TeamMember   TeamMember `gorm:"foreignKey:TeamMemberID;references:ID"`
}
