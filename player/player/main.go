package main

import (
	"flag"

	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/log"

	"github.com/nickcorin/unsure/player/playerpb"
	"github.com/nickcorin/unsure/player/server"
	"github.com/nickcorin/unsure/player/state"
	"github.com/nickcorin/unsure/player/ops"
)

var grpcAddress = flag.String("grpc_address", "", "player grpc address")

func main() {
	flag.Parse()

	s, err := state.New()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create player state"))
	}

	go serveGRPCForever(s)
	ops.StartLoops(s)
}

func serveGRPCForever(s *state.State) {
	grpcSrv, err := server.NewGRPCServer(*grpcAddress)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create grpc server"))
	}

	playerpb.RegisterPlayerServer(grpcSrv.GRPCServer(), server.New(s))

	log.Fatal(grpcSrv.ServeForever())
}
