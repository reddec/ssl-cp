package db

import (
	"fmt"

	"gorm.io/gorm"
)

func New(dialector gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("initialize database: %w", err)
	}

	err = db.AutoMigrate(&Certificate{}, &Domain{}, &Unit{}, &Address{})
	if err != nil {
		return nil, fmt.Errorf("apply migrations: %w", err)
	}

	return db, err
}
