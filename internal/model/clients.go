package model

import pb "github.com/arandich/marketplace-proto/api/proto/services"

type Clients struct {
	IdService    pb.IdServiceClient
	OrderService pb.OrderServiceClient
	GoodsService pb.GoodsServiceClient
}
