package player

import (
	"errors"
	"time"

	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/pariz/gountries"
)

func ToAddPlayerReq(input *pb.PlayerDetails) (*AddPlayerReq, error) {

	position, err := toPlayerPosition(&input.Position)
	if err != nil {
		return nil, apperror.NewBadRequestError(err.Error())
	}
	preferedfoot, err := toPreferredFoot(&input.PreferedFoot)
	if err != nil {
		return nil, apperror.NewBadRequestError(err.Error())
	}


	return &AddPlayerReq{
		UserID:        input.UserId,
		FullName:      input.FullName,
		DateOfBirth:   input.Dob.AsTime(),
		Position:      position,
		PreferredFoot: preferedfoot,
		Nationality:   input.Nationality,
		RefreshToken:  input.RefreshToken,
		Height:        input.Height,
		Weight:        input.Weight,
	}, nil
}

func ToUpdatePlayerReq(input *pb.UpdatePlayerRequest) (*UpdatePlayerReq, error) {
	var dateOfBirth *time.Time
	if input.DateOfBirth != nil {
		t := input.DateOfBirth.AsTime()
		dateOfBirth = &t
	}

	position, err := toPlayerPosition(input.Position)
	if err != nil {
		return nil, apperror.NewBadRequestError(err.Error())

	}

	preferedFoot, err := toPreferredFoot(input.PreferredFoot)
	if err != nil {
		return nil, apperror.NewBadRequestError(err.Error())
	}

	return &UpdatePlayerReq{
		UserID:        input.UserId,
		FullName:      input.FullName,
		DateOfBirth:   dateOfBirth,
		Position:      &position,
		PreferredFoot: &preferedFoot,
		Weight:        input.Weight,
		Height:        input.Height,
		Nationality:   input.Nationality,
	}, nil
}

func ToGetPlayerReq(input *pb.GetPlayerReq) *GetPlayerReq {
	return &GetPlayerReq{
		PlayerID: input.PlayerId,
	}
}

func GetCountry(nationality string) (string, error) {
	if nationality == "" {
		return "", nil
	}
	query := gountries.New()
	country, err := query.FindCountryByAlpha(nationality)
	if err != nil {
		return "", errors.New("invalid nationality")
	}
	return country.Name.Common, nil
}
