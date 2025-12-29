package igrpc

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/Verthandii/spring/ilogger"
)

type Option func(s *Server)

// Timeout with server timeout.
func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// Listener with server lis
func Listener(lis net.Listener) Option {
	return func(s *Server) {
		s.lis = lis
	}
}

// Options with grpc options.
func Options(opts ...grpc.ServerOption) Option {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}

type Server struct {
	*grpc.Server
	lis        net.Listener
	timeout    time.Duration
	health     *health.Server
	middleware Matcher
	grpcOpts   []grpc.ServerOption
}

func NewServer(opts ...Option) (*Server, error) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}

	srv := &Server{
		lis:        lis,
		timeout:    1 * time.Second,
		health:     health.NewServer(),
		middleware: New(),
	}
	for _, o := range opts {
		o(srv)
	}

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(srv.unaryServerInterceptor()),
		grpc.ChainStreamInterceptor(srv.streamServerInterceptor()),
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}

	srv.Server = grpc.NewServer(grpcOpts...)

	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	reflection.Register(srv.Server)

	return srv, nil
}

func (s *Server) Use(selector string, m ...Middleware) {
	s.middleware.Add(selector, m...)
}

func (s *Server) Run() error {
	ilogger.Info("grpc server ready", "addr", s.lis.Addr().String())
	return s.Serve(s.lis)
}

func (s *Server) Close(_ context.Context) {
	ilogger.Info("grpc server close")
	s.GracefulStop()
}
