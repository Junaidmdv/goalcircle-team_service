package player

import (
	"bytes"
	"context"
	"io"
	"time"

	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/usecase/player"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/imageutil"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/validater"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PlayerHandler struct {
	playerUc  player.PlayerUsecase
	logger    logger.Logger
	timeout   *time.Duration
	validater *validater.Validater
	pb.UnimplementedPlayerServiceServer
}

func NewPlayerHandler(puc player.PlayerUsecase, logger logger.Logger, timeout *time.Duration, validater *validater.Validater) *PlayerHandler {
	return &PlayerHandler{
		playerUc:  puc,
		logger:    logger,
		timeout:   timeout,
		validater: validater,
	}
}

func (ph *PlayerHandler) AddNewPlayer(stream grpc.ClientStreamingServer[pb.AddPlayerReq, pb.AddPlayerRes]) error {

	ctx := stream.Context()

	context, cancel := context.WithTimeout(ctx, *ph.timeout)
	defer cancel()

	var (
		playerDetails *pb.PlayerDetails
		buffer        bytes.Buffer
		imageSize     int
	)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to receive stream: %v", err)
		}

		switch data := req.Player.(type) {
		case *pb.AddPlayerReq_PlayerDetails:
			if playerDetails != nil {
				return status.Error(codes.InvalidArgument, "player details already sent")
			}
			playerDetails = data.PlayerDetails
		case *pb.AddPlayerReq_PlayerImageChunks:
			if playerDetails == nil {
				return status.Error(codes.InvalidArgument, "player details is missing")
			}
			buffer.Write(data.PlayerImageChunks)
			imageSize += len(data.PlayerImageChunks)
		}
	}

	contentType, err := imageutil.ValidateImage(buffer.Bytes(), imageutil.PlayerImage)
	if err != nil {
		return apperror.GRPCStatus(err)
	}

	data := ToAddPlayerReq(playerDetails)

	if validationError := ph.validater.Validation(data); validationError != nil {
		stWithDetails, err := validater.ValidationError(validationError)
		if err != nil {
			return status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return stWithDetails.Err()
	}

	playr, err := ph.playerUc.AddNewPlayer(context, &player.AddPlayerReq{
		TeamID:       data.TeamID,
		UserID:       data.UserID,
		DOB:          data.DateOfBirth,
		JerseyNumber: data.JerseyNumber,
		Postion:      data.Position,
		Height:       data.Height,
		Weight:       data.Weight,
		ImageBytes:   buffer.Bytes(),
		ContentType:  contentType,
		ImageSize:    imageSize,
	})

	if err != nil {
		return apperror.GRPCStatus(err)
	}

	position := MapProtoPlayerPosition(playr.Position)
	status := MapProtoPlayerStatus(playr.Status)

	return stream.SendAndClose(&pb.AddPlayerRes{
		PlayerId:       playr.ID.String(),
		TeamMemberId:   playr.TeamMemberID.String(),
		FullName:       playr.FullName,
		JerseyNumber:   playr.JerseyNumber,
		PlayerStatus:   status,
		PlayerPosition: position,
	})
}

func (ph *PlayerHandler) UpdateUserStatus(ctx context.Context, input *pb.UpdatePlayerStatusReq) (*pb.UpdatePlayerStatusRes, error) {
	context, cancel := context.WithTimeout(ctx, *ph.timeout)
	defer cancel()
	data := ToUpdateStatusReq(input)

	if validationError := ph.validater.Validation(data); validationError != nil {
		stWithDetails, err := validater.ValidationError(validationError)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := ph.playerUc.UpdatePlayerStatus(context, &player.UpdatPlayerStatusReq{
		PlayerID: data.PlayerID,
		Status:   data.PlayerStatus,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	ph.logger.Info("user login succefull", "response", res)

	return &pb.UpdatePlayerStatusRes{
		Success: res.Success,
	}, nil
}

func (ph *PlayerHandler) ListTeamPlayer(ctx context.Context, input *pb.ListTeamPlayerReq) (*pb.ListTeamPlayerRes, error) {
	context, cancel := context.WithTimeout(ctx, *ph.timeout)
	defer cancel()

	users, Paginates, err := ph.playerUc.ListTeamPlayers(context, &player.ListTeamPlayersReq{
		TeamID: input.TeamId,
	})
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	var pbPlayers []*pb.PlayerList

	for _, user := range users {
		position := MapProtoPlayerPosition(user.Position)
		pbPlayers = append(pbPlayers, &pb.PlayerList{
			PlayerId:     user.ID.String(),
			TeamMemberId: user.TeamMemberID.String(),
			FullName:     user.FullName,
			JerseyNumber: user.JerseyNumber,
			Position:     position,
		})
	}

	paginateDetails := &pb.PaginationDetails{
		TotalPage: int32(Paginates.Total),
		Limit:     Paginates.Limit,
		TotalItem: Paginates.Total,
		Page:      Paginates.Page,
	}

	return &pb.ListTeamPlayerRes{
		Players:  pbPlayers,
		Paginate: paginateDetails,
	}, nil
}

func (ph *PlayerHandler) GetPlayer(ctx context.Context, input *pb.GetPlayerReq) (*pb.GetPlayerRes, error) {
	context, cancel := context.WithTimeout(ctx, *ph.timeout)
	defer cancel()

	data := ToGetPlayerReq(input)

	if validationError := ph.validater.Validation(data); validationError != nil {
		stWithDetails, err := validater.ValidationError(validationError)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	playerRes, err := ph.playerUc.GetPlayer(context, &player.GetPlayerReq{
		PlayerID: data.PlayerID,
	})
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	playerPosition := MapProtoPlayerPosition(playerRes.Position)
	playerStatus := MapProtoPlayerStatus(playerRes.Status)

	return &pb.GetPlayerRes{
		PlayerId:       playerRes.PlayerID.String(),
		TeamMemberId:   playerRes.TeamMemberID.String(),
		FullName:       playerRes.FullName,
		Dob:            timestamppb.New(playerRes.DateOfBirth),
		JerseyNumber:   playerRes.JerseyNumber,
		PlayerPosition: playerPosition,
		Height:         playerRes.Height,
		Weight:         playerRes.Weight,
		PlayerStatus:   playerStatus,
		CreatedAt:      timestamppb.New(playerRes.CreatedAt),
	}, nil
}

func (ph *PlayerHandler) ReleasePlayer(ctx context.Context, input *pb.ReleasePlayerReq) (*pb.ReleasePlayerRes, error) {
	return nil, nil
}
