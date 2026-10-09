package teaminvite

// type TeamInviteRepository interface {
// 	CreateInvitation(context.Context, *entity.TeamInvite) (*entity.TeamInvite, error)
// 	GetInvitation(context.Context, *uuid.UUID) (*entity.TeamInvite, bool, error)
// 	GetStaffInvitationByCode(context.Context, string) (*entity.TeamInvite, error)
// 	UpdateTeamInvite(ctxm context.Context, requestID uuid.UUID, teamInvite *entity.TeamInvite) error
// 	ExistPendingInvite(ctx context.Context, PlayerID uuid.UUID) (bool, error)
// }

// type teamInviteRepository struct {
// 	db     *gorm.DB
// 	logger logger.Logger
// }

// func NewTeamMemberInviteRepository(db *gorm.DB, logger logger.Logger) TeamInviteRepository {
// 	return &teamInviteRepository{
// 		db:     db,
// 		logger: logger,
// 	}
// }

// func (tr *teamInviteRepository) CreateInvitation(ctx context.Context, invite *entity.TeamInvite) (*entity.TeamInvite, error) {
// 	if err := tr.db.WithContext(ctx).Create(invite).Error; err != nil {
// 		tr.logger.Error("databaser error", "eroror", err)
// 		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
// 	}

// 	return invite, nil
// }

// func (tr *teamInviteRepository) GetInvitation(ctx context.Context, teamMemberID *uuid.UUID) (*entity.TeamInvite, bool, error) {

// 	var invite entity.TeamInvite
// 	result := tr.db.WithContext(ctx).Where("team_member_id=?", teamMemberID).Take(&invite)

// 	if result.Error != nil {
// 		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			return nil, false, nil
// 		}
// 		tr.logger.Error("database error", "error", result.Error, "method", "GetInvitation")
// 		return nil, false, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)

// 	}

// 	return &invite, true, nil

// }

// func (tr *teamInviteRepository) GetStaffInvitationByCode(ctx context.Context, code string) (*entity.TeamInvite, error) {

// 	var invite entity.TeamInvite

// 	result := tr.db.WithContext(ctx).Where("si.invite_code=?", code).
// 		Joins("JOIN staff_invite_details si ON team_invites.id=si.team_invite_id").
// 		Take(&invite)

// 	if result.Error != nil {
// 		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			return nil, apperror.NewNotFoundError("invitation not exist")
// 		}
// 		tr.logger.Error("database error", "error", result.Error, "method", "GetInvitation")
// 		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, result.Error)

// 	}

// 	return &invite, nil
// }

// func (tr *teamInviteRepository) UpdateTeamInvite(ctx context.Context, inviteID uuid.UUID, teamInvite *entity.TeamInvite) error {
// 	if err := tr.db.WithContext(ctx).Model(&entity.TeamInvite{}).Where("id=?", teamInvite).Updates(teamInvite).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return apperror.NewNotFoundError("team invitation details not exist")
// 		}
// 		tr.logger.Error("database error", "error", err, "method", "repo.updateTeamInvite")
// 		return apperror.NewInternalError(apperror.InternalErrorMsg, err)
// 	}

// 	return nil
// }
