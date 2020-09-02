package service

import (
	"context"

	healthpb "github.com/istsh/go-grpc-health-probe-sample/app/interface/rpc/v1/health"
)

type healthService struct {
}

// NewHealthService generates the `HealthServer` implementation.
func NewHealthService() healthpb.HealthServer {
	return &healthService{}
}

func (*healthService) Check(context.Context, *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	// 	0: UNKNOWN
	//	1: SERVING
	//	2: NOT_SERVING
	//	3: SERVICE_UNKNOWN
	return &healthpb.HealthCheckResponse{
		Status: healthpb.HealthCheckResponse_SERVING,
	}, nil
}
