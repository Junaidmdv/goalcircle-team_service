package playerinviteuc

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/permission"
	playerrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/player"
	playerinviterepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/playerinvite"
	staffrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staff"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teamplayer"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage"

	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
)

type PlayerInviteUsecase interface {
	SendPlayerRequest(context.Context, *SendPlayerInvitationReq) (*SendPlayerInvitationRes, error)
	TeamInviteList(ctx context.Context, req *TeamInviteListReq) (*TeamInviteListRes, error)
	AcceptPlayerInvite(ctx context.Context, req *AcceptPlayerInviteReq) (*AcceptPlayerInviteRes, error)
	RejectPlayerInvite(ctx context.Context, req *RejectPlayerInviteReq) (*RejectPlayerInviteRes, error)
	PlayerInviteCancel(ctx context.Context, req *PlayerInviteCancelReq) (*PlayerInviteCancelRes, error)
	PlayerInvites(ctx context.Context, req *PlayerInvitesReq) (*PlayerInvitesRes, error)
}

type playerInviteUsecase struct {
	staffRepo           staffrepo.StaffRepository
	teamPlayerRepo      teamplayer.TeamPlayerRepository
	playerInviteRepo    playerinviterepo.PlayerInviteRepository
	teamMemberRepo      teammember.TeamMemberRepository
	playerRepo          playerrepo.PlayerRepository
	ObjectStorage       storage.ObjectStorage
	ObjectStorageConfig config.ObjectStorageConfig
}

func NewPlayerInviteUsecase(
	staffrepo staffrepo.StaffRepository,
	teamplayerrepo teamplayer.TeamPlayerRepository,
	playerinviterepo playerinviterepo.PlayerInviteRepository,
	playerrepo playerrepo.PlayerRepository,
	objectstorage storage.ObjectStorage,
	teammemberrepo teammember.TeamMemberRepository,

) PlayerInviteUsecase {
	return &playerInviteUsecase{
		staffRepo:        staffrepo,
		teamPlayerRepo:   teamplayerrepo,
		playerInviteRepo: playerinviterepo,
		playerRepo:       playerrepo,
		ObjectStorage:    objectstorage,
		teamMemberRepo:   teammemberrepo,
	}
}

func (u *playerInviteUsecase) SendPlayerRequest(ctx context.Context, req *SendPlayerInvitationReq) (*SendPlayerInvitationRes, error) {

	userId, err := uuid.Parse(req.SendUserID)
	if err != nil {
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, errors.New("invalid token"))
	}

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	desig, err := u.staffRepo.GetStaffDesignation(ctx, userId, teamID)
	if err != nil {
		return nil, err
	}

	isPermited := permission.HasPermission(desig, permission.PermissionAddPlayer)
	if !isPermited {
		return nil, apperror.NewBadRequestError("user not permitted add player")
	}

	playerID, err := uuid.Parse(req.PlayerID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid player id")
	}

	isMember, err := u.teamPlayerRepo.IsPlayerPartOfTeam(ctx, teamID, playerID)
	if err != nil {
		return nil, err
	}

	if isMember {
		return nil, apperror.NewBadRequestError("player already part of squad")
	}

	isInviteSent, err := u.playerInviteRepo.ExistPendingInvite(ctx, playerID)
	if err != nil {
		return nil, err
	}

	if isInviteSent {
		return nil, apperror.NewConflictError("an invite has already been sent to this player")
	}

	if req.OfferedJerseyNum != nil {
		isJerseyNumTaken, err := u.teamPlayerRepo.IsJerseyNumberTaken(ctx, teamID, req.OfferedJerseyNum)
		if err != nil {
			return nil, err
		}
		if isJerseyNumTaken {
			return nil, apperror.NewConflictError("jersey number already taken")
		}
	}

	player, err := u.playerRepo.GetPlayer(ctx, playerID)
	if err != nil {
		return nil, err
	}

	playerInviteDetails, err := u.playerInviteRepo.AddPlayerInviteDetails(ctx, &entity.PlayerTeamInvite{
		ID:               uuid.New(),
		TeamID:           teamID,
		PlayerID:         playerID,
		OfferedJerseyNum: req.OfferedJerseyNum,
		Message:          req.Message,
		InvitedByUserID:  userId,
		InvitedUserID:    player.UserID,
		Status:           entity.InvitationPending,
		ExpiresAt:        time.Now().Add(time.Hour * 24 * entity.PlayerRequestExpiry),
	})
	if err != nil {
		return nil, err
	}

	return &SendPlayerInvitationRes{
		RequestID: playerInviteDetails.ID.String(),
		Status:    string(playerInviteDetails.Status),
		CreatedAt: playerInviteDetails.CreatedAt,
		ExpiresAt: playerInviteDetails.ExpiresAt,
	}, nil
}

