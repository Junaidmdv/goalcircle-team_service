package player

import (
	"time"

	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	teamv1 "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
)

func ToAddPlayerReq(input *pb.PlayerDetails) *AddPlayerReq {

	position := MapPlayerPosition(teamv1.PlayerPosition(input.Position))

	return &AddPlayerReq{
		TeamID:       input.TeamId,
		UserID:       input.UserId,
		FullName:     input.FullName,
		DateOfBirth:  input.DateOfBirth.AsTime(),
		JerseyNumber: input.JerseyNumber,
		Position:     position,
		Height:       input.Height,
		Weight:       input.Weight,
	}
}

func ToUpdatePlayerReq(input *pb.UpdatePlayerRequest) *UpdatePlayerReq {
	var dateOfBirth *time.Time
	if input.DateOfBirth != nil {
		t := input.DateOfBirth.AsTime()
		dateOfBirth = &t
	}

	var position *entity.PlayerPosition
	if input.Position != nil {
		p := MapPlayerPosition(*input.Position)
		position = &p
	}

	var status *entity.PlayerStatus
	if input.Status != nil {
		s := MapPlayerStatus(*input.Status)
		status = &s
	}

	return &UpdatePlayerReq{
		UserID:       input.UserId,
		TeamID:       input.TeamId,
		PlayerID:     input.PlayerId,
		FullName:     input.FullName,
		DateOfBirth:  dateOfBirth,
		JerseyNumber: input.JerseyNumber,
		Position:     position,
		Height:       input.Height,
		Weight:       input.Weight,
		Status:       status,
	}
}

func ToGetPlayerReq(input *pb.GetPlayerReq) *GetPlayerReq {
	return &GetPlayerReq{
		PlayerID: input.PlayerId,
	}
}
