package teaminvite

import (
	"context"

	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"gorm.io/gorm"
)

type TeamInviteRepository interface {
	CreateInvitation(context.Context, *entity.TeamInvite) (*entity.TeamInvite, error)
}

type teamInviteRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewTeamMemberInviteRepository(db *gorm.DB) TeamInviteRepository {
	return &teamInviteRepository{
		db: db,
	}
}

func (tr *teamInviteRepository) CreateInvitation(ctx context.Context, invite *entity.TeamInvite) (*entity.TeamInvite, error) {
	if err := tr.db.WithContext(ctx).Create(invite).Error; err != nil {
		tr.logger.Error("databaser error", "eroror", err)
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return invite, nil
}


