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
		}
	}

	err := imageutil.ValidateImage(buffer.Bytes(), imageutil.PlayerImage)
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
		FullName:     data.FullName,
		DOB:          data.DateOfBirth,
		JerseyNumber: data.JerseyNumber,
		Postion:      data.Position,
		Height:       data.Height,
		Weight:       data.Weight,
		ImageBytes:   buffer.Bytes(),
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
		PlayerStatus:   status.String(),
		PlayerPosition: position.String(),
		PresignedUrl:   playr.PresignedUrl,
	})
}

func (ph *PlayerHandler) UpdatePlayerDetails(ctx context.Context, input *pb.UpdatePlayerRequest) (*pb.UpdatePlayersResponse, error) {
	context, cancel := context.WithTimeout(ctx, *ph.timeout)
	defer cancel()
	data := ToUpdatePlayerReq(input)

	if validationError := ph.validater.Validation(data); validationError != nil {
		stWithDetails, err := validater.ValidationError(validationError)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := ph.playerUc.UpdatePlayerDetails(context, &player.UpdatePlayerReq{
		UserID:       data.UserID,
		PlayerID:     data.PlayerID,
		TeamID:       data.TeamID,
		FullName:     data.FullName,
		DateOfBirth:  data.DateOfBirth,
		JerseyNumber: data.JerseyNumber,
		Position:     data.Position,
		Height:       data.Height,
		Weight:       data.Weight,
		Status:       data.Status,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.UpdatePlayersResponse{
		PlayerId:     res.PlayerID.String(),
		TeamMemberId: res.TeamMemberID.String(),
		FullName:     res.FullName,
		DateOfBirth:  timestamppb.New(res.DateOfBirth),
		JerseyNumber: res.JerseyNumber,
		Position:     string(res.Position),
		Height:       res.Height,
		Weight:       res.Weight,
		Status:       string(res.Status),
	}, nil
}

func (ph *PlayerHandler) ListTeamPlayer(ctx context.Context, input *pb.ListTeamPlayerReq) (*pb.ListTeamPlayerRes, error) {
	context, cancel := context.WithTimeout(ctx, *ph.timeout)
	defer cancel()

	status := MapPlayerStatus(*input.PlayerStatus)
	position := MapPlayerPosition(*input.Position)

	players, Paginates, err := ph.playerUc.ListTeamPlayers(context, &player.ListTeamPlayersReq{
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
			JerseyNumber: player.JerseyNumber,
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
	res, err := ph.playerUc.ReleasePlayer(ctx, &player.ReleasePlayerReq{
		UserID:   input.UserId,
		TeamID:   input.TeamId,
		PlayerID: input.PlayerId,
	})
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	position := MapProtoPlayerPosition(res.Position)
	status := MapProtoPlayerStatus(res.Status)
	return &pb.ReleasePlayerRes{
		Id:           res.ID.String(),
		TeamMemberId: res.TeamMemberID.String(),
		FullName:     res.FullName,
		JerseyNumber: res.JerseyNumber,
		Position:     position,
		Status:       status,
		ReleasedAt:   timestamppb.New(res.ReleasedAt),
		JoinedAt:     timestamppb.New(res.JoinedAt),
	}, nil
}

func (ph *PlayerHandler) UpdatePlayerImage(stream grpc.ClientStreamingServer[pb.UpdatePlayerImageReq, pb.UpdatePlayerImageRes]) error {
	ctx := stream.Context()

	ctx, cancel := context.WithTimeout(ctx, *ph.timeout)
	defer cancel()

	var (
		metadata *pb.UpdatePlayerImageMeta
		buffer   bytes.Buffer
	)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to receive stream: %v", err)
		}

		switch data := req.Data.(type) {
		case *pb.UpdatePlayerImageReq_Meta:
			if metadata != nil {
				return status.Error(codes.InvalidArgument, "player details already sent")
			}
			metadata = data.Meta
		case *pb.UpdatePlayerImageReq_Chunks:
			if metadata == nil {
				return status.Error(codes.InvalidArgument, "player details is missing")
			}
			buffer.Write(data.Chunks)
		}
	}
	err := imageutil.ValidateImage(buffer.Bytes(), imageutil.PlayerImage)
	if err != nil {
		return apperror.GRPCStatus(err)
	}

	res, err := ph.playerUc.UpdateImage(ctx, &player.UpdatePlayerImageReq{
		UserID:    metadata.UserId,
		TeamID:    metadata.TeamId,
		PlayerID:  metadata.PlayerId,
		ImageData: buffer.Bytes(),
	})

	return stream.SendAndClose(&pb.UpdatePlayerImageRes{
		TeamId:       res.TeamID,
		PlayerId:     res.PlayerID,
		PresignedUrl: res.PresignedUrl,
	})

}

func (ph *PlayerHandler) GetPlayerPresignedUrl(ctx context.Context, req *pb.GetPlayerPresignedUrlReq) (*pb.GetPlayerPresignedUrlRes, error) {
	ctx, cancel := context.WithTimeout(ctx, *ph.timeout)
	defer cancel()

	res, err := ph.playerUc.GetPlayerPresignedUrl(ctx, &player.GetPlayerPresignedUrlReq{
		UserID:   req.UserId,
		PlayerID: req.PlayerId,
		TeamID:   req.TeamId,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}
	return &pb.GetPlayerPresignedUrlRes{
		TeamId:       res.TeamID,
		PlayerId:     res.PlayerId,
		PresignedUrl: res.PresignedUrl,
	}, nil
}
