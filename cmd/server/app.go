package main

import (
	"context"
	"github.com/arandich/marketplace-gateway/internal/config"
	"github.com/arandich/marketplace-gateway/internal/model"
	"github.com/arandich/marketplace-gateway/internal/repository"
	"github.com/arandich/marketplace-gateway/internal/service"
	"github.com/arandich/marketplace-gateway/internal/transport/gateway"
	grpcTransport "github.com/arandich/marketplace-gateway/internal/transport/grpc"
	httpTransport "github.com/arandich/marketplace-gateway/internal/transport/http"
	"github.com/arandich/marketplace-gateway/internal/transport/http/handlers"
	pb "github.com/arandich/marketplace-proto/api/proto/services"
	image_cdn "github.com/arandich/marketplace-sdk/image-cdn"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func runApp(ctx context.Context, cfg config.Config) {
	logger := zerolog.Ctx(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt, os.Kill)

	// Prometheus.
	promMetrics := initMetrics(ctx, cfg.Prometheus)

	idClientConn, err := initIdService(ctx, cfg.IdClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing id client conn")
	}
	defer func() {
		if err = idClientConn.Close(); err != nil {
			logger.Error().Err(err).Msg("error closing id service client connection")
		}
	}()

	orderClientConn, err := initOrderdService(ctx, cfg.OrdersClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing orders client conn")
	}
	defer func() {
		if err = orderClientConn.Close(); err != nil {
			logger.Error().Err(err).Msg("error closing orders service client connection")
		}
	}()

	goodsClientConn, err := initGoodsService(ctx, cfg.GoodsClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing goods client conn")
	}
	defer func() {
		if err = goodsClientConn.Close(); err != nil {
			logger.Error().Err(err).Msg("error closing goods service client connection")
		}
	}()

	clients := model.Clients{
		IdService:    pb.NewIdServiceClient(idClientConn),
		OrderService: pb.NewOrderServiceClient(orderClientConn),
		GoodsService: pb.NewGoodsServiceClient(goodsClientConn),
	}
	services := model.Services{}
	gatewayService := &service.MarketplaceService{}

	gatewayService = service.New(repository.New(ctx, services, clients, cfg))

	// HTTP.
	httpServer := httpTransport.NewHTTPServer(cfg.HTTP)
	if err := httpServer.StartHTTPServer(ctx, cfg.HTTP); err != nil {
		logger.Error().Err(err).Msg("error starting http server")
	}

	// GRPC.
	grpcServer := grpcTransport.NewGRPCServer(ctx, cfg.GRPC, gatewayService, promMetrics)
	if err != nil {
		logger.Fatal().Err(err).Msg("error creating grpc server")
	}
	if err = grpcServer.StartGRPCServer(ctx); err != nil {
		logger.Error().Err(err).Msg("error starting grpc server")
	}

	// CDN
	cdnClient := image_cdn.New(&http.Client{}, image_cdn.Config{
		Host: cfg.CDN.Host,
		Port: cfg.CDN.Port,
	})

	handlerCdn := handlers.NewCdnClientHandler(cdnClient)

	// GRPC Gateway.
	httpGateway := gateway.NewGRPCGateway(grpcServer.GetGRPCServer())
	if err = httpGateway.StartGRPCGateway(ctx, gatewayService, cfg.HttpGateway, handlerCdn); err != nil {
		logger.Error().Err(err).Msg("error starting grpc-gateway server")
	}

	logger.Info().Str("service", cfg.App.Name).Msg("service started")

	sig := <-c
	logger.Warn().Str("signal", sig.String()).Msg("received shutdown signal")
}
