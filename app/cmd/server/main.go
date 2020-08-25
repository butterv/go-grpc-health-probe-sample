package main

import (
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/reflection"

	"github.com/istsh/go-grpc-health-probe-sample/app/domain/service"
	healthpb "github.com/istsh/go-grpc-health-probe-sample/app/interface/rpc/v1/health"
)

func main() {
	listenPort, err := net.Listen("tcp", ":9090")
	if err != nil {
		logrus.Fatalln(err)
	}

	s := newGRPCServer()
	reflection.Register(s)
	_ = s.Serve(listenPort)
	s.GracefulStop()
}

func newGRPCServer() *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	healthpb.RegisterHealthServer(s, service.NewHealthService())

	return s
}
