package player

import (
	"bytes"
	"context"
	"io"
	"time"

	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	playerUsecase "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/player"
	playeruc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/player"
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
	playerUc  playerUsecase.PlayerUsecase
	logger    logger.Logger
	timeout   *time.Duration
	validater *validater.Validater
	pb.UnimplementedPlayerServiceServer
}

func NewPlayerHandler(puc playeruc.PlayerUsecase, logger logger.Logger, timeout *time.Duration, validater *validater.Validater) *PlayerHandler {
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

	data, err := ToAddPlayerReq(playerDetails)
	if err != nil {
		return apperror.GRPCStatus(err)
	}

	if validationError := ph.validater.Validation(data); validationError != nil {
		stWithDetails, err := validater.ValidationError(validationError)
		if err != nil {
			return status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return stWithDetails.Err()
	}

	playr, err := ph.playerUc.AddNewPlayer(context, &playerUsecase.AddPlayerReq{})

	if err != nil {
		return apperror.GRPCStatus(err)
	}

	return stream.SendAndClose(&pb.AddPlayerRes{
		PlayerId: playr.ID.String(),
	})
}

func (ph *PlayerHandler) UpdatePlayerDetails(ctx context.Context, input *pb.UpdatePlayerRequest) (*pb.UpdatePlayersResponse, error) {
	context, cancel := context.WithTimeout(ctx, *ph.timeout)
	defer cancel()
	data, err := ToUpdatePlayerReq(input)
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	if validationError := ph.validater.Validation(data); validationError != nil {
		stWithDetails, err := validater.ValidationError(validationError)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to attach details")
		}
		return nil, stWithDetails.Err()
	}

	res, err := ph.playerUc.UpdatePlayerDetails(context, &playeruc.UpdatePlayerReq{
		PlayerID:      data.PlayerID,
		UserID:        data.UserID,
		FullName:      data.FullName,
		DOB:           data.DateOfBirth,
		Nationality:   data.Nationality,
		PreferredFoot: data.PreferredFoot,
		Position:      data.Position,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.UpdatePlayersResponse{
		PlayerId:      res.PlayerID.String(),
		FullName:      res.FullName,
		DateOfBirth:   timestamppb.New(res.DOB),
		Position:      string(res.Position),
		Status:        string(res.Status),
		PreferredFoot: string(res.PreferredFoot),
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

	playerRes, err := ph.playerUc.GetPlayer(context, &playeruc.GetPlayerReq{
		PlayerID: data.PlayerID,
	})
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.GetPlayerRes{
		PlayerId:      playerRes.PlayerID.String(),
		FullName:      playerRes.FullName,
		DateOfBirth:   timestamppb.New(playerRes.DateOfBirth),
		Position:      string(playerRes.Position),
		PreferredFoot: string(playerRes.PreferedFoot),
		Nationality:   playerRes.Nationality,
		Status:        string(playerRes.Status),
		CreatedAt:     timestamppb.New(playerRes.CreatedAt),
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

	res, err := ph.playerUc.UpdateImage(ctx, &playeruc.UpdatePlayerImageReq{
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

	res, err := ph.playerUc.GetPlayerPresignedUrl(ctx, &playeruc.GetPlayerPresignedUrlReq{
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