func (u *playerInviteUsecase) TeamInviteList(ctx context.Context, req *TeamInviteListReq) (*TeamInviteListRes, error) {

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	switch {
	case req.Limit > entity.MaxDefaultPaginateLimit:
		req.Limit = entity.MaxDefaultPaginateLimit
	case req.Limit <= 0:
		req.Limit = entity.MinDefaultPagination
	}

	teamPlayers, total, err := u.playerInviteRepo.GetTeamPlayerInvites(ctx, &playerinviterepo.PaginateDetails{
		TeamID: teamID,
		Page:   req.Page,
		Limit:  req.Limit,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	var invites []InvitationDetails

	for i, tp := range teamPlayers {
		player, err := u.playerRepo.GetPlayer(ctx, tp.PlayerID)
		if err != nil {
			log.Printf("failed to get player %d for invite %d: %v", tp.PlayerID, tp.ID, err)
			continue
		}

		status := tp.Status
		if status == entity.InvitationPending && tp.ExpiresAt.Before(time.Now()) {
			status = entity.InvitationExpired
		}

		invites = append(invites, InvitationDetails{
			RequestID:   tp.ID.String(),
			Status:      status,
			CreatedAt:   tp.CreatedAt,
			ExpiresAt:   tp.ExpiresAt,
			RespondedAt: tp.RespondedAt,
			Player: PlayerDetails{
				ID:            player.ID.String(),
				FullName:      player.FullName,
				Position:      player.Position,
				PreferredFoot: string(player.PreferredFoot),
			},
		})
		presignedUrl, err := u.ObjectStorage.GetPresignedURL(ctx, u.ObjectStorageConfig.Bucket, *player.ProfileImageKey, u.ObjectStorageConfig.PresignedURLExpiry)
		if err == nil {
			invites[i].Player.ImageUrl = presignedUrl
		}

	}

	totalPage := total / int64(req.Limit)
	if total%int64(req.Limit) > 0 {
		totalPage += 1
	}

	return &TeamInviteListRes{
		Data: invites,
		Paginate: &PaginationDetails{
			TotalPage: int32(totalPage),
			Page:      int32(req.Page),
			Limit:     int32(req.Limit),
			Total:     int64(total),
		},
	}, nil
}

func (u *playerInviteUsecase) AcceptPlayerInvite(ctx context.Context, req *AcceptPlayerInviteReq) (*AcceptPlayerInviteRes, error) {
	invitedID, err := uuid.Parse(req.InviteID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid invitation id")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid user id")
	}

	playerInvite, err := u.playerInviteRepo.GetTeamPlayerInvite(ctx, invitedID)
	if err != nil {
		return nil, err
	}

	if playerInvite.InvitedUserID != userID {
		return nil, apperror.NewPermissionDeniedError("you are not authorized to accept this invitation")
	}

	if playerInvite.Status != entity.InvitationPending || time.Now().After(playerInvite.ExpiresAt) {
		return nil, apperror.NewConflictError("invitation is no longer valid")
	}

	if err := u.playerInviteRepo.UpdatePlayerInviteRes(ctx, invitedID, entity.InvitationAccepted); err != nil {
		return nil, err
	}

	member, err := u.teamMemberRepo.AddTeamMember(ctx, &entity.TeamMember{
		ID:     uuid.New(),
		TeamID: playerInvite.TeamID,
		UserID: playerInvite.InvitedUserID,
		Role:   entity.TeamMemberPlayer,
		Status: entity.TeamMemberStatusActive,
	})
	if err != nil {
		return nil, err
	}

	teamplyer, err := u.teamPlayerRepo.AddTeamPlayer(ctx, &entity.TeamPlayer{
		ID:           uuid.New(),
		PlayerID:     playerInvite.PlayerID,
		TeamMemberID: member.ID,
		JerseyNumber: playerInvite.OfferedJerseyNum,
		JoinedAt:     time.Now(),
		Status:       entity.TeamPlayerStatusActive,
	})
	if err != nil {
		return nil, err
	}

	return &AcceptPlayerInviteRes{
		InviteID:      invitedID,
		RequestStatus: string(playerInvite.Status),
		RespondedAt:   playerInvite.RespondedAt,
		TeamPlayerID:  teamplyer.ID,
		TeamID:        member.TeamID,
	}, nil
}

func (u *playerInviteUsecase) RejectPlayerInvite(ctx context.Context, req *RejectPlayerInviteReq) (*RejectPlayerInviteRes, error) {

	invitedID, err := uuid.Parse(req.InviteID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid invited id")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid user id")
	}

	playerInvite, err := u.playerInviteRepo.GetTeamPlayerInvite(ctx, invitedID)
	if err != nil {
		return nil, err
	}

	if playerInvite.InvitedUserID != userID {
		return nil, apperror.NewPermissionDeniedError("you are not authorized to accept this invitation")
	}

	if playerInvite.Status != entity.InvitationPending || time.Now().After(playerInvite.ExpiresAt) {
		return nil, apperror.NewConflictError("invitation is no longer valid")
	}

	if err := u.playerInviteRepo.UpdatePlayerInviteRes(ctx, invitedID, entity.InvitationRejected); err != nil {
		return nil, err
	}

	return &RejectPlayerInviteRes{
		InviteID:    playerInvite.ID,
		RespondedAt: playerInvite.RespondedAt,
	}, nil
}

