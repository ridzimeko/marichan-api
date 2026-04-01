package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(env *Env) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		env.DBHost,
		env.DBPort,
		env.DBUser,
		env.DBPass,
		env.DBName,
		env.DBSSLMode,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
