package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	server "github.com/vbncursed/medialog/auth-service/internal/api/auth_service_api"
	"github.com/vbncursed/medialog/auth-service/internal/pb/auth_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AppRun поднимает gRPC сервер и gRPC-Gateway (HTTP).
func AppRun(api server.AuthServiceAPI, grpcAddr, httpAddr string) {
	go func() {
		if err := runGRPCServer(api, grpcAddr); err != nil {
			panic(fmt.Errorf("failed to run gRPC server: %v", err))
		}
	}()

	if err := runGatewayServer(httpAddr, grpcAddr); err != nil {
		panic(fmt.Errorf("failed to run gateway server: %v", err))
	}
}

func runGRPCServer(api server.AuthServiceAPI, grpcAddr string) error {
	addr := grpcAddr
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(unaryLoggingInterceptor()),
	)
	auth_api.RegisterAuthServiceServer(s, &api)

	slog.Info("gRPC-server listening", "addr", addr)
	return s.Serve(lis)
}

func runGatewayServer(httpAddr, grpcAddr string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	swaggerPath := os.Getenv("swaggerPath")
	if swaggerPath != "" {
		if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
			panic(fmt.Errorf("swagger file not found: %s", swaggerPath))
		}
	}

	r := chi.NewRouter()

	if swaggerPath != "" {
		r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, swaggerPath)
		})

		r.Get("/docs/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger.json"),
		))
	}

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := auth_api.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return err
	}

	r.Mount("/", mux)

	slog.Info("gRPC-Gateway server listening", "addr", httpAddr)
	return http.ListenAndServe(httpAddr, r)
}
