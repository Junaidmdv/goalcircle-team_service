package team

import (
	"context"

	"gorm.io/gorm"
)

type TeamRepository interface {
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{
		db: db,
	}
}




func(tr *teamRepository)CreateTeam(ctx context.Context,req  string){
     
}


func(tr *teamRepository)DeleteTeam(ctx context.Context,uuid string){

}