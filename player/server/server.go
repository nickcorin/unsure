package server

import (
	"context"
	"github.com/luno/jettison/j"
	"github.com/nickcorin/unsure/player/ops"
	"github.com/nickcorin/unsure/player/playerpb/protocp"

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
	pl, err := ops.GetParts(ctx, srv.b, req.ExternalId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list parts for round",
			j.KV("external_id", req.ExternalId))
	}

	// Convert parts to proto.
	var parts []*pb.Part
	for _, p := range pl {
		partProto, err := protocp.PartToProto(&p)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert part to proto")
		}

		parts = append(parts, partProto)
	}

	return &pb.GetPartsResp{Parts: parts}, nil
}

// GetRound returns a local rounds from a Player's DB.
func (srv *Server) GetRound(ctx context.Context, req *pb.GetRoundReq) (
	*pb.GetRoundResp, error) {
	r, err := rounds.Lookup(ctx, srv.b.PlayerDB(), req.RoundId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup round")
	}

	// Convert round to proto.
	roundProto, err := protocp.RoundToProto(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert round to proto")
	}
	return &pb.GetRoundResp{Round: roundProto}, nil
}
