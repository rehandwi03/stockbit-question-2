package main

import (
	"fmt"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/joho/godotenv"
	http_handler "github.com/rehandwi03/stockbit-question-2/handler/http"
	"github.com/rehandwi03/stockbit-question-2/handler/rpc"
	"github.com/rehandwi03/stockbit-question-2/middleware"
	"github.com/rehandwi03/stockbit-question-2/model"
	"github.com/rehandwi03/stockbit-question-2/repository"
	"github.com/rehandwi03/stockbit-question-2/service"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"os"
)

func initDB() (db *gorm.DB) {
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASSWORD")
	dbPort := os.Getenv("DATABASE_PORT")
	dbName := os.Getenv("DATABASE_LOG_NAME")
	dbUrl := os.Getenv("DATABASE_URL")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbUrl,
		dbPort, dbName,
	)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatalf("error when connecting to db: %v", err)
	}

	return db

}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	log.Println("Init application...")

	// init godotenv
	if err := godotenv.Load(); err != nil {
		log.Println("error when loading env file ", err)
	}

	// init connection to mysql database
	db := initDB()

	// create table from model
	if err := db.AutoMigrate(&model.Log{}); err != nil {
		log.Println("error when migrating table to db ", err)
	}

	// init log repository
	logRepo := repository.NewLogRepository(db)

	// inject dependency log repo to middleware
	middleware.NewLogConn(logRepo)

	// init router
	gMux := gorilla_mux.NewRouter()
	// use middleware Log
	gMux.Use(middleware.Log)

	// init movie repo, svc, and http handler
	repo := repository.NewMovieRepository()
	svc := service.NewMovieService(repo)
	http_handler.NewMovieHttpHandler(gMux, svc)

	// use port 10000 for run application
	listener, err := net.Listen("tcp", ":"+os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	// init cmux
	m := cmux.New(listener)
	httpListener := m.Match(cmux.HTTP1Fast())
	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	http2 := m.Match(cmux.HTTP2())

	g := new(errgroup.Group)

	g.Go(
		func() error {
			return rpc.NewGrpcHandler(grpcListener, svc)
		},
	)
	g.Go(
		func() error {
			s := &http.Server{Handler: gMux}
			return s.Serve(httpListener)
		},
	)

	g.Go(
		func() error {
			s := &http.Server{Handler: gMux}
			return s.Serve(http2)
		},
	)

	g.Go(
		func() error {
			return m.Serve()
		},
	)
	log.Println("Run HTTP and GRPC server on port:", os.Getenv("PORT"))
	log.Println("Application started")

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
	// grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/json"))

	// run http and grpc simultaneously
	// g := new(errgroup.Group)
	// g.Go(func() error { return rpc.NewGrpcHandler(grpcListener, svc) })
	// g.Go(
	// 	func() error {
	// 		s := &http.Server{Handler: gMux}
	// 		return s.Serve(httpListener)
	// 	},
	// )
	// g.Go(func() error { return m.Serve() })

	// if err := g.Wait(); err != nil {
	// 	log.Fatal(err)
	// }

}
