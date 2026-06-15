package team

import (
	"context"

	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	team_uc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/team"
)

type TeamHandler struct {
	tu team_uc.TeamUsecase
	pb.UnimplementedTeamServiceServer
}

func NewTeamHandler(tu team_uc.TeamUsecase) *TeamHandler {
	return &TeamHandler{
		tu: tu,
	}
}



func (th *TeamHandler) CreateTeam(ctx context.Context, req *pb.CreateTeamReq) (*pb.CreateTeamRes, error) {
   
	

	return nil, nil
} 


func(th *TeamHandler)UpdateTeamDetails(ctx context.Context){
  
}



