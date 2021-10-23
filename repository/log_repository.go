package repository

import (
	"context"
	"github.com/rehandwi03/stockbit-question-2/model"
	"gorm.io/gorm"
)

type LogRepository interface {
	Save(ctx context.Context, model model.Log) (int64, error)
}

type logRepository struct {
	conn *gorm.DB
}

func NewLogRepository(conn *gorm.DB) LogRepository {
	return &logRepository{conn: conn}
}

func (l logRepository) Save(ctx context.Context, model model.Log) (int64, error) {
	err := l.conn.WithContext(ctx).Save(&model).Error
	if err != nil {
		return 0, err
	}

	return model.ID, nil
}