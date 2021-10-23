package main

import (
	"context"
	"github.com/rehandwi03/stockbit-question-2/proto/movie"
	"google.golang.org/grpc"
	"log"
)

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		log.Printf("error dial: %v", err)
	}
	defer conn.Close()

	client := movie.NewMovieClient(conn)

	movies, err := client.Fetch(context.Background(), &movie.MovieRequest{Searchworld: "Batman"})
	if err != nil {
		log.Printf("error fetch: %v", err)
	}

	movie, err := client.GetByID(context.Background(), &movie.MovieDetailRequest{Id: "tt0372784"})
	if err != nil {
		log.Printf("error fetch: %v", err)
	}

	log.Println("movies: ", movies)
	log.Println("movie: ", movie)

}
