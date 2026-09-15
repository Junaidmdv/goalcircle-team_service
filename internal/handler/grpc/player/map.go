package player

import (
	"time"

	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
)

func ToAddPlayerReq(input *pb.PlayerDetails) (*AddPlayerReq, error) {

	position, err := toPlayerPosition(input.Position)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid player position")
	}
	preferedfoot, err := toPreferredFoot(input.PreferedFoot)
	if err != nil {
		return nil, apperror.NewBadRequestError(err.Error())
	}

	return &AddPlayerReq{
		UserID:        input.UserId,
		FullName:      input.FullName,
		DateOfBirth:   input.Dob.AsTime(),
		Position:      position,
		PreferredFoot: preferedfoot,
		Nationality:   input.Natinality,
		RefreshToken:  input.RefreshToken,
	}, nil
}

func ToUpdatePlayerReq(input *pb.UpdatePlayerRequest) (*UpdatePlayerReq, error) {
	var dateOfBirth *time.Time
	if input.DateOfBirth != nil {
		t := input.DateOfBirth.AsTime()
		dateOfBirth = &t
	}

	var position *entity.PlayerPosition
	if input.Position != nil {
		p, err := toPlayerPosition(*input.Position)
		if err != nil {
			return nil, apperror.NewBadRequestError(err.Error())
		}
		position = &p
	}

	var preferredFoot *entity.PreferredFoot
	if input.PreferredFoot != nil {
		f, err := toPreferredFoot(*input.PreferredFoot)
		if err != nil {
			return nil, apperror.NewBadRequestError(err.Error())
		}
		preferredFoot = &f
	}

	return &UpdatePlayerReq{
		UserID:        input.UserId,
		PlayerID:      input.PlayerId,
		FullName:      input.FullName,
		DateOfBirth:   dateOfBirth,
		Position:      position,
		PreferredFoot: preferredFoot,
	}, nil
}

func ToGetPlayerReq(input *pb.GetPlayerReq) *GetPlayerReq {
	return &GetPlayerReq{
		PlayerID: input.PlayerId,
	}
}
