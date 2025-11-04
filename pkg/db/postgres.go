package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Gorm -.
type Gorm struct {
	*gorm.DB
}

// New -.

func New(url string) (*Gorm, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {

		return nil, err
	}

	return &Gorm{db}, nil
}
