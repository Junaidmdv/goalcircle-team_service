package team

import "context"







type TeamUsecase interface{
	 
}


func NewTeamUsecase()TeamUsecase{
	return &teamUsecase{}
}





type teamUsecase struct{

}



func(tu *teamUsecase)CreateTeam(ctx context.Context,)
