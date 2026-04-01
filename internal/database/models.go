package database

import "time"

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
