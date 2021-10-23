package service

import (
	"context"
	"github.com/rehandwi03/stockbit-question-2/model"
)

type Service interface {
	Fetch(ctx context.Context, params map[string]interface{}) (movies []model.Movie, err error)
	GetByID(ctx context.Context, id string) (movie model.MovieDetail, err error)
}