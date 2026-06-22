package player

import (
	"time"

	"github.com/google/uuid"
)

type AddPlayerReq struct {
	TeamID       uuid.UUID
	FullName     string
	DOB          time.Time
	JerseyNumber int32
	Postion      string
	Height       float32
	Weight       float32
}

type AddPlayerRes struct {
}

type RemovePlayerReq struct {
}

type RemovePlayerRes struct {
}

type SuspendPlayerReq struct {
}

type SuspendPlayerRes struct {
}

type GetTeamPlayersReq struct {
}

type GetTeamPlayersRes struct {
}

type GetPlayerReq struct {
}

type GetPlayerRes struct {
}
