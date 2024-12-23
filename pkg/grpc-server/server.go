package grpc_server

import (
	"context"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type RegisterFunc func(*grpc.Server)

func Serve(ctx context.Context, registerFunc RegisterFunc, opts ...grpc.ServerOption) error {
	server := grpc.NewServer(opts...)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	reflection.Register(server)
	for name := range server.GetServiceInfo() {
		healthServer.SetServingStatus(name, grpc_health_v1.HealthCheckResponse_SERVING)
	}
	registerFunc(server)

	listener, err := net.Listen("tcp", ":9556")
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		slog.Info("shutting down grpc server..")
		healthServer.Shutdown()
		server.GracefulStop()
		listener.Close()
	}()
	slog.Info("starting grpc server..", slog.String("address", listener.Addr().String()))
	return server.Serve(listener)
}
