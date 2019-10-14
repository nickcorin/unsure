package grpc

import (
	"context"
	"flag"

	"github.com/luno/jettison/errors"
	"github.com/luno/reflex"
	"github.com/luno/reflex/reflexpb"
	"google.golang.org/grpc"

	"github.com/nickcorin/unsure/player"
	pb "github.com/nickcorin/unsure/player/playerpb"
	"github.com/nickcorin/unsure/player/server"
)

var _ player.Client = (*client)(nil)

var grpcAddr = flag.String("player_address", "",
	"host:port of the player gRPC service")

// IsEnabled returns whether a gRPC address has been provided to the Player.
func IsEnabled() bool {
	return *grpcAddr != ""
}

type clientOpt func(c *client)

// WithAddress provides an option to specify the gRPC address of a Player
// client.
func WithAddress(address string) clientOpt {
	return func(c *client) {
		c.address = address
	}
}

// New returns a gRPC client for a Player.
func New(opts ...clientOpt) (player.Client, error) {
	c := client{
		address: *grpcAddr,
	}

	for _, o := range opts {
		o(&c)
	}

	var err error
	c.rpcConn, err = server.NewGRPCClient(c.address)
	if err != nil {
		return nil, err
	}

	c.rpcClient = pb.NewPlayerClient(c.rpcConn)

	return &c, nil
}

type client struct {
	address   string
	rpcConn   *grpc.ClientConn
	rpcClient pb.PlayerClient
}

// StreamEvents returns a reflex.StreamClient that can be used to
// stream reflex events from a Player.
func (c *client) StreamEvents(ctx context.Context, after string,
	opts ...reflex.StreamOption) (reflex.StreamClient, error) {

	streamFn := reflex.WrapStreamPB(func(ctx context.Context,
		req *reflexpb.StreamRequest) (reflex.StreamClientPB, error) {
		return c.rpcClient.StreamRoundEvents(ctx, req)
	})

	return streamFn(ctx, after, opts...)
}

// GetParts returns a Player's parts received for a given round.
func (c *client) GetParts(ctx context.Context, roundID int64) (
	[]player.Part, error) {
	res, err := c.rpcClient.GetParts(ctx, &pb.GetPartsReq{
		RoundId: roundID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get parts")
	}

	return res.Parts, nil
}

// GetRank returns a Player's rank received for a given round.
func (c *client) GetRank(ctx context.Context, roundID int64) (int32, error) {
	res, err := c.rpcClient.GetRank(ctx, &pb.GetRankReq{
		RoundId: roundID,
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to get rank")
	}

	return res.Rank, nil
}
