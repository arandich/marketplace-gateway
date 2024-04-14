package gateway

import (
	"context"
	"github.com/arandich/marketplace-gateway/internal/service"
	pb "github.com/arandich/marketplace-proto/api/proto/services"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func (g *Gateway) registerGRPCServices(ctx context.Context, mux *runtime.ServeMux, service *service.MarketplaceService) error {

	err := pb.RegisterGoodsServiceHandlerServer(ctx, mux, service)
	if err != nil {
		return err
	}

	err = pb.RegisterOrderServiceHandlerServer(ctx, mux, service)
	if err != nil {
		return err
	}

	err = pb.RegisterIdServiceHandlerServer(ctx, mux, service)
	if err != nil {
		return err
	}

	return nil
}
