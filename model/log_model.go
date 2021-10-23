package model

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type Log struct {
	ID         int64
	ClientIP   string
	ServerIP   string
	Method     string
	URL        string
	Protocol   string
	CreatedAt  sql.NullTime
	ModifiedAt sql.NullTime
}

func (l *Log) BeforeCreate(tx *gorm.DB) (err error) {
	l.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	l.ModifiedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return
}
