package player

import (
	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	teamv1 "github.com/Junaidmdv/goalcircle-protos/team/v1"
)

func ToAddPlayerReq(input *pb.AddPlayerReq) *AddPlayerReq {

	position := MapPlayerPosition(teamv1.PlayerPosition(input.Position))

	return &AddPlayerReq{
		TeamID:       input.TeamId,
		TeamMemberID: input.TeamMemeberId,
		FullName:     input.FullName,
		DateOfBirth:  input.DateOfBirth.AsTime(),
		JerseyNumber: input.JerseyNumber,
		Position:     position,
		Height:       input.Height,
		Weight:       input.Weight,
	}
}

func ToUpdateStatusReq(input *pb.UpdatePlayerStatusReq) *UpdatePlayerStatusReq {

	playerstatus := MapPlayerStatus(input.PlayerStatus)

	return &UpdatePlayerStatusReq{
		PlayerID:     input.PlayerId,
		PlayerStatus: playerstatus,
	}
}

func ToGetPlayerReq(input *pb.GetPlayerReq) *GetPlayerReq {
	return &GetPlayerReq{
		PlayerID: input.PlayerId,
	}
}
