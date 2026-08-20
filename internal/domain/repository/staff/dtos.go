package staffrepo

import "github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"

type ListStaffDetails struct {
	Page        int32
	Limit       int32
	Role        entity.StaffRole
	Designation entity.StaffDesignation
	Search      string
}


