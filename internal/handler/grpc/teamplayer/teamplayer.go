package teamplayerhandler

import (
	"context"
	"time"

	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	teamplayeruc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/teamplayer"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TeamPlayerHandler struct {
	teamPlayerUse teamplayeruc.TeamPlayerUsecase
	timeout       time.Duration
}

func NewTeamPlayerHandler() *TeamPlayerHandler {
	return &TeamPlayerHandler{}
}

func (ph *TeamPlayerHandler) ReleasePlayer(ctx context.Context, input *pb.ReleasePlayerReq) (*pb.ReleasePlayerRes, error) {

	ctx, cancel := context.WithTimeout(ctx, ph.timeout)
	defer cancel()

	res, err := ph.teamPlayerUse.ReleasePlayer(ctx, &teamplayeruc.ReleasePlayerReq{
		UserID:   input.UserId,
		TeamID:   input.TeamId,
		PlayerID: input.PlayerId,
	})
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.ReleasePlayerRes{
		Id:           res.ID.String(),
		TeamMemberId: res.TeamMemberID.String(),
		FullName:     res.FullName,
		JerseyNumber: res.JerseyNumber,
		Position:     string(res.Position),
		Status:       string(res.Status),
		ReleasedAt:   timestamppb.New(res.ReleasedAt),
		JoinedAt:     timestamppb.New(res.JoinedAt),
	}, nil
}

func (ph *TeamPlayerHandler) ListTeamPlayer(ctx context.Context, input *pb.ListTeamPlayerReq) (*pb.ListTeamPlayerRes, error) {
	context, cancel := context.WithTimeout(ctx, ph.timeout)
	defer cancel()

	status, err := toPlayerStatus(*input.PlayerStatus)
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}
	position, err := toPlayerPosition(*input.Position)
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	players, Paginates, err := ph.teamPlayerUse.ListTeamPlayers(context, &teamplayeruc.ListTeamPlayersReq{
		TeamID:       input.TeamId,
		Page:         input.Page,
		Limit:        input.Limit,
		PlayerStatus: status,
		Position:     position,
		Search:       *input.Search,
	})
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	var pbPlayers []*pb.PlayerList

	for _, player := range players {

		pbPlayers = append(pbPlayers, &pb.PlayerList{
			PlayerId:     player.ID.String(),
			TeamMemberId: player.TeamMemberID.String(),
			FullName:     player.FullName,
			Position:     string(player.Position),
			ImageUrl:     player.ImageUrl,
		})
	}

	paginateDetails := &pb.PaginationDetails{
		TotalPage: int32(Paginates.TotalPage),
		Limit:     Paginates.Limit,
		TotalItem: Paginates.Total,
		Page:      Paginates.Page,
	}

	return &pb.ListTeamPlayerRes{
		Players:  pbPlayers,
		Paginate: paginateDetails,
	}, nil
}
