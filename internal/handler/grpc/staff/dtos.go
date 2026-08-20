package staff

import (
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
)

type AddStaffReq struct {
	UserID      string                  `validate:"required"`
	TeamID      string                  `validate:"required"`
	FullName    string                  `validate:"required,min=2,max=100"`
	Role        entity.StaffRole        `validate:"required,staff-role"`
	Designation entity.StaffDesignation `validate:"required,staff-desig"`
	DOB         time.Time               `validate:"required"`
}

type UpdateStaffReq struct { 
	UserID      string                   `validate:"required"` 
	TeamID      string                   `validate:"required"`
	StaffID     string                   `validate:"required"`
	FullName    *string                  `validate:"omitempty,min=2,max=100"`
	DateOfBirth *time.Time               `validate:"omitempty"`
	Role        *entity.StaffRole        `validate:"omitempty"`
	Designation *entity.StaffDesignation `validate:"omitempty"`
	PhoneNum    *string                  `validate:"omitempty,e164"`
}
