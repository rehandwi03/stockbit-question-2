package service

import (
	"context"
	"github.com/rehandwi03/stockbit-question-2/model"
	"github.com/rehandwi03/stockbit-question-2/repository"
	"log"
	"time"
)

type service struct {
	repository repository.Repository
}

func NewMovieService(repo repository.Repository) Service {
	return &service{repository: repo}
}

func (s *service) Fetch(ctx context.Context, params map[string]interface{}) ([]model.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	res, err := s.repository.Fetch(ctx, params)
	if err != nil {
		log.Printf("error fetch in service: %v", err)
		return nil, err
	}

	return res.Search, nil
}

func (s *service) GetByID(ctx context.Context, id string) (movie model.MovieDetail, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	res, err := s.repository.GetByID(ctx, id)
	if err != nil {
		log.Printf("error fetch in service: %v", err)
		return movie, err
	}

	return res, nil
}
