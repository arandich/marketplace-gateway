package http

import (
	"context"
	"crypto/tls"
	"github.com/arandich/marketplace-gateway/internal/config"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
	"net/http/pprof"
)

type Server struct {
	server *http.Server
	cfg    config.HttpConfig
}

func NewHTTPServer(cfg config.HttpConfig) *Server {
	server := &http.Server{
		Addr:              cfg.Address,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return &Server{
		server: server,
		cfg:    cfg,
	}
}

func (s *Server) StartHTTPServer(ctx context.Context, cfg config.HttpConfig) error {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("starting http server")

	router := httprouter.New()

	router.Handler(http.MethodGet, "/metrics", promhttp.Handler())

	if s.cfg.ProfilingEnabled {
		router.HandlerFunc(http.MethodGet, "/debug/pprof/", pprof.Index)
		router.HandlerFunc(http.MethodGet, "/debug/pprof/cmdline", pprof.Cmdline)
		router.HandlerFunc(http.MethodGet, "/debug/pprof/profile", pprof.Profile)
		router.HandlerFunc(http.MethodGet, "/debug/pprof/symbol", pprof.Symbol)
		router.HandlerFunc(http.MethodGet, "/debug/pprof/trace", pprof.Trace)
		router.Handler(http.MethodGet, "/debug/pprof/heap", pprof.Handler("heap"))
		router.Handler(http.MethodGet, "/debug/pprof/goroutine", pprof.Handler("goroutine"))
		router.Handler(http.MethodGet, "/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		router.Handler(http.MethodGet, "/debug/pprof/block", pprof.Handler("block"))
	}

	s.server.Handler = router

	listener, err := net.Listen(cfg.Network, cfg.Address)
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.server.Serve(listener)
	}()

	return nil
}

func (s *Server) Close(ctx context.Context) {
	logger := zerolog.Ctx(ctx)

	ctx, cancel := context.WithTimeout(ctx, time.Minute*10)
	defer cancel()

	if err := s.server.Close(); err != nil {
		logger.Error().Err(err).Msg("error closing http server")
	}
}
