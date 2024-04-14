package grpc

import (
	"context"
	"fmt"
	"github.com/arandich/marketplace-gateway/internal/config"
	"github.com/arandich/marketplace-gateway/internal/service"
	pb "github.com/arandich/marketplace-proto/api/proto/services"
	loggerInterceptor "github.com/arandich/marketplace-sdk/interceptors/logger"
	recoveryInterceptor "github.com/arandich/marketplace-sdk/interceptors/recovery"
	sdkPrometheus "github.com/arandich/marketplace-sdk/prometheus"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	srv    *grpc.Server
	cfg    config.GrpcConfig
	logger *zerolog.Logger
}

func NewGRPCServer(ctx context.Context, cfg config.GrpcConfig, marketplaceService *service.MarketplaceService, promMetrics sdkPrometheus.Metrics) *Server {
	opts := []grpc.ServerOption{
		grpc.ConnectionTimeout(cfg.ConnectionTimeout),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     cfg.MaxConnectionIdle,
			MaxConnectionAge:      cfg.MaxConnectionAge,
			MaxConnectionAgeGrace: cfg.MaxConnectionAgeGrace,
			Timeout:               cfg.KeepAliveTimeout,
		}),
		grpc.ChainUnaryInterceptor(
			recoveryInterceptor.UnaryServerInterceptor(),
			loggerInterceptor.NewUnaryLoggerInterceptor(ctx),
			promMetrics.UnaryServerInterceptor(ctx),
		),
	}

	srv := grpc.NewServer(opts...)

	pb.RegisterGoodsServiceServer(srv, marketplaceService)
	pb.RegisterIdServiceServer(srv, marketplaceService)
	pb.RegisterOrderServiceServer(srv, marketplaceService)

	// Can be enabled if develop branch.
	reflection.Register(srv)

	return &Server{
		srv:    srv,
		logger: zerolog.Ctx(ctx),
		cfg:    cfg,
	}
}

func (s *Server) StartGRPCServer(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("starting grpc server")

	listener, err := net.Listen("tcp", s.cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to listen on address %s: %w", s.cfg.Address, err)
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.srv.Serve(listener)
	}()

	return nil
}

func (s *Server) GetGRPCServer() *grpc.Server {
	return s.srv
}
