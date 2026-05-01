package database

import (
	"time"

	"github.com/pgvector/pgvector-go"
)

type RequestLog struct {
	ID        uint   `gorm:"primaryKey"`
	Feature   string `gorm:"size:100:index"`
	Status    string `gorm:"size:50"`
	Method    string `gorm:"size:20"`
	Path      string `gorm:"size:255"`
	ClientIP  string `gorm:"size:100"`
	UserAgent string `gorm:"type:text"`
	CreatedAt time.Time
}

type ChatMemory struct {
	ID        uint            `gorm:"primaryKey"`
	SessionID string          `gorm:"size:100;index"`
	Role      string          `gorm:"size:20"`
	Content   string          `gorm:"type:text"`
	Embedding pgvector.Vector `gorm:"type:vector(768)"` // Updated to 768 dimensions
	CreatedAt time.Time       `gorm:"index"`
}
