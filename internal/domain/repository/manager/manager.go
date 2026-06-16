package manager

import (
	"context"

	"gorm.io/gorm"
)



type ManagerRepository interface{

}

type managerRepository struct{
    db *gorm.DB
}


func(mr *managerRepository)AddManager(ctx context.Context){

}


