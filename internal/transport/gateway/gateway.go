package gateway

import (
	"context"
	"github.com/arandich/marketplace-gateway/internal/config"
	"github.com/arandich/marketplace-gateway/internal/service"
	"github.com/arandich/marketplace-gateway/internal/transport/http/handlers"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
)

type Gateway struct {
	wrappedGrpc *grpcweb.WrappedGrpcServer
}

func NewGRPCGateway(grpcSrv *grpc.Server) *Gateway {
	originFunc := func(origin string) bool {
		return true
	}

	options := []grpcweb.Option{
		grpcweb.WithOriginFunc(originFunc),
		grpcweb.WithAllowedRequestHeaders(defaultAllowedHTTPHeaders),
	}

	wrappedGrpc := grpcweb.WrapServer(grpcSrv, options...)

	return &Gateway{wrappedGrpc: wrappedGrpc}
}

func (g *Gateway) StartGRPCGateway(ctx context.Context, marketplaceService *service.MarketplaceService, cfg config.HttpGatewayConfig, cdnHandler handlers.CdnClient) error {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("starting grpc-gateway server")

	handler := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					Multiline:      true,
					Indent:         "    ",
					AllowPartial:   true,
					UseProtoNames:  true,
					UseEnumNumbers: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			authStr := req.Header.Get("Authorization")
			return metadata.New(map[string]string{
				"authorization": authStr,
			})
		}),
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
		runtime.WithRoutingErrorHandler(runtime.DefaultRoutingErrorHandler),
	)

	router := mux.NewRouter()
	router.Path("/private/cdn/upload").
		HandlerFunc(cdnHandler.UploadImg).
		Methods(defaultAllowedHTTPMethods...)

	router.PathPrefix("/").
		Handler(handler).
		Methods(defaultAllowedHTTPMethods...)

	// Register GRPC services for grpc-gateway.
	if err := g.registerGRPCServices(ctx, handler, marketplaceService); err != nil {
		return err
	}

	// Cors options.
	opts := cors.Options{
		AllowedOrigins:      []string{"*"},
		AllowedMethods:      defaultAllowedHTTPMethods,
		AllowedHeaders:      defaultAllowedHTTPHeaders,
		AllowPrivateNetwork: true,
	}
	c := cors.New(opts)

	srv := &http.Server{
		Handler:           c.Handler(router),
		Addr:              cfg.Address,
		WriteTimeout:      cfg.WriteTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.ListenAndServe()
	}()

	return nil
}
