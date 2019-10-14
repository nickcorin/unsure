package server

import (
	"context"

	"github.com/luno/jettison/errors"
	"github.com/luno/reflex"
	"github.com/luno/reflex/reflexpb"

	"github.com/nickcorin/unsure/player/internal/db/rounds"
	pb "github.com/nickcorin/unsure/player/playerpb"
)

var _ pb.PlayerServer = (*Server)(nil)

// Server defines the dependencies required for a Player's gRPC server.
type Server struct {
	b       Backends
	rserver *reflex.Server
	stream  reflex.StreamFunc
}

// New returns an instance to the Player's gRPC server.
func New(b Backends) *Server {
	return &Server{
		b:       b,
		rserver: reflex.NewServer(),
		stream:  rounds.EventStream(b.PlayerDB()),
	}
}

// StreamRoundEvents returns a reflex.StreamClient that can be used to
// stream reflex events from a Player.
func (srv *Server) StreamRoundEvents(req *reflexpb.StreamRequest,
	ss pb.Player_StreamRoundEventsServer) error {
	return nil
}

// GetParts returns a Player's parts received for a given round.
func (srv *Server) GetParts(ctx context.Context, req *pb.GetPartsReq) (
	*pb.GetPartsResp, error) {
	return nil, errors.New("not implemented")
}

// GetRank returns a Player's rank received for a given round.
func (srv *Server) GetRank(ctx context.Context, req *pb.GetRankReq) (
	*pb.GetRankResp, error) {
	return nil, errors.New("not implemented")
}
