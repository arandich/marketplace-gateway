package service

import (
	"context"
	pb "github.com/arandich/marketplace-proto/api/proto/services"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MarketplaceRepository interface {
	IdRepository
	GoodsRepository
	OrdersRepository
}

type GoodsRepository interface {
	GetGood(ctx context.Context, req *pb.GetGoodRequest) (*pb.GetGoodResponse, error)
	GetGoods(ctx context.Context, req *pb.GetGoodsRequest) (*pb.GetGoodsResponse, error)
	AddGood(ctx context.Context, req *pb.AddGoodRequest) (*pb.AddGoodResponse, error)
}

type OrdersRepository interface {
	GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error)
	InitOrder(ctx context.Context, req *pb.InitOrderRequest) (*pb.InitOrderResponse, error)
}

type IdRepository interface {
	Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error)
	InitHold(ctx context.Context, req *pb.InitHoldRequest) (*pb.InitHoldResponse, error)
	GetUser(ctx context.Context, req *emptypb.Empty) (*pb.GetUserResponse, error)
	RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error)
}

var _ MarketplaceRepository = (*MarketplaceService)(nil)

type MarketplaceService struct {
	pb.UnimplementedIdServiceServer
	pb.UnimplementedGoodsServiceServer
	pb.UnimplementedOrderServiceServer
	repository MarketplaceRepository
}

func New(repository MarketplaceRepository) *MarketplaceService {
	return &MarketplaceService{
		repository: repository,
	}
}

func (s *MarketplaceService) GetGood(ctx context.Context, req *pb.GetGoodRequest) (*pb.GetGoodResponse, error) {
	return s.repository.GetGood(ctx, req)
}

func (s *MarketplaceService) GetGoods(ctx context.Context, req *pb.GetGoodsRequest) (*pb.GetGoodsResponse, error) {
	return s.repository.GetGoods(ctx, req)
}

func (s *MarketplaceService) AddGood(ctx context.Context, req *pb.AddGoodRequest) (*pb.AddGoodResponse, error) {
	return s.repository.AddGood(ctx, req)
}

func (s *MarketplaceService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	return s.repository.GetOrder(ctx, req)
}

func (s *MarketplaceService) InitOrder(ctx context.Context, req *pb.InitOrderRequest) (*pb.InitOrderResponse, error) {
	return s.repository.InitOrder(ctx, req)
}

func (s *MarketplaceService) InitHold(ctx context.Context, req *pb.InitHoldRequest) (*pb.InitHoldResponse, error) {
	return s.repository.InitHold(ctx, req)
}

func (s *MarketplaceService) GetUser(ctx context.Context, req *emptypb.Empty) (*pb.GetUserResponse, error) {
	return s.repository.GetUser(ctx, req)
}

func (s *MarketplaceService) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	return s.repository.Auth(ctx, req)
}

func (s *MarketplaceService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	return s.repository.RegisterUser(ctx, req)
}
