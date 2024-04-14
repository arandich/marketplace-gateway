package repository

import (
	"context"
	"github.com/arandich/marketplace-gateway/internal/config"
	"github.com/arandich/marketplace-gateway/internal/model"
	pb "github.com/arandich/marketplace-proto/api/proto/services"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MarketplaceRepository struct {
	services model.Services
	clients  model.Clients
	logger   *zerolog.Logger
	cfg      config.Config
}

func New(ctx context.Context, services model.Services, clients model.Clients, cfg config.Config) *MarketplaceRepository {
	return &MarketplaceRepository{
		services: services,
		clients:  clients,
		logger:   zerolog.Ctx(ctx),
		cfg:      cfg,
	}
}

func (m MarketplaceRepository) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	return m.clients.IdService.Auth(ctx, req)
}

func (m MarketplaceRepository) InitHold(ctx context.Context, req *pb.InitHoldRequest) (*pb.InitHoldResponse, error) {
	return m.clients.IdService.InitHold(ctx, req)
}

func (m MarketplaceRepository) GetUser(ctx context.Context, req *emptypb.Empty) (*pb.GetUserResponse, error) {
	return m.clients.IdService.GetUser(ctx, req)
}

func (m MarketplaceRepository) GetGood(ctx context.Context, req *pb.GetGoodRequest) (*pb.GetGoodResponse, error) {
	return m.clients.GoodsService.GetGood(ctx, req)
}

func (m MarketplaceRepository) GetGoods(ctx context.Context, req *pb.GetGoodsRequest) (*pb.GetGoodsResponse, error) {
	return m.clients.GoodsService.GetGoods(ctx, req)
}

func (m MarketplaceRepository) AddGood(ctx context.Context, req *pb.AddGoodRequest) (*pb.AddGoodResponse, error) {
	return m.clients.GoodsService.AddGood(ctx, req)
}

func (m MarketplaceRepository) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	return m.clients.OrderService.GetOrder(ctx, req)
}

func (m MarketplaceRepository) InitOrder(ctx context.Context, req *pb.InitOrderRequest) (*pb.InitOrderResponse, error) {
	return m.clients.OrderService.InitOrder(ctx, req)
}

func (m MarketplaceRepository) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	return m.clients.IdService.RegisterUser(ctx, req)
}
