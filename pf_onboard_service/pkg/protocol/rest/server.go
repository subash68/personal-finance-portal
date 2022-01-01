package rest

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/subash68/ate/ate_token_service/pkg/api/token"
	"github.com/subash68/ate/ate_token_service/pkg/logger"
	"github.com/subash68/ate/ate_token_service/pkg/protocol/rest/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, grpcPort, httpPort string) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := v1.RegisterTokenServiceHandlerFromEndpoint(ctx, mux, ":"+grpcPort, opts); err != nil {
		logger.Log.Fatal("failed to start HTTP gateway", zap.String("reason", err.Error()))
	}

	srv := &http.Server{
		Addr: ":" + httpPort,
		Handler: middleware.AddRequestID(
			middleware.AddLogger(logger.Log, mux)),
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			// sig is a ^C, handle it
			// How ? :-)
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	logger.Log.Info("starting HTTP/REST gateway...")
	return srv.ListenAndServe()
}
