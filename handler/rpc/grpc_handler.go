package rpc

import (
	"context"
	"errors"
	"github.com/rehandwi03/stockbit-question-2/middleware"
	"github.com/rehandwi03/stockbit-question-2/proto/movie"
	"github.com/rehandwi03/stockbit-question-2/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

type movieServer struct {
	service service.Service
}

func (m *movieServer) GetByID(ctx context.Context, request *movie.MovieDetailRequest) (
	*movie.MovieDetailResponse, error,
) {
	if request.Id == "" {
		return nil, errors.New("id param can't be null")
	}

	res, err := m.service.GetByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := new(movie.MovieDetailResponse)
	response.Title = res.Title
	response.Year = res.Year
	response.Rated = res.Rated
	response.Released = res.Released
	response.Runtime = res.Runtime
	response.Genre = res.Genre
	response.Director = res.Director
	response.Writer = res.Writer
	response.Actors = res.Actors
	response.Plot = res.Plot
	response.Language = res.Language
	response.Country = res.Country
	response.Awards = res.Awards
	response.Poster = res.Awards
	response.Metasource = res.Metasource
	response.Type = res.Type
	response.DVD = res.DVD
	response.BoxOffice = res.BoxOffice
	response.Production = res.Production
	response.Website = res.Website
	response.Response = res.Response

	for _, val := range res.Ratings {
		var data movie.Rating

		data.Value = val.Value
		data.Source = val.Source

		response.Ratings = append(response.Ratings, &data)
	}

	return response, nil
}

func (m *movieServer) Fetch(ctx context.Context, request *movie.MovieRequest) (*movie.MovieResponse, error) {
	pagination := ""
	if request.Pagination == "" {
		pagination = "1"
	}

	if request.Searchworld == "" {
		return nil, errors.New("searchworld can't be null")
	}

	params := map[string]interface{}{
		"query": map[string]interface{}{
			"page": pagination,
			"s":    request.Searchworld,
		},
	}

	res, err := m.service.Fetch(ctx, params)
	if err != nil {
		return nil, err
	}

	movies := make([]*movie.MovieModel, 0)

	for _, val := range res {
		var data movie.MovieModel

		data.Title = val.Title
		data.Type = val.Type
		data.ImdbID = val.ImdbID
		data.Poster = val.Poster
		data.Year = val.Year

		movies = append(movies, &data)
	}

	response := movie.MovieResponse{
		Message: "success",
		Movies:  movies,
	}

	return &response, nil
}

func NewGrpcHandler(l net.Listener, service service.Service) error {
	s := grpc.NewServer(
		withServerUnaryInterceptor(),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionIdle: 5 * time.Minute,
			},
		),
	)
	movie.RegisterMovieServer(s, &movieServer{service: service})

	return s.Serve(l)
}

func withServerUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(middleware.GrpcInterceptor)
}
