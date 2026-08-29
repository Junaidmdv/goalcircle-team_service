package transaction

import (
	staffrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staff"
	teamrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
	teammemberrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
)

type Repositories interface {
	Staff() staffrepo.StaffRepository
	Team() teamrepo.TeamRepository
	TeamMember() teammemberrepo.TeamMemberRepository
}

type repositories struct {
	teamRepo       teamrepo.TeamRepository
	staffRepo      staffrepo.StaffRepository
	teamMemberRepo teammemberrepo.TeamMemberRepository
}

func NewRepositories(tmRepo teamrepo.TeamRepository, staffrepo staffrepo.StaffRepository, tmMemberRepo teammemberrepo.TeamMemberRepository) Repositories {
	return &repositories{
		teamRepo:       tmRepo,
		staffRepo:      staffrepo,
		teamMemberRepo: tmMemberRepo,
	}
}

func (r *repositories) Team() teamrepo.TeamRepository {
	return r.teamRepo
}

func (r *repositories) Staff() staffrepo.StaffRepository {
	return r.staffRepo
}

func (r *repositories) TeamMember() teammemberrepo.TeamMemberRepository {
	return r.teamMemberRepo
}
