package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"

	healthpb "github.com/istsh/go-grpc-health-probe-sample/app/interface/rpc/v1/health"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
)

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: false, EmitDefaults: true}),
	)

	if err := registerHandlers(ctx, mux); err != nil {
		return err
	}

	handler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "Accept-Encoding", "Accept"}),
	)(mux)

	addr := ":8080"
	fmt.Printf("http server started on %s\n", addr)
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(addr, handler)
}

func registerHandlers(ctx context.Context, mux *runtime.ServeMux) error {
	opts := grpcDialOptions()

	if err := healthpb.RegisterHealthHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
		return err
	}

	return nil
}

func grpcDialOptions() []grpc.DialOption {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	}

	return opts
}
