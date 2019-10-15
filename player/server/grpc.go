package server

import (
	"context"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
	"github.com/luno/jettison/log"
	"google.golang.org/grpc"
	"net"
)

// GRPCServer wraps a grpc.Server.
type GRPCServer struct {
	listener net.Listener
	srv      *grpc.Server
}

// Listener returns the net.Listener for the GRPCServer type.
func (srv *GRPCServer) Listener() net.Listener {
	return srv.listener
}

// GRPCServer returns the *grpc.Server for the GRPCServer type.
func (srv *GRPCServer) GRPCServer() *grpc.Server {
	return srv.srv
}

// ServeForever starts the gRPC server until it encounters an error or is
// stopped.
func (srv *GRPCServer) ServeForever() error {
	log.Info(context.Background(), "Player online", j.KV("addr",
		srv.listener.Addr()))
	return srv.srv.Serve(srv.listener)
}

// Stop attempts to perform a graceful stop on the server.
func (srv *GRPCServer) Stop() {
	srv.srv.GracefulStop()
}

// NewGRPCServer returns a gRPC server that we can register our service
// servers with.
func NewGRPCServer(address string) (*GRPCServer, error) {
	if address == "" {
		return nil, errors.New("no grpc address provided")
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create listener")
	}

	return &GRPCServer{
		listener: listener,
		srv:      grpc.NewServer(),
	}, nil
}

// NewGRPCClient returns a client connection to the given url.
func NewGRPCClient(url string) (*grpc.ClientConn, error) {
	return grpc.Dial(url, grpc.WithInsecure())
}
