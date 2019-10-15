package main

import (
	"context"
	"flag"
	"github.com/nickcorin/unsure/player/ops"
	"os"
	"os/signal"
	"time"

	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/log"

	"github.com/nickcorin/unsure/player/playerpb"
	"github.com/nickcorin/unsure/player/server"
	"github.com/nickcorin/unsure/player/state"
)

var grpcAddress = flag.String("grpc_address", "", "player grpc address")

func main() {
	flag.Parse()

	s, err := state.New()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create player state"))
	}

	grpcSrv, err := server.NewGRPCServer(*grpcAddress)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create grpc server"))
	}

	// Register player gRPC server.
	playerpb.RegisterPlayerServer(grpcSrv.GRPCServer(), server.New(s))

	// Serve forever.
	go func() {
		log.Fatal(grpcSrv.ServeForever())
	}()

	// Start event loops.
	ops.StartLoops(s)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive interrupt.
	<-c

	// Wait for deadline.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)
	defer cancel()

	// Gracefully shutdown.
	grpcSrv.Stop()
	log.Info(ctx, "Shutting down")
	os.Exit(0)
}