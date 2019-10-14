package protocp

import (
	"github.com/nickcorin/unsure/player"
	pb "github.com/nickcorin/unsure/player/playerpb"
)

// PartFromProto converts a pb.Part to a player.Part.
func PartFromProto(in *pb.Part) (*player.Part, error) {
	return &player.Part{
		ID:        in.Id,
		RoundID:   in.RoundId,
		Name:      in.Name,
		Rank:      in.Rank,
		Value:     in.Value,
		CreatedAt: in.CreatedAt,
	}, nil
}

// PartToProto converts a player.Part to a pb.Part.
func PartToProto(in *player.Part) (*pb.Part, error) {
	return &pb.Part{
		Id:        in.ID,
		RoundId:   in.RoundID,
		Name:      in.Name,
		Rank:      in.Rank,
		Value:     in.Value,
		CreatedAt: in.CreatedAt,
	}, nil
}
