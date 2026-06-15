package manager

import "gorm.io/gorm"



type ManagerRepository interface{

}

type managerRepository struct{
    db *gorm.DB
}



