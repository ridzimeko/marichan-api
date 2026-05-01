package database

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	// Enable pgvector extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		return err
	}

	return db.AutoMigrate(
		&RequestLog{},
		&ChatMemory{},
	)
}
