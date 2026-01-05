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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/vbncursed/medialog/auth/config"
	"github.com/vbncursed/medialog/auth/internal/api/auth_service_api"
	"github.com/vbncursed/medialog/auth/internal/pb/auth_api"
)

func AppRun(api *auth_service_api.AuthServiceAPI, cfg *config.Config) {
	go func() {
		if err := runGRPCServer(api, cfg); err != nil {
			panic(fmt.Errorf("failed to run gRPC server: %v", err))
		}
	}()

	if err := runGatewayServer(cfg); err != nil {
		panic(fmt.Errorf("failed to run gateway server: %v", err))
	}
}

func runGRPCServer(api *auth_service_api.AuthServiceAPI, cfg *config.Config) error {
	lis, err := net.Listen("tcp", cfg.Server.GRPCAddr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	auth_api.RegisterAuthServiceServer(s, api)

	slog.Info("gRPC server listening", "addr", cfg.Server.GRPCAddr)
	return s.Serve(lis)
}

func runGatewayServer(cfg *config.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	swaggerPath := os.Getenv("swaggerPath")
	if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
		slog.Warn("swagger file not found, skipping swagger UI", "path", swaggerPath)
	}

	r := chi.NewRouter()

	// OpenAPI JSON endpoint
	if swaggerPath != "" {
		r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, swaggerPath)
		})

		// Scalar API Reference UI
		r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
      body {
        margin: 0;
      }
    </style>
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/swagger.json"
    ></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference@latest"></script>
  </body>
</html>`)
		})
	}

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := auth_api.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, cfg.Server.GRPCAddr, opts)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	r.Mount("/", mux)

	slog.Info("gRPC-Gateway server listening", "addr", cfg.Server.HTTPAddr)
	return http.ListenAndServe(cfg.Server.HTTPAddr, r)
}
