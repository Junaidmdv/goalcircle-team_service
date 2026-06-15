package team

import "gorm.io/gorm"

type TeamRepository interface {
}

type teamRepository struct {
	db *gorm.DB
}

