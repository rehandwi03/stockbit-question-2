package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rehandwi03/stockbit-question-2/model"
	"github.com/rehandwi03/stockbit-question-2/util"
	"log"
	"os"
)

type Repository interface {
	Save(ctx context.Context) error
	Fetch(ctx context.Context, params map[string]interface{}) (movies model.Movies, err error)
	GetByID(ctx context.Context, id string) (movie model.MovieDetail, err error)
}

type repository struct {
}

func NewMovieRepository() Repository {
	return &repository{}
}

func (r repository) doHit(ctx context.Context, params map[string]interface{}) (
	res *resty.Response,
	err error,
) {
	url := os.Getenv("OMDB_URL")

	headers := map[string]string{}

	resHttp, err := util.ApiClientGet(ctx, url, params["endpoint"].(string), headers)
	if err != nil {
		log.Println(err)
		return res, err
	}

	if resHttp.StatusCode() != 200 {
		log.Println("failed")
		return res, errors.New("get product failed")
	}

	return resHttp, nil
}

func (r repository) GetByID(ctx context.Context, id string) (movie model.MovieDetail, err error) {
	paramHit := map[string]interface{}{
		"endpoint": fmt.Sprintf("/?apikey=%s&i=%s", os.Getenv("OMDB_API_KEY"), id),
	}

	res, err := r.doHit(ctx, paramHit)
	if err != nil {
		log.Printf("error fetch in repo: %v", err)
		return movie, err
	}

	err = json.Unmarshal(res.Body(), &movie)
	if err != nil {
		return movie, err
	}

	if movie.Response != "True" {
		var response model.MovieDetailResponse
		err = json.Unmarshal(res.Body(), &response)
		if err != nil {
			return movie, err
		}
		return movie, errors.New(response.Error)
	}

	return movie, nil

}

func (r repository) Fetch(ctx context.Context, params map[string]interface{}) (movies model.Movies, err error) {
	queryString := ""

	if params["query"] != nil {
		for index, value := range params["query"].(map[string]interface{}) {
			if queryString == "" {
				queryString = fmt.Sprintf("%s%s=%v", queryString, index, value)
			} else {
				queryString = fmt.Sprintf("%s&%s=%v", queryString, index, value)
			}
		}
	}

	paramHit := map[string]interface{}{
		"endpoint": fmt.Sprintf("/?apikey=%s&%s", os.Getenv("OMDB_API_KEY"), queryString),
	}

	res, err := r.doHit(ctx, paramHit)
	if err != nil {
		log.Printf("error fetch in repo: %v", err)
		return movies, err
	}

	err = json.Unmarshal(res.Body(), &movies)
	if err != nil {
		log.Println(err)
		return movies, err
	}

	if movies.Response == "False" {
		log.Println(movies.Error)
		return movies, errors.New(movies.Error)

	}

	return movies, nil
}

func (r repository) Save(ctx context.Context) error {
	panic("implement me")
}