func (u *playerInviteUsecase) PlayerInviteCancel(ctx context.Context, req *PlayerInviteCancelReq) (*PlayerInviteCancelRes, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid user id")
	}

	teamID, err := uuid.Parse(req.InviteID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid team id")
	}

	desig, err := u.staffRepo.GetStaffDesignation(ctx, userID, teamID)
	if err != nil {
		return nil, err
	}

	isPermited := permission.HasPermission(desig, permission.PermissionCancelPlayerInvite)
	if !isPermited {
		return nil, apperror.NewBadRequestError("user not permitted cancel player request")
	}

	invitedID, err := uuid.Parse(req.InviteID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid request id")
	}

	playerInvite, err := u.playerInviteRepo.GetTeamPlayerInvite(ctx, invitedID)
	if err != nil {
		return nil, err
	}

	if playerInvite.TeamID != teamID {
		return nil, apperror.NewBadRequestError("invalid team invitation")
	}

	if playerInvite.Status != entity.InvitationPending || time.Now().After(playerInvite.ExpiresAt) {
		return nil, apperror.NewConflictError("invitation is no longer valid")
	}

	if err := u.playerInviteRepo.UpdatePlayerInviteRes(ctx, invitedID, entity.InvitationCancelled); err != nil {
		return nil, err
	}

	return &PlayerInviteCancelRes{
		InviteID: playerInvite.ID,
	}, nil
}

func (u *playerInviteUsecase) PlayerInvites(ctx context.Context, req *PlayerInvitesReq) (*PlayerInvitesRes, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid player id")
	}

	playerinvites, err := u.playerInviteRepo.PlayerInvitationList(ctx, userID)
	if err != nil {
		return nil, err
	}

	var invites []PlayerInvitationDetails

	for _, invite := range playerinvites {
		invites = append(invites, PlayerInvitationDetails{
			InviteID:         invite.ID,
			InvitedByUser:    invite.InvitedByUserID,
			TeamID:           invite.TeamID,
			OfferedJerseyNum: invite.OfferedJerseyNum,
			Message:          invite.Message,
			CreatedAt:        invite.CreatedAt,
		})
	}

	return &PlayerInvitesRes{
		UserID:            playerinvites[0].InvitedUserID,
		InvitationDetails: invites,
	}, nil
}
