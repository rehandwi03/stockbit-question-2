package middleware

import (
	"context"
	"github.com/rehandwi03/stockbit-question-2/model"
	"github.com/rehandwi03/stockbit-question-2/repository"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"google.golang.org/grpc"
	"log"
	"net/http"
)

var logRepository repository.LogRepository

func NewLogConn(repo repository.LogRepository) {
	logRepository = repo
}

func GrpcInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	p, _ := peer.FromContext(ctx)

	md, _ := metadata.FromIncomingContext(ctx)

	endpoint := info.FullMethod
	clientIP := p.Addr.String()
	serverIP := md[":authority"]

	done := make(chan bool, 1)
	go func(done chan bool) {
		_, err := logRepository.Save(
			context.Background(), model.Log{
				ClientIP: clientIP,
				ServerIP: serverIP[0],
				Method:   "",
				URL:      endpoint,
				Protocol: "RPC",
			},
		)

		if err != nil {
			log.Printf("can't insert to db: %v", err)
		}

		done <- true
	}(done)

	log.Printf(
		"[Client IP - %s] [Server IP - %s] [Method - %s] [Endpoint - %s] [Protocol - %s]", clientIP, serverIP,
		"", endpoint,
		"RPC",
	)

	h, err := handler(ctx, req)

	<-done

	return h, err
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			done := make(chan bool, 1)
			go func(done chan bool) {
				_, err := logRepository.Save(
					context.Background(), model.Log{
						ClientIP: getClientIp(r),
						ServerIP: r.RemoteAddr,
						Method:   r.Method,
						URL:      r.URL.String(),
						Protocol: "HTTP",
					},
				)

				if err != nil {
					log.Printf("can't insert to db: %v", err)
				}

				done <- true
			}(done)
			log.Printf(
				"[Client IP - %s] [Server IP - %s] [Method - %s] [Endpoint - %s] [Protocol - %s]", getClientIp(r),
				r.RemoteAddr, r.Method, r.URL, "HTTP",
			)

			<-done

			next.ServeHTTP(w, r)
		},
	)
}

func getClientIp(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}

	return r.RemoteAddr
}
