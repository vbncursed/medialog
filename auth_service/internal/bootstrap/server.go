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
	server "github.com/vbncursed/medialog/auth_service/internal/api/auth_service_api"
	"github.com/vbncursed/medialog/auth_service/internal/pb/auth_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AppRun(api server.AuthServiceAPI) {
	go func() {
		if err := runGRPCServer(api); err != nil {
			panic(fmt.Errorf("failed to run gRPC server: %v", err))
		}
	}()

	if err := runGatewayServer(); err != nil {
		panic(fmt.Errorf("failed to run gateway server: %v", err))
	}
}

func runGRPCServer(api server.AuthServiceAPI) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(unaryLoggingInterceptor()),
	)
	auth_api.RegisterAuthServiceServer(s, &api)

	slog.Info("gRPC-server server listening on :50051")
	return s.Serve(lis)
}

func runGatewayServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	swaggerPath := os.Getenv("swaggerPath")
	if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
		panic(fmt.Errorf("swagger file not found: %s", swaggerPath))
	}

	r := chi.NewRouter()
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

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := auth_api.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, ":50051", opts)
	if err != nil {
		panic(err)
	}

	r.Mount("/", mux)

	slog.Info("gRPC-Gateway server listening on :8081")
	return http.ListenAndServe(":8081", r)
}
