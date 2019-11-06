package main

import (
	"flag"
	"github.com/corverroos/unsure"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/log"
	"unsure/player/ops"
	"unsure/player/playerpb"

	"unsure/player/server"
	"unsure/player/state"
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

	unsure.WaitForShutdown()
}

func serveGRPCForever(s *state.State) {
	grpcServer, err := unsure.NewServer(*grpcAddress)
	if err != nil {
		unsure.Fatal(errors.Wrap(err, "new grpctls server"))
	}

	playerSrv := server.New(s)
	playerpb.RegisterPlayerServer(grpcServer.GRPCServer(), playerSrv)

	unsure.RegisterNoErr(func() {
		playerSrv.Stop()
		grpcServer.Stop()
	})

	unsure.Fatal(grpcServer.ServeForever())
}
