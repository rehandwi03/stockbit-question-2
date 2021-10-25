package tests

import (
	"context"
	"fmt"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	http_handler "github.com/rehandwi03/stockbit-question-2/handler/http"
	"github.com/rehandwi03/stockbit-question-2/handler/rpc"
	"github.com/rehandwi03/stockbit-question-2/middleware"
	"github.com/rehandwi03/stockbit-question-2/model"
	"github.com/rehandwi03/stockbit-question-2/proto/movie"
	"github.com/rehandwi03/stockbit-question-2/repository"
	"github.com/rehandwi03/stockbit-question-2/service"
	"github.com/soheilhy/cmux"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
)

var (
	db          *gorm.DB
	resourceVar *dockertest.Resource
	poolVar     *dockertest.Pool
	url         = "http://localhost:10000"
	movieRepo   repository.Repository
	movieSvc    service.Service
)

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	poolVar = pool

	opts := dockertest.RunOptions{
		Repository: "mariadb",
		Tag:        "latest",
		Env: []string{
			"MYSQL_DATABASE=stockbit",
			"MYSQL_ROOT_PASSWORD=secret",
		},
		ExposedPorts: []string{"3306"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"3306": {
				{HostIP: "0.0.0.0", HostPort: "3307"},
			},
		},
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	resourceVar = resource

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(
		func() error {
			var err error
			dsn := fmt.Sprintf(
				"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "secret", "localhost",
				"3307", "stockbit",
			)
			db, err = gorm.Open(mysql.Open(dsn))

			if err != nil {
				return err
			}
			dbPing, _ := db.DB()
			return dbPing.Ping()
		},
	); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// init godotenv
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("error when loading env file ", err)
	}

	if err := db.AutoMigrate(&model.Log{}); err != nil {
		log.Fatalf("error migrate table: %v", err)
	}

	// init log repository
	logRepo := repository.NewLogRepository(db)

	// inject dependency log repo to middleware
	middleware.NewLogConn(logRepo)

	// init router
	gMux := gorilla_mux.NewRouter()
	// use middleware HttpLog
	gMux.Use(middleware.HttpLog)

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
	cmuxInit := cmux.New(listener)
	httpListener := cmuxInit.Match(cmux.HTTP1Fast())
	grpcListener := cmuxInit.MatchWithWriters(
		cmux.HTTP2MatchHeaderFieldSendSettings(
			"content-type", "application/grpc",
		),
	)
	http2 := cmuxInit.Match(cmux.HTTP2())

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
			return cmuxInit.Serve()
		},
	)
	log.Println("Run HTTP and GRPC server on port:", os.Getenv("PORT"))
	log.Println("Application started")

	// if err := g.Wait(); err != nil {
	// 	log.Fatal(err)
	// }

	// You can't defer this because os.Exit doesn't care for defer
	// if err := pool.Purge(resource); err != nil {
	// 	log.Fatalf("Could not purge resource: %s", err)
	// }

	code := m.Run()

	os.Exit(code)
}

func TestMovieFetchHTTP(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodGet, fmt.Sprintf("%s/movies?searchword=%s", url, "Batman"),
		nil,
	)

	assert.Nil(t, err)

	client := http.Client{}

	response, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	response.Body.Close()
}

func TestMovieGetByIDHTTP(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodGet, fmt.Sprintf("%s/movies/%s", url, "tt0372784"),
		nil,
	)

	assert.Nil(t, err)

	client := http.Client{}

	response, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	response.Body.Close()
}

func TestMovieFetchGRPC(t *testing.T) {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	assert.Nil(t, err)

	defer conn.Close()

	client := movie.NewMovieClient(conn)

	movies, err := client.Fetch(context.Background(), &movie.MovieRequest{Searchworld: "Batman"})
	assert.Nil(t, err)
	assert.NotNil(t, movies)
	assert.NotEmpty(t, movies)
}

func TestMovieGetByIDGRPC(t *testing.T) {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	assert.Nil(t, err)

	defer conn.Close()

	client := movie.NewMovieClient(conn)

	movies, err := client.GetByID(context.Background(), &movie.MovieDetailRequest{Id: "tt0372784"})
	assert.Nil(t, err)
	assert.NotNil(t, movies)
	assert.NotEmpty(t, movies)
}

func TestTearDown(t *testing.T) {
	if err := poolVar.Purge(resourceVar); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
