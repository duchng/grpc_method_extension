package app

import (
	"context"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	gardenservicev1 "grpc_method_extension/gen/api/garden-service/v1"
	apiv1 "grpc_method_extension/gen/api/v1"
	"grpc_method_extension/internal/infra/handlers"
	grpc_server "grpc_method_extension/pkg/grpc-server"
)

func ServerGrpc(ctx context.Context) error {
	gardenService := handlers.NewGardenService()
	return grpc_server.Serve(
		ctx, func(server *grpc.Server) {
			gardenservicev1.RegisterGardenServiceServer(server, gardenService)
		}, grpc.UnaryInterceptor(RequireTierUnaryInterceptor),
	)
}

func RequireTierUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "metadata not found")
	}
	var (
		userTier int
		err      error
	)
	for _, t := range md.Get("tier") {
		userTier, err = strconv.Atoi(t)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid tier")
		}
	}
	if userTier < int(getMinimumUserTier[int32](info.FullMethod, apiv1.E_MinimumTier)) {
		return nil, status.Error(codes.PermissionDenied, "not allowed to get flowers")
	}
	return handler(ctx, req)
}

func getMinimumUserTier[T any](fullMethodName string, extensionType protoreflect.ExtensionType) T {
	var defaultTier T
	methodParts := strings.Split(fullMethodName, "/")
	if len(methodParts) != 3 {
		return defaultTier
	}
	serviceName, methodName := methodParts[1], methodParts[2]

	// Find the service descriptor
	serviceDescriptor, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(serviceName))
	if err != nil {
		return defaultTier
	}

	// Find the method descriptor
	serviceDesc, ok := serviceDescriptor.(protoreflect.ServiceDescriptor)
	if !ok {
		return defaultTier
	}
	methodDesc := serviceDesc.Methods().ByName(protoreflect.Name(methodName))
	if methodDesc == nil {
		return defaultTier
	}
	// Check for the custom option
	opts := methodDesc.Options().(*descriptorpb.MethodOptions)
	if proto.HasExtension(opts, extensionType) {
		return proto.GetExtension(opts, extensionType).(T)
	}

	return defaultTier
}
