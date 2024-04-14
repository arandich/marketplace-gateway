package main

import (
	"context"
	"github.com/arandich/marketplace-gateway/internal/config"
	"github.com/arandich/marketplace-sdk/authorization/jwtInterceptor"
	sdkPrometheus "github.com/arandich/marketplace-sdk/prometheus"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

func initHTTP(ctx context.Context, cfg config.HttpConfig) (net.Listener, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("initializing HTTP server listener")

	lis, err := net.Listen(cfg.Network, cfg.Address)
	if err != nil {
		return nil, err
	}

	return lis, nil
}

func initGRPC(ctx context.Context, cfg config.GrpcConfig) (net.Listener, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("initializing GRPC server listener")

	lis, err := net.Listen(cfg.Network, cfg.Address)
	if err != nil {
		return nil, err
	}

	return lis, nil
}

func initMetrics(ctx context.Context, cfg config.PrometheusConfig) sdkPrometheus.Metrics {
	logger := zerolog.Ctx(ctx)
	logger.Info().Str("namespace", cfg.Namespace).Str("subsystem", cfg.Subsystem).Msg("initializing prometheus metrics")

	promCfg := sdkPrometheus.Config{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
	}

	baseMetrics := sdkPrometheus.New(promCfg)

	return baseMetrics
}

func initIdService(ctx context.Context, cfg config.IdClientConfig) (*grpc.ClientConn, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("initializing referrals grpc client")

	opts := []grpc.DialOption{
		grpc.WithIdleTimeout(cfg.IdleTimeout),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(jwtInterceptor.ChainUnaryInterceptor),
	}

	conn, err := grpc.DialContext(ctx, cfg.ConnString, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func initOrderdService(ctx context.Context, cfg config.OrdersClientConfig) (*grpc.ClientConn, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("initializing referrals grpc client")

	opts := []grpc.DialOption{
		grpc.WithIdleTimeout(cfg.IdleTimeout),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(jwtInterceptor.ChainUnaryInterceptor),
	}

	conn, err := grpc.DialContext(ctx, cfg.ConnString, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func initGoodsService(ctx context.Context, cfg config.GoodsClientConfig) (*grpc.ClientConn, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("initializing referrals grpc client")

	opts := []grpc.DialOption{
		grpc.WithIdleTimeout(cfg.IdleTimeout),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(jwtInterceptor.ChainUnaryInterceptor),
	}

	conn, err := grpc.DialContext(ctx, cfg.ConnString, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
